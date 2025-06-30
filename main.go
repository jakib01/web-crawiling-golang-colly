package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func main() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
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

	for len(productURLs) < 200 {
		pageURL := "https://www.adidas.jp/メンズ"
		if start > 0 {
			pageURL = fmt.Sprintf("https://www.adidas.jp/メンズ?start=%d", start)
		}

		fmt.Println("Visiting:", pageURL)

		var html string
		err := chromedp.Run(ctx,
			chromedp.Navigate(pageURL),
			chromedp.Sleep(6*time.Second), // adjust if slow
			chromedp.OuterHTML("html", &html),
		)
		if err != nil {
			log.Println("Failed to load:", pageURL, err)
			break
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			log.Println("Failed to parse HTML:", err)
			break
		}

		foundOnPage := 0
		doc.Find("a[href$='.html']").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists && strings.HasSuffix(href, ".html") && strings.Count(href, "/") == 2 {
				fullURL := "https://www.adidas.jp" + href
				if !productURLs[fullURL] {
					productURLs[fullURL] = true
					fmt.Println(fullURL)
					foundOnPage++
					if len(productURLs) >= 200 {
						return
					}
				}
			}
		})

		if foundOnPage == 0 {
			fmt.Println("🔚 No more products found. Stopping.")
			break
		}

		start += step
	}

	fmt.Printf("\n✅ Collected %d product URLs.\n", len(productURLs))
}
