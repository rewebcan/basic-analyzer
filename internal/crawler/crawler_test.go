package crawler

import (
	"github.com/rewebcan/url-fetcher-home24/internal/fetcher"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCrawler(t *testing.T) {
	r, err := Crawl("https://crawler-test.com/mobile/separate_desktop_with_different_h1", fetcher.NewFakeFetcher())

	assert.NotNil(t, r)
	assert.Nil(t, err)

	assert.NotNil(t, r.FailedURLs)
	assert.Len(t, r.FailedURLs, 2)
}

func TestCrawler_Fail(t *testing.T) {
	r, err := Crawl("https://google.com", fetcher.NewFakeFetcher())

	assert.Nil(t, r)
	assert.NotNil(t, err)
}
