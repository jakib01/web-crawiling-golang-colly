package main

import (
	"flag"
	"fmt"
	"github.com/jakib01/web-crawiling-golang-colly/internal/crawler/adidas"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/jakib01/web-crawiling-golang-colly/internal/config"
	"github.com/jakib01/web-crawiling-golang-colly/internal/logger"
	//"github.com/jakib01/web-crawiling-golang-colly/internal/repository/postgres"
)

func main() {
	envFile := flag.String("env", ".env", "path to env file")
	limit := flag.Int("limit", 300, "max number of products to crawl")
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
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	// ─── Connect to DB ─────────────────────────────────────────
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		sugar.Fatalf("db connection failed: %v", err)
	}
	defer db.Close()

	// ─── Connect to DB ─────────────────────────────────────────
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
