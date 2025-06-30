package adidas

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/jakib01/web-crawiling-golang-colly/internal/model"
	"go.uber.org/zap"
)

const step = 48

func collectProductURLs(limit int, logger *zap.SugaredLogger) ([]model.ProductURL, error) {
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

	productMap := map[string]bool{}
	var productList []model.ProductURL
	start := 0

	for len(productList) < limit {
		pageURL := "https://www.adidas.jp/メンズ"
		if start > 0 {
			pageURL = fmt.Sprintf("%s?start=%d", pageURL, start)
		}

		var html string
		err := chromedp.Run(ctx,
			chromedp.Navigate(pageURL),
			chromedp.Sleep(6*time.Second),
			chromedp.OuterHTML("html", &html),
		)
		if err != nil {
			logger.Errorf("Failed to load %s: %v", pageURL, err)
			break
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			logger.Errorf("Failed to parse HTML: %v", err)
			break
		}

		found := 0
		doc.Find("a[href$='.html']").Each(func(_ int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if !exists || strings.Count(href, "/") != 2 || !strings.HasSuffix(href, ".html") {
				return
			}

			fullURL := "https://www.adidas.jp" + href
			if productMap[fullURL] {
				return
			}

			// ✅ Extract code from last segment of path
			parts := strings.Split(href, "/")
			code := strings.TrimSuffix(parts[len(parts)-1], ".html")

			imgURL := ""
			if img := s.Find("img"); img.Length() > 0 {
				imgURL, _ = img.Attr("src")
			}

			productList = append(productList, model.ProductURL{
				Code:      code,
				URL:       fullURL,
				ImageURL:  imgURL,
				ScrapedAt: time.Now(),
			})
			productMap[fullURL] = true
			found++

			if len(productList) >= limit {
				return
			}
		})

		if found == 0 {
			logger.Info("No more products found. Ending pagination.")
			break
		}

		start += step
	}

	return productList, nil
}
