package crawler

import (
	"context"
	"log/slog"
	"net/url"

	"github.com/rewebcan/url-fetcher-home24/internal/fetcher"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

// CrawlOption is a function that configures crawl behavior
type CrawlOption func(*crawlConfig)

// crawlConfig holds internal crawl configuration
type crawlConfig struct {
	concurrencyLimit int
}

// WithConcurrencyLimit sets the maximum number of concurrent link pings
func WithConcurrencyLimit(limit int) CrawlOption {
	return func(c *crawlConfig) {
		c.concurrencyLimit = limit
	}
}

type Crawler interface {
	Crawl(ctx context.Context, url string) (*CrawlResult, error)
}

type crawler struct {
	f           fetcher.Fetcher
	crawlConfig *crawlConfig
	logger      *slog.Logger
}

// NewCrawler creates a new Crawler using the provided fetcher and optional
// configuration options.
func NewCrawler(f fetcher.Fetcher, logger *slog.Logger, opts ...CrawlOption) Crawler {
	// Default configuration
	config := &crawlConfig{
		concurrencyLimit: 10,
	}

	// Apply options
	for _, opt := range opts {
		opt(config)
	}

	return &crawler{f: f, crawlConfig: config, logger: logger}
}

type CrawlResult struct {
	fetcher.FetchResult
	FailedURLs []fetcher.Anchor
}

// Crawl
// Crawls a page with given URL
func (c *crawler) Crawl(ctx context.Context, urlRaw string) (*CrawlResult, error) {
	c.logger.Info("Starting crawl", "url", urlRaw, "concurrency_limit", c.crawlConfig.concurrencyLimit)

	sem := semaphore.NewWeighted(int64(c.crawlConfig.concurrencyLimit))

	baseUrl, err := url.Parse(urlRaw)
	if err != nil {
		c.logger.Error("Failed to parse URL", "url", urlRaw, "error", err.Error())
		return nil, err
	}

	c.logger.Info("Fetching main page", "url", urlRaw)
	result, err := c.f.Fetch(ctx, urlRaw)
	if err != nil {
		c.logger.Error("Failed to fetch main page", "url", urlRaw, "error", err.Error())
		return nil, err
	}

	c.logger.Info("Main page fetched successfully", "url", urlRaw, "anchors_found", len(result.Anchors))

	var failedURLs []fetcher.Anchor

	failedUrlsC := make(chan fetcher.Anchor, 16)
	g := new(errgroup.Group)

	for _, a := range result.Anchors {
		a := a
		g.Go(func() error {
			if err := sem.Acquire(ctx, 1); err != nil {
				return err
			}
			defer sem.Release(1)

			urlStr := a.URL
			if !a.External {
				if rel, err := url.Parse(a.URL); err == nil {
					urlStr = baseUrl.ResolveReference(rel).String()
				} else {
					failedUrlsC <- a
					return nil
				}
			}
			if err := c.f.Ping(ctx, urlStr); err != nil {
				failedUrlsC <- a
			}
			return nil
		})
	}

	go func() {
		_ = g.Wait()
		close(failedUrlsC)
	}()

	for a := range failedUrlsC {
		failedURLs = append(failedURLs, a)
	}

	c.logger.Info("Crawl completed", "url", urlRaw, "total_anchors", len(result.Anchors), "failed_urls", len(failedURLs))

	return &CrawlResult{
		FetchResult: *result,
		FailedURLs:  failedURLs,
	}, nil
}
