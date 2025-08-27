package crawler

import (
	"github.com/rewebcan/url-fetcher-home24/internal/fetcher"
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

	for _, a := range result.Anchors {
		urlStr := baseUrl.String()

		if false == a.External {
			rel, _ := url.Parse(a.URL)
			urlStr = baseUrl.ResolveReference(rel).String()
		}

		if _, err := f.Fetch(urlStr); err != nil {
			failedURLs = append(failedURLs, a)
		}
	}

	return &CrawlResult{
		FetchResult: *result,
		FailedURLs:  failedURLs,
	}, nil
}
