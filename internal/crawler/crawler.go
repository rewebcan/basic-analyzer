package crawler

import (
	"context"
	"net/url"

	"github.com/rewebcan/url-fetcher-home24/internal/fetcher"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type CrawlResult struct {
	fetcher.FetchResult
	FailedURLs []fetcher.Anchor
}

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

// Crawl
// Crawls a page with given fetcher and URL, applying crawl options using Commander pattern
func Crawl(f fetcher.Fetcher, urlRaw string, opts ...CrawlOption) (*CrawlResult, error) {
	// Default configuration
	config := &crawlConfig{
		concurrencyLimit: 10,
	}

	// Apply options
	for _, opt := range opts {
		opt(config)
	}

	ctx := context.Background()
	sem := semaphore.NewWeighted(int64(config.concurrencyLimit))

	baseUrl, err := url.Parse(urlRaw)
	if err != nil {
		return nil, err
	}

	result, err := f.Fetch(urlRaw)
	if err != nil {
		return nil, err
	}
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
			if err := f.Ping(urlStr); err != nil {
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

	return &CrawlResult{
		FetchResult: *result,
		FailedURLs:  failedURLs,
	}, nil
}
