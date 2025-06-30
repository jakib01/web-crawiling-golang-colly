package adidas

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
)

func collectProductURLs(limit int, logger *zap.SugaredLogger) ([]string, error) {
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

	productURLs := map[string]bool{}
	start := 0
	step := 48

	for len(productURLs) < limit {
		pageURL := "https://www.adidas.jp/メンズ"
		if start > 0 {
			pageURL = fmt.Sprintf("https://www.adidas.jp/メンズ?start=%d", start)
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

		foundOnPage := 0
		doc.Find("a[href$='.html']").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists && strings.HasSuffix(href, ".html") && strings.Count(href, "/") == 2 {
				fullURL := "https://www.adidas.jp" + href
				if !productURLs[fullURL] {
					productURLs[fullURL] = true
					foundOnPage++
					if len(productURLs) >= limit {
						return
					}
				}
			}
		})

		if foundOnPage == 0 {
			logger.Info("No more products found. Ending pagination.")
			break
		}

		start += step
	}

	var urls []string
	for url := range productURLs {
		urls = append(urls, url)
	}
	return urls, nil
}
