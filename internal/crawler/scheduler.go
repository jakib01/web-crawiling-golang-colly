package crawler

import (
	"context"
	"log"
	"sync"

	"github.com/gocolly/colly"
	"github.com/jmoiron/sqlx"

	"github.com/jakib01/web-crawiling-golang-colly/internal/config"
	"github.com/jakib01/web-crawiling-golang-colly/internal/crawler/adidas"
)

type Scheduler struct {
	db    *sqlx.DB
	cfg   config.CrawlerConfig
	mutex sync.Mutex
	seen  map[string]bool
}

func NewScheduler(db *sqlx.DB, cfg config.CrawlerConfig) *Scheduler {
	return &Scheduler{
		db:   db,
		cfg:  cfg,
		seen: make(map[string]bool),
	}
}

func (s *Scheduler) Run(ctx context.Context, limit int) error {
	c := colly.NewCollector(
		colly.Async(true),
	)
	c.Limit(&colly.LimitRule{
		Parallelism: s.cfg.Concurrency,
	})

	count := 0
	c.OnHTML("a[href^='/products/']", func(e *colly.HTMLElement) {
		s.mutex.Lock()
		if count >= limit || s.seen[e.Request.AbsoluteURL(e.Attr("href"))] {
			s.mutex.Unlock()
			return
		}
		count++
		s.seen[e.Request.AbsoluteURL(e.Attr("href"))] = true
		s.mutex.Unlock()

		c.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})

	c.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting: %s", r.URL.String())
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {
		if !adidas.IsProductPage(e.Request.URL.String()) {
			return
		}
		dto, err := adidas.ParseProduct(e)
		if err != nil {
			log.Printf("failed to parse product: %v", err)
			return
		}
		log.Printf("Parsed product: %s (%s)", dto.Name, dto.ProductURL)
		// TODO: map dto → model → repo.Save()
	})

	if err := c.Visit(s.cfg.StartURL); err != nil {
		return err
	}
	c.Wait()
	return nil
}
