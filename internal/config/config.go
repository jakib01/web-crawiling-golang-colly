package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// CrawlerConfig holds settings specific to the crawler
type CrawlerConfig struct {
	StartURL    string
	Concurrency int
}

// Config is the application configuration
type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	Crawler    CrawlerConfig
	LogLevel   string
}

// Load reads configuration from the given env file (e.g. ".env") and environment variables.
func Load(envFile string) (*Config, error) {
	viper.SetConfigFile(envFile)
	viper.AutomaticEnv()

	// Default values
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("CRAWLER_START_URL", "https://shop.adidas.jp/men/")
	viper.SetDefault("CRAWLER_CONCURRENCY", 4)
	viper.SetDefault("LOG_LEVEL", "info")

	// Read from file (if present)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// missing file is fine; we'll use env vars/defaults
	}

	cfg := &Config{
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetInt("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		DBSSLMode:  viper.GetString("DB_SSLMODE"),
		Crawler: CrawlerConfig{
			StartURL:    viper.GetString("CRAWLER_START_URL"),
			Concurrency: viper.GetInt("CRAWLER_CONCURRENCY"),
		},
		LogLevel: viper.GetString("LOG_LEVEL"),
	}

	// basic validation
	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("missing one or more required DB credentials")
	}

	return cfg, nil
}
