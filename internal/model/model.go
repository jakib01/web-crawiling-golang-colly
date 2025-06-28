package model

type Product struct {
	ProductNumber    string
	Name             string
	Category         string
	PriceYen         float64
	SenseOfSize      string
	DetailsURL       string
	TotalReviews     int
	RecommendedRate  float64
	Images           []ProductImage
	Sizes            []ProductSize
	Keywords         []string
	CoordinatedItems []CoordinatedProduct
	Reviews          []Review
}

type ProductImage struct {
	URL    string
	IsMain bool
}

type ProductSize struct {
	Label             string
	ChestCM           float64
	BackLengthCM      float64
	OtherMeasurements map[string]string
	SpecialFunctions  map[string]string
}

type CoordinatedProduct struct {
	ProductNumber string
	Name          string
	PriceYen      float64
	ImageURL      string
	PageURL       string
}

type Review struct {
	ReviewerID    string
	ReviewDate    string
	OverallRating float64
	Title         string
	Body          string
	AspectRatings map[string]float64
}
