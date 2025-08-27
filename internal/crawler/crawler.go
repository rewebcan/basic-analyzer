package crawler

import (
	"github.com/rewebcan/url-fetcher-home24/internal/fetcher"
	"golang.org/x/sync/errgroup"
	"net/url"
)

type CrawlResult struct {
	fetcher.FetchResult
	FailedURLs []fetcher.Anchor
}

// Crawl
// Crawls a page with given Fetcher and Crawls again only `one more` time for the Anchors returned
func Crawl(urlRaw string, f fetcher.Fetcher) (*CrawlResult, error) {
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
