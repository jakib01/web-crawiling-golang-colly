package adidas

import (
	"strings"

	_ "github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// ProductDTO holds minimal parsed fields for now
type ProductDTO struct {
	ProductURL string
	Name       string
}

// IsProductPage checks if the URL is a product detail page
func IsProductPage(url string) bool {
	return strings.Contains(url, "/products/")
}

// ParseProduct extracts product info from the page
func ParseProduct(e *colly.HTMLElement) (*ProductDTO, error) {
	doc := e.DOM
	dto := &ProductDTO{
		ProductURL: e.Request.URL.String(),
		Name:       strings.TrimSpace(doc.Find(".product-name h1").Text()),
	}
	return dto, nil
}
