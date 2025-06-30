package main

import (
	"flag"
	"fmt"
	"github.com/jakib01/web-crawiling-golang-colly/internal/crawler/adidas"
	"os"

	_ "github.com/lib/pq"

	"github.com/jakib01/web-crawiling-golang-colly/internal/config"
	"github.com/jakib01/web-crawiling-golang-colly/internal/logger"
	//"github.com/jakib01/web-crawiling-golang-colly/internal/repository/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	envFile := flag.String("env", ".env", "path to env file")
	limit := flag.Int("limit", 10, "max number of products to crawl")
	flag.Parse()

	// ─── Load config ───────────────────────────────────────────
	cfg, err := config.Load(*envFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// ─── Init logger ───────────────────────────────────────────
	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	// get a SugaredLogger for fmt-style methods
	sugar := log.Sugar()
	sugar.Infof("Starting crawler with limit=%d", *limit)

	// ─── Build DB connection string ───────────────────────────
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
	)

	// ─── Connect to DB (GORM) ─────────────────────────────────
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		sugar.Fatalf("db connection failed: %v", err)
	}

	// ─── Start crawl ──────────────────────────────────────────
	c := adidas.NewAdidasCrawler(db, sugar)
	products, err := c.CrawlProducts(*limit)
	if err != nil {
		sugar.Fatalf("crawl failed: %v", err)
	}

	//repo := postgres.NewProductRepository(db)
	//if err := repo.BulkInsert(products); err != nil {
	//	sugar.Fatalf("insert failed: %v", err)
	//}

	sugar.Infof("✅ Successfully crawled and stored %d products", len(products))
}
