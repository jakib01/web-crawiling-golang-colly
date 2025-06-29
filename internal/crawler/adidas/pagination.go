package adidas

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func CrawlProductURLs(max int) ([]string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 90*time.Second)
	defer cancel()

	var allLinks []string
	seen := map[string]bool{}
	for start := 0; len(allLinks) < max; start += 48 {
		pageURL := fmt.Sprintf("https://www.adidas.jp/メンズ?start=%d", start)
		var html string
		if err := chromedp.Run(ctx,
			chromedp.Navigate(pageURL),
			chromedp.Sleep(4*time.Second),
			chromedp.OuterHTML("html", &html),
		); err != nil {
			return nil, err
		}

		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			href, ok := s.Attr("href")
			if !ok || !strings.HasSuffix(href, ".html") || strings.Contains(href, "/products/") {
				return
			}
			if !seen[href] {
				seen[href] = true
				full := "https://www.adidas.jp" + href
				if decoded, err := url.PathUnescape(full); err == nil {
					allLinks = append(allLinks, decoded)
				}
			}
		})
		if len(allLinks) >= max {
			break
		}
	}
	return allLinks[:min(len(allLinks), max)], nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
