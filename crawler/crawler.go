package crawler

import (
	"fmt"
	"github.com/rewebcan/url-fetcher-home24/fetcher"
	"net/url"
	"strings"
)

type CrawlResult struct {
	fetcher.FetchResult
	FailedURLs []fetcher.Anchor
}

// Crawl
// Crawls a page with given Fetcher and Crawls again only `one more` time for the Anchors returned
func Crawl(urlRaw string, f fetcher.Fetcher) (*CrawlResult, error) {
	result, err := f.Fetch(urlRaw)

	if err != nil {
		return nil, err
	}

	var failedURLs []fetcher.Anchor

	for _, a := range result.Anchors {
		urlStr := a.URL
		if false == a.External {
			urlStr, _ = url.JoinPath(urlRaw, a.URL)
		}

		urlStr = strings.TrimRight(urlStr, "/")

		if _, err := f.Fetch(urlStr); err != nil {
			fmt.Printf("Failed to fetch URL: %s\n", urlStr)
			failedURLs = append(failedURLs, a)
		}
	}

	return &CrawlResult{
		FetchResult: *result,
		FailedURLs:  failedURLs,
	}, nil
}
