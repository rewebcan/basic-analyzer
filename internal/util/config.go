package util

import (
	"os"
	"strconv"
	"time"
)

// CrawlerConfig holds configuration for the crawler and fetcher
type CrawlerConfig struct {
	CrawlerTimeout   time.Duration
	BodySizeLimit    int64
	ConcurrencyLimit int
	CrawlDepth       int
}

// NewDefaultCrawlerConfig creates a default configuration
func NewDefaultCrawlerConfig() *CrawlerConfig {
	return &CrawlerConfig{
		CrawlerTimeout:   10 * time.Second,
		BodySizeLimit:    10 << 20, // 10MB
		ConcurrencyLimit: 10,
		CrawlDepth:       1,
	}
}

// LoadCrawlerConfigFromEnv loads configuration from environment variables
func LoadCrawlerConfigFromEnv() *CrawlerConfig {
	config := NewDefaultCrawlerConfig()

	if timeoutStr := os.Getenv("CRAWLER_TIMEOUT"); timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil {
			config.CrawlerTimeout = timeout
		}
	}

	if bodySizeStr := os.Getenv("CRAWLER_BODY_SIZE_LIMIT"); bodySizeStr != "" {
		if bodySize, err := strconv.ParseInt(bodySizeStr, 10, 64); err == nil {
			config.BodySizeLimit = bodySize
		}
	}

	if concurrencyStr := os.Getenv("CRAWLER_CONCURRENCY_LIMIT"); concurrencyStr != "" {
		if concurrency, err := strconv.Atoi(concurrencyStr); err == nil && concurrency > 0 {
			config.ConcurrencyLimit = concurrency
		}
	}

	return config
}
