package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/xuri/excelize/v2"
)

// --- your data models ---
type Product struct {
	ID                         int            `json:"ID"`
	ProductCode                string         `json:"ProductCode"`
	Name                       string         `json:"Name"`
	Category                   string         `json:"Category"`
	PriceYen                   int            `json:"PriceYen"`
	SenseOfSize                string         `json:"SenseOfSize"`
	DetailsURL                 string         `json:"DetailsURL"`
	TotalReviews               int            `json:"TotalReviews"`
	OverallRating              float64        `json:"OverallRating"`
	TitleDescription           string         `json:"TitleDescription"`
	GeneralDescription         string         `json:"GeneralDescription"`
	ItemGeneralDescription     string         `json:"ItemGeneralDescription"`
	SpecialFunctionDescription string         `json:"SpecialFunctionDescription"`
	Images                     []Image        `json:"Images"`
	Sizes                      []Size         `json:"Sizes"`
	Reviews                    []Review       `json:"Reviews"`
	AspectRatings              []AspectRating `json:"AspectRatings"`
	Coordinated                []Coordinated  `json:"Coordinated"`
}

type Image struct {
	ProductID int    `json:"ProductID"`
	URL       string `json:"URL"`
	IsMain    bool   `json:"IsMain"`
}

type Size struct {
	ProductID    int    `json:"ProductID"`
	SizeLabel    string `json:"SizeLabel"`
	Availability int    `json:"Availability"`
}

type Review struct {
	ProductID  int     `json:"ProductID"`
	ReviewDate string  `json:"ReviewDate"`
	Rating     float64 `json:"Rating"`
	Title      string  `json:"Title"`
	Body       string  `json:"Body"`
}

type AspectRating struct {
	ReviewID int     `json:"ReviewID"`
	Aspect   string  `json:"Aspect"`
	Rating   float64 `json:"Rating"`
}

type Coordinated struct {
	SourceProductID int    `json:"SourceProductID"`
	ProductNumber   string `json:"ProductNumber"`
	Name            string `json:"Name"`
	PriceYen        int    `json:"PriceYen"`
	ImageURL        string `json:"ImageURL"`
	ProductPageURL  string `json:"ProductPageURL"`
}

// --- helper to write one sheet ---
func writeSheet(f *excelize.File, sheet string, header []string, rows [][]interface{}) error {
	// Create sheet
	index, _ := f.NewSheet(sheet)
	// Write header
	for c, h := range header {
		cell, _ := excelize.CoordinatesToCellName(c+1, 1)
		if err := f.SetCellValue(sheet, cell, h); err != nil {
			return err
		}
	}
	// Write rows
	for r, row := range rows {
		for c, val := range row {
			cell, _ := excelize.CoordinatesToCellName(c+1, r+2)
			if err := f.SetCellValue(sheet, cell, val); err != nil {
				return err
			}
		}
	}
	f.SetActiveSheet(index)
	return nil
}

func main() {
	// 1) load JSON
	data, err := ioutil.ReadFile("all_products.json")
	if err != nil {
		log.Fatal(err)
	}
	var products []Product
	if err := json.Unmarshal(data, &products); err != nil {
		log.Fatal(err)
	}

	// 2) prepare rows for each sheet
	var (
		prodRows, imgRows, sizeRows, revRows, aspectRows, coordRows [][]interface{}
	)
	for _, p := range products {
		// product row
		prodRows = append(prodRows, []interface{}{
			p.ID, p.ProductCode, p.Name, p.Category, p.PriceYen, p.DetailsURL, p.TotalReviews, p.OverallRating,
		})
		// images
		for _, img := range p.Images {
			imgRows = append(imgRows, []interface{}{p.ID, img.URL, img.IsMain})
		}
		// sizes
		for _, s := range p.Sizes {
			sizeRows = append(sizeRows, []interface{}{p.ID, s.SizeLabel, s.Availability})
		}
		// reviews
		for _, r := range p.Reviews {
			revRows = append(revRows, []interface{}{p.ID, r.ReviewDate, r.Rating, r.Title})
		}
		// aspect ratings
		for _, a := range p.AspectRatings {
			aspectRows = append(aspectRows, []interface{}{a.ReviewID, a.Aspect, a.Rating})
		}
		// coordinated
		for _, c := range p.Coordinated {
			coordRows = append(coordRows, []interface{}{p.ID, c.ProductNumber, c.Name, c.PriceYen, c.ImageURL})
		}
	}

	// 3) create workbook & sheets
	f := excelize.NewFile()

	if err := writeSheet(f, "Products",
		[]string{"ProductID", "Code", "Name", "Category", "PriceYen", "DetailsURL", "TotalReviews", "OverallRating"},
		prodRows,
	); err != nil {
		log.Fatal(err)
	}
	if err := writeSheet(f, "Images",
		[]string{"ProductID", "URL", "IsMain"},
		imgRows,
	); err != nil {
		log.Fatal(err)
	}
	if err := writeSheet(f, "Sizes",
		[]string{"ProductID", "SizeLabel", "Availability"},
		sizeRows,
	); err != nil {
		log.Fatal(err)
	}
	if err := writeSheet(f, "Reviews",
		[]string{"ProductID", "ReviewDate", "Rating", "Title"},
		revRows,
	); err != nil {
		log.Fatal(err)
	}
	if err := writeSheet(f, "AspectRatings",
		[]string{"ReviewID", "Aspect", "Rating"},
		aspectRows,
	); err != nil {
		log.Fatal(err)
	}
	if err := writeSheet(f, "Coordinated",
		[]string{"ProductID", "ProductNumber", "Name", "PriceYen", "ImageURL"},
		coordRows,
	); err != nil {
		log.Fatal(err)
	}

	// 4) save file
	if err := f.SaveAs("product_details.xlsx"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Wrote product_details.xlsx")
}
