package model

import "time"

type Product struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Slug      string    `db:"slug"`
	PriceRaw  string    `db:"price_raw"`
	ImageURL  string    `db:"image_url"`
	ScrapedAt time.Time `db:"scraped_at"`
}
