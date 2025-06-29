package postgres

import (
	"github.com/jakib01/web-crawiling-golang-colly/internal/model"
	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) BulkInsert(products []model.Product) error {
	tx := r.db.MustBegin()
	stmt := `INSERT INTO products (id, name, slug, price_raw, image_url, scraped_at)
	         VALUES (:id, :name, :slug, :price_raw, :image_url, :scraped_at)
	         ON CONFLICT (id) DO NOTHING`
	for _, p := range products {
		if _, err := tx.NamedExec(stmt, p); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
