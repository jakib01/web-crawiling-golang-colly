package adidas

import (
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

func (c *AdidasCrawler) CrawlProducts(limit int) ([]string, error) {
	urls, err := collectProductURLs(limit, c.logger)
	if err != nil {
		return nil, err
	}
	return urls, nil
}
