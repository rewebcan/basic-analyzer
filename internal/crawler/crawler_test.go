package crawler

import (
	"io"
	"log/slog"
	"testing"

	"github.com/rewebcan/url-fetcher-home24/internal/fetcher"
	"github.com/stretchr/testify/assert"
)

func TestCrawler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	c := NewCrawler(fetcher.NewFakeFetcher(), logger)
	r, err := c.Crawl("https://crawler-test.com/mobile/separate_desktop_with_different_h1")

	assert.NotNil(t, r)
	assert.Nil(t, err)

	assert.NotNil(t, r.FailedURLs)
	assert.Len(t, r.FailedURLs, 2)
}

func TestCrawler_Fail(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	c := NewCrawler(fetcher.NewFakeFetcher(), logger)
	r, err := c.Crawl("https://google.com")

	assert.Nil(t, r)
	assert.NotNil(t, err)
}
