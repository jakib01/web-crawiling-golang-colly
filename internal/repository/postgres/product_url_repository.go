package postgres

import (
	"github.com/jakib01/web-crawiling-golang-colly/internal/model"
	"gorm.io/gorm"
)

func StoreProductURLs(db *gorm.DB, entries []model.ProductURL) error {
	for _, p := range entries {
		var existing model.ProductURL
		err := db.Where("url = ?", p.URL).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&p).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err // other errors
		}
	}
	return nil
}
