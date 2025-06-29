package adidas

import (
	_ "context"
	"time"

	"github.com/jakib01/web-crawiling-golang-colly/internal/model"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type AdidasCrawler struct {
	db     *sqlx.DB
	logger *zap.SugaredLogger
}

func NewAdidasCrawler(db *sqlx.DB, logger *zap.SugaredLogger) *AdidasCrawler {
	return &AdidasCrawler{db: db, logger: logger}
}

func (c *AdidasCrawler) CrawlProducts(limit int) ([]model.Product, error) {
	urls, err := CrawlProductURLs(limit)
	if err != nil {
		return nil, err
	}

	var products []model.Product
	for _, url := range urls {
		p, err := FetchAndParseDetailPage(url)
		if err != nil {
			c.logger.Warnf("failed to fetch product: %s, error: %v", url, err)
			continue
		}
		p.ScrapedAt = time.Now()
		products = append(products, p)
	}

	return products, nil
}
