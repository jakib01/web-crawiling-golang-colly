package adidas

import (
	"encoding/json"
	"fmt"
	"github.com/jakib01/web-crawiling-golang-colly/internal/model"
	"github.com/jakib01/web-crawiling-golang-colly/internal/repository/postgres"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
)

type AdidasCrawler struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewAdidasCrawler(db *gorm.DB, logger *zap.SugaredLogger) *AdidasCrawler {
	return &AdidasCrawler{db: db, logger: logger}
}

func (c *AdidasCrawler) CrawlProducts(limit int) ([]model.ProductURL, error) {
	products, err := collectProductURLs(limit, c.logger)
	if err != nil {
		return nil, err
	}

	if err := postgres.StoreProductURLs(c.db, products); err != nil {
		return nil, err
	}

	var allDetails []model.Product

	for _, p := range products {
		detail, err := FetchAndParseDetailPage(p.URL, p.Code)
		if err != nil {
			c.logger.Warnf("failed to fetch detail for %s: %v", p.URL, err)
			continue
		}
		allDetails = append(allDetails, detail)

		fmt.Println(detail)

		// Marshal the entire slice as pretty-printed JSON
		out, err := json.MarshalIndent(allDetails, "", "  ")
		if err != nil {
			return nil, err
		}

		// Write to a single file
		if err := os.WriteFile("all_products.json", out, 0644); err != nil {
			return nil, err
		}

		//// save detail to DB (use GORM create or update)
		//if err := postgres.StoreProductDetail(c.db, detail); err != nil {
		//	c.logger.Warnf("failed to store detail for %s: %v", p.URL, err)
		//}
	}

	return products, nil
}
