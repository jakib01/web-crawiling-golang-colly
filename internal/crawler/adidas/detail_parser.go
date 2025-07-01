package adidas

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/jakib01/web-crawiling-golang-colly/internal/model"
)

func FetchAndParseDetailPage(url string, code string) (model.Product, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var html string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(3*time.Second),
		chromedp.OuterHTML("html", &html),
	); err != nil {
		return model.Product{}, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return model.Product{}, err
	}

	// Extract fields
	name := strings.TrimSpace(doc.Find(`h1[data-auto-id="product-title"]`).Text())
	priceStr := doc.Find(`[data-testid="main-price"] span`).Last().Text()

	var category, titleDescription, generalDescription string
	reviewCount := 0

	doc.Find(`script[type="application/ld+json"]`).Each(func(i int, s *goquery.Selection) {
		var data map[string]interface{}
		raw := strings.TrimSpace(s.Text())
		if raw == "" || !strings.Contains(raw, `"@type"`) {
			return
		}
		if err := json.Unmarshal([]byte(raw), &data); err != nil {
			return
		}
		if data["@type"] == "Product" {
			if val, ok := data["category"].(string); ok {
				category = val
			}
			if val, ok := data["description"].(string); ok {
				generalDescription = val
			}
			// optional fallback for titleDescription
			//if titleDescription == "" {
			//	titleDescription = val
			//}
		}
	})

	// fallback titleDescription from DOM if available
	if domTitle := strings.TrimSpace(doc.Find("h3").First().Text()); domTitle != "" {
		titleDescription = domTitle
	}

	// reviewCount fallback from DOM
	if reviewStr := strings.TrimSpace(doc.Find(`button[data-auto-id="product-rating-review-count"]`).Text()); reviewStr != "" {
		reviewCount, _ = strconv.Atoi(strings.Trim(reviewStr, "件のレビュー "))
	}

	// Parse price (¥16,500 → 16500.00)
	priceYen := 0.0
	if priceStr != "" {
		cleaned := strings.ReplaceAll(priceStr, "¥", "")
		cleaned = strings.ReplaceAll(cleaned, ",", "")
		priceYen, _ = strconv.ParseFloat(cleaned, 64)
	}

	sizes, err := ExtractProductSizes(ctx)
	if err != nil {
		return model.Product{}, fmt.Errorf("extract sizes failed: %w", err)
	}

	// Extract reviews
	reviews, err := ExtractReviews(ctx)
	if err != nil {
		return model.Product{}, fmt.Errorf("extract reviews failed: %w", err)
	}

	// Extract reviews
	aspectRatings, err := ExtractAspectRatings(ctx)
	if err != nil {
		return model.Product{}, fmt.Errorf("extract aspectRatings failed: %w", err)
	}

	coordinatedItems, err := ExtractCoordinatedItems(doc)
	if err != nil {
		return model.Product{}, fmt.Errorf("extract coordinatedItems failed: %w", err)
	}

	// Extract product images
	var images []model.ProductImage
	doc.Find(`picture[data-testid="pdp-gallery-picture"] img`).Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists || strings.HasPrefix(src, "data:image") {
			return // skip placeholders or missing src
		}
		images = append(images, model.ProductImage{
			URL:    src,
			IsMain: false,
		})
	})

	var data = model.Product{
		ProductCode:                code,
		Name:                       name,
		Category:                   category,
		PriceYen:                   priceYen,
		SenseOfSize:                "",
		TotalReviews:               reviewCount,
		DetailsURL:                 url,
		TitleDescription:           titleDescription,
		GeneralDescription:         generalDescription,
		ItemGeneralDescription:     "",
		SpecialFunctionDescription: "",
		Sizes:                      sizes,
		Reviews:                    reviews,
		Images:                     images,
		AspectRatings:              aspectRatings,
		Coordinated:                coordinatedItems,
	}
	return data, nil
}

func ExtractProductSizes(ctx context.Context) ([]model.ProductSize, error) {
	// Wait for the size-selector section to become visible
	if err := chromedp.Run(ctx,
		chromedp.WaitVisible(`div.size-selector___2kfnl`, chromedp.ByQuery),
		chromedp.Sleep(1*time.Second),
	); err != nil {
		return nil, fmt.Errorf("size container not visible: %w", err)
	}

	// Pull all span texts under the sizes container, filtering out placeholder 'AAA'
	var labels []string
	js := `
        Array.from(
            document.querySelectorAll('div.sizes___2jQjF .gl-label span')
        ).map(el => el.textContent.trim())
         .filter(lbl => lbl && lbl !== 'AAA')
    `
	if err := chromedp.Run(ctx,
		chromedp.Evaluate(js, &labels),
	); err != nil {
		return nil, fmt.Errorf("evaluate size labels error: %w", err)
	}

	if len(labels) == 0 {
		return nil, fmt.Errorf("no sizes found in DOM via JS")
	}

	sizes := make([]model.ProductSize, len(labels))
	for i, lbl := range labels {
		sizes[i] = model.ProductSize{
			ProductID:    0,
			SizeLabel:    lbl,
			Availability: 1,
		}
	}
	return sizes, nil
}

func ExtractReviews(ctx context.Context) ([]model.Review, error) {
	// Expand the reviews accordion
	if err := chromedp.Run(ctx,
		chromedp.Click(`div[data-testid="accordion"] button.accordion__header___3Pii5`, chromedp.ByQuery),
		chromedp.Sleep(500*time.Millisecond),
	); err != nil {
		return nil, fmt.Errorf("accordion click error: %w", err)
	}

	// JS snippet to extract reviews
	js := `
        Array.from(
          document.querySelectorAll('div[data-auto-id="single-review-mobile"]')
        ).map(div => {
          const masks = Array.from(div.querySelectorAll('.gl-star-rating__mask')).map(m => parseFloat(m.style.width));
          const rating = masks.length ? (masks.reduce((a,b)=>a+b,0)/masks.length)/20 : 0;
          const dateEl = div.querySelector('.review-date___sEaVk');
          const dateText = dateEl ? dateEl.textContent.trim().replace('年','-').replace('月','-').replace('日','') : '';
          const titleEl = div.querySelector('.review-title___1382M strong');
          const bodyEl = div.querySelector('.review-description___21UXW .clamped___3Fp2g');
          const userEl = div.querySelector('.user-name___1n05v');
          const votes = div.querySelectorAll('.votes___3Q6JI span');
          const helpful = votes.length>1 ? votes[1].textContent.trim() : '';
          return {
            Title: titleEl? titleEl.textContent.trim() : '',
            Body: bodyEl? bodyEl.textContent.trim() : '',
            DateText: dateText,
            Rating: rating,
            UserName: userEl? userEl.textContent.trim(): '',
            Helpful: helpful
          };
        })
    `
	var raw []struct {
		Title    string
		Body     string
		DateText string
		Rating   float64
		UserName string
		Helpful  string
	}
	if err := chromedp.Run(ctx, chromedp.Evaluate(js, &raw)); err != nil {
		return nil, fmt.Errorf("evaluate reviews JS error: %w", err)
	}

	var reviews []model.Review
	for _, r := range raw {
		t, _ := time.Parse("2006-1-2", r.DateText)
		reviews = append(reviews, model.Review{
			Title:      r.Title,
			Body:       r.Body,
			Rating:     r.Rating,
			ReviewDate: t,
		})
	}
	return reviews, nil
}

// ExtractAspectRatings scrapes aspect ratings from the expanded review section via Chromedp.
func ExtractAspectRatings(ctx context.Context) ([]model.ReviewAspectRating, error) {
	// 1) Expand the reviews accordion if it's collapsed
	if err := chromedp.Run(ctx,
		chromedp.Click(`div[data-testid="accordion"] button.accordion__header___3Pii5`, chromedp.ByQuery),
		chromedp.Sleep(300*time.Millisecond),
	); err != nil {
		return nil, fmt.Errorf("accordion click error: %w", err)
	}

	// 2) JS snippet to read each comparison bar
	js := `
      Array.from(
        document.querySelectorAll('.sub-ratings___1pAhV .gl-comparison-bar')
      ).map(bar => {
        const aspect = bar.querySelector('.gl-comparison-bar__title strong')?.textContent.trim() || '';
        const style = bar.querySelector('.gl-comparison-bar__indicator')?.getAttribute('style') || '';
        const m = style.match(/left:\s*(\d+(?:\.\d+)?)%/);
        const pct = m ? parseFloat(m[1]) : 0;
        return { Aspect: aspect, Rating: pct };
      });
    `
	var raw []struct {
		Aspect string
		Rating float64
	}
	if err := chromedp.Run(ctx, chromedp.Evaluate(js, &raw)); err != nil {
		return nil, fmt.Errorf("evaluate aspects JS error: %w", err)
	}

	// 3) Map into your model.ReviewAspectRating
	aspects := make([]model.ReviewAspectRating, len(raw))
	for i, a := range raw {
		aspects[i] = model.ReviewAspectRating{Aspect: a.Aspect, Rating: a.Rating}
	}
	return aspects, nil
}

// ExtractCoordinatedItems pulls both style‐lookbook cards and the "complete the look" product carousel
func ExtractCoordinatedItems(doc *goquery.Document) ([]model.CoordinatedItem, error) {
	var items []model.CoordinatedItem

	// 1) Lookbook / style cards
	doc.Find(`div[data-testid="styles-carousel"] a[data-testid="style-card"]`).Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		// first image in the card
		img, _ := s.Find("span._imageWrap_1hxoi_9 img").First().Attr("src")
		headline := strings.TrimSpace(s.Find(`[data-testid="style-card-headline"]`).Text())
		desc := strings.TrimSpace(s.Find(`[data-testid="style-card-description"]`).Text())
		items = append(items, model.CoordinatedItem{
			ProductNumber:  "", // none for lookbook
			Name:           fmt.Sprintf("%s %s", headline, desc),
			PriceYen:       0, // no price
			ImageURL:       img,
			ProductPageURL: href,
		})
	})

	// 2) "Complete the look" product recommendations
	doc.Find(`#gl-carousel-system-product-carousel-complete-the-look-recs-content li`).Each(func(_ int, s *goquery.Selection) {
		id, _ := s.Attr("id")
		card := s.Find(`a._product-card__link_o6rgp_73`)
		href, _ := card.Attr("href")
		img, _ := card.Find(`img`).First().Attr("src")
		name := strings.TrimSpace(card.Find("h4").Text())
		priceStr := strings.TrimSpace(card.Find(`[data-testid="main-price"]`).Text())
		// clean "¥12,100" → "12100"
		clean := strings.ReplaceAll(strings.ReplaceAll(priceStr, "¥", ""), ",", "")
		priceVal, _ := strconv.ParseFloat(clean, 64)

		items = append(items, model.CoordinatedItem{
			ProductNumber:  id,
			Name:           name,
			PriceYen:       priceVal,
			ImageURL:       img,
			ProductPageURL: href,
		})
	})

	return items, nil
}
