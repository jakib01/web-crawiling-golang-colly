package model

import "time"

type ProductURL struct {
	ID        uint   `gorm:"primaryKey"`
	Code      string `gorm:"not null"`
	URL       string `gorm:"uniqueIndex;not null"`
	ImageURL  string
	ScrapedAt time.Time
}
