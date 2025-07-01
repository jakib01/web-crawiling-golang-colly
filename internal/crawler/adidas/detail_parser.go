package adidas

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	// Fetch sizes via API
	sizes, _ := FetchProductSizes(code) // handle error optionally

	// Extract reviews
	reviews := ExtractReviews(doc)

	// Extract reviews
	aspectRatings := ExtractAspectRatings(doc)

	coordinatedItems := ExtractCoordinatedItems(doc)

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

	return model.Product{
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
	}, nil
}

func FetchProductSizes(code string) ([]model.ProductSize, error) {
	apiURL := fmt.Sprintf("https://www.adidas.jp/api/products/%s/availability", code)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		VariationList []struct {
			Size         string  `json:"size"`
			Availability float64 `json:"availability"`
		} `json:"variation_list"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	var sizes []model.ProductSize
	for _, v := range data.VariationList {
		sizes = append(sizes, model.ProductSize{
			ProductID:         0,
			SizeLabel:         v.Size,
			Availability:      v.Availability,
			ChestCM:           0.0,
			BackLengthCM:      0.0,
			OtherMeasurements: "",
			SpecialFunctions:  "",
		})
	}

	return sizes, nil
}

func ExtractReviews(doc *goquery.Document) []model.Review {
	var reviews []model.Review

	doc.Find(`[data-auto-id="review"]`).Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(`.title___`).Text())
		reviewDate := strings.TrimSpace(s.Find(`.date`).Text())
		body := strings.TrimSpace(s.Find(`.text___`).Text())

		// Count .gl-star-rating__item under the review's .gl-star-rating container
		rating := s.Find(".gl-star-rating .gl-star-rating__item").Length()
		parsedDate, err := time.Parse("2006年1月2日", reviewDate)
		if err != nil {
			parsedDate = time.Time{}
		}

		reviews = append(reviews, model.Review{
			Title:         title,
			ReviewDate:    parsedDate,
			Rating:        float64(rating),
			OverallRating: 0, // As per your input
			Body:          body,
		})
	})

	return reviews
}

func ExtractAspectRatings(doc *goquery.Document) []model.ReviewAspectRating {
	var aspectRatings []model.ReviewAspectRating

	doc.Find(".gl-comparison-bar").Each(func(i int, s *goquery.Selection) {
		aspect := strings.TrimSpace(s.Find(".gl-comparison-bar__title strong").Text())

		// Get "left: 62.5%;" from style attribute
		style, exists := s.Find(".gl-comparison-bar__indicator").Attr("style")
		rating := 0.0
		if exists {
			// Example style: transform: translateX(-62.5%); left: 62.5%;
			if parts := strings.Split(style, "left:"); len(parts) > 1 {
				raw := strings.TrimSpace(strings.TrimSuffix(parts[1], "%;"))
				rating, _ = strconv.ParseFloat(raw, 64)
			}
		}

		aspectRatings = append(aspectRatings, model.ReviewAspectRating{
			ReviewID: 0, // placeholder, assign after DB insert or join
			Aspect:   aspect,
			Rating:   rating,
		})
	})

	return aspectRatings
}

func ExtractCoordinatedItems(doc *goquery.Document) []model.CoordinatedItem {
	var items []model.CoordinatedItem

	doc.Find(`#gl-carousel-system-product-carousel-complete-the-look-recs-content > li`).Each(func(i int, s *goquery.Selection) {
		productNumber, _ := s.Attr("id")

		name := strings.TrimSpace(s.Find("h4._product-card-content-main__name_36dpn_83").Text())

		priceStr := s.Find(`[data-testid="main-price"] span`).Last().Text()
		priceYen := 0.0
		if priceStr != "" {
			clean := strings.ReplaceAll(priceStr, "¥", "")
			clean = strings.ReplaceAll(clean, ",", "")
			priceYen, _ = strconv.ParseFloat(clean, 64)
		}

		imageURL, _ := s.Find("img").Attr("src")
		productPageURL, _ := s.Find("a").Attr("href")

		items = append(items, model.CoordinatedItem{
			ProductNumber:  productNumber,
			Name:           name,
			PriceYen:       priceYen,
			ImageURL:       imageURL,
			ProductPageURL: productPageURL,
		})
	})

	return items
}
