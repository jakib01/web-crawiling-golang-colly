package adidas

import (
	"strings"

	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/jakib01/web-crawiling-golang-colly/internal/model"
)

func FetchAndParseDetailPage(link string) (model.Product, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var html string
	err := chromedp.Run(ctx,
		chromedp.Navigate(link),
		chromedp.Sleep(3),
		chromedp.OuterHTML("html", &html),
	)
	if err != nil {
		return model.Product{}, err
	}

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	name := doc.Find("h1").First().Text()
	price := doc.Find("div[data-auto-id='product-price']").Text()
	imageURL, _ := doc.Find("img").First().Attr("src")

	return model.Product{
		ID:       ExtractSKUFromURL(link),
		Name:     strings.TrimSpace(name),
		Slug:     link,
		PriceRaw: strings.TrimSpace(price),
		ImageURL: imageURL,
	}, nil
}

func ExtractSKUFromURL(link string) string {
	parts := strings.Split(link, "/")
	last := parts[len(parts)-1]
	return strings.TrimSuffix(last, ".html")
}
