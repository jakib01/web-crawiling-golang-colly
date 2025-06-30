package adidas

import (
	"github.com/jakib01/web-crawiling-golang-colly/internal/model"
	"github.com/jakib01/web-crawiling-golang-colly/internal/repository/postgres"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

	return products, nil
}
