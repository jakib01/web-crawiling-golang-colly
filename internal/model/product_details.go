package model

import "time"

type Product struct {
	ID                         uint    `gorm:"primaryKey"`
	ProductCode                string  `gorm:"size:50;uniqueIndex;not null"`
	Name                       string  `gorm:"size:500;not null"`
	Category                   string  `gorm:"size:500;not null"`
	PriceYen                   float64 `gorm:"type:numeric(10,2);not null"`
	SenseOfSize                string  `gorm:"size:100"`
	DetailsURL                 string  `gorm:"type:text;not null"`
	TotalReviews               int     `gorm:"default:0"`
	OverallRating              float64 `gorm:"type:numeric(3,2);not null"`
	TitleDescription           string  `gorm:"type:text;not null"`
	GeneralDescription         string  `gorm:"type:text;not null"`
	ItemGeneralDescription     string  `gorm:"type:text;not null"`
	SpecialFunctionDescription string  `gorm:"type:text;not null"`

	Images        []ProductImage       `gorm:"foreignKey:ProductID"`
	Sizes         []ProductSize        `gorm:"foreignKey:ProductID"`
	Keywords      []Keyword            `gorm:"many2many:product_keywords"`
	Reviews       []Review             `gorm:"foreignKey:ProductID"`
	AspectRatings []ReviewAspectRating `gorm:"foreignKey:ReviewID"`
	Coordinated   []CoordinatedItem    `gorm:"foreignKey:SourceProductID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ProductImage struct {
	ID        uint   `gorm:"primaryKey"`
	ProductID uint   `gorm:"index"`
	URL       string `gorm:"type:text;not null"`
	IsMain    bool   `gorm:"default:false"`
}

type ProductSize struct {
	ID                uint    `gorm:"primaryKey"`
	ProductID         uint    `gorm:"index"`
	SizeLabel         string  `gorm:"size:20;not null"`
	ChestCM           float64 `gorm:"type:numeric(5,2)"`
	Availability      float64 `gorm:"type:numeric(5,2)"`
	BackLengthCM      float64 `gorm:"type:numeric(5,2)"`
	OtherMeasurements string  `gorm:"type:text"`
	SpecialFunctions  string  `gorm:"type:text"`
}

type CoordinatedItem struct {
	ID              uint    `gorm:"primaryKey"`
	SourceProductID uint    `gorm:"index"`
	ProductNumber   string  `gorm:"size:50;not null"`
	Name            string  `gorm:"size:500;not null"`
	PriceYen        float64 `gorm:"type:numeric(10,2);not null"`
	ImageURL        string  `gorm:"type:text;not null"`
	ProductPageURL  string  `gorm:"type:text;not null"`
}

type Keyword struct {
	ID uint   `gorm:"primaryKey"`
	Kw string `gorm:"size:100;unique;not null"`
}

type ProductKeyword struct {
	ProductID uint `gorm:"primaryKey"`
	KeywordID uint `gorm:"primaryKey"`
}

type Review struct {
	ID            uint      `gorm:"primaryKey"`
	ProductID     uint      `gorm:"index"`
	ReviewDate    time.Time `gorm:"type:date;not null"`
	Rating        float64   `gorm:"type:numeric(3,2);not null"`
	OverallRating float64   `gorm:"type:numeric(3,2);not null"`
	Title         string    `gorm:"size:255"`
	Body          string    `gorm:"type:text"`
}

type ReviewAspectRating struct {
	ID       uint    `gorm:"primaryKey"`
	ReviewID uint    `gorm:"index"`
	Aspect   string  `gorm:"size:100;not null"`
	Rating   float64 `gorm:"type:numeric(3,2);not null"`
}
type ProductDetail struct {
	ID          uint   `gorm:"primaryKey"`
	Code        string `gorm:"uniqueIndex;not null"`
	Name        string
	Price       string
	Description string
	ImageURL    string
	URL         string
	ScrapedAt   time.Time
}
