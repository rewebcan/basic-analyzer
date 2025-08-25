package crawler

import "github.com/rewebcan/url-fetcher-home24/fetcher"

type CrawlResult struct {
	fetcher.FetchResult
	FailedURLs []fetcher.Anchor
}

// Crawl
// Crawls a page with given Fetcher and Crawls again only `one more` time for the Anchors returned
func Crawl(url string, f fetcher.Fetcher) (*CrawlResult, error) {
	result, err := f.Fetch(url)

	if err != nil {
		return nil, err
	}

	var failedURLs []fetcher.Anchor

	for _, a := range result.Anchors {
		if _, err := f.Fetch(a.URL); err != nil {
			failedURLs = append(failedURLs, a)
		}
	}

	return &CrawlResult{
		FetchResult: *result,
		FailedURLs:  failedURLs,
	}, nil
}
