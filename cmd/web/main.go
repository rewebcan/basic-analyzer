package main

import (
	"log"
	"net/http"

	"github.com/rewebcan/url-fetcher-home24/internal/crawler"
	"github.com/rewebcan/url-fetcher-home24/internal/fetcher"
	"github.com/rewebcan/url-fetcher-home24/internal/util"
)

var (
	config = util.LoadCrawlerConfigFromEnv()
)

func main() {
	app := http.NewServeMux()

	hc := &http.Client{Timeout: config.CrawlerTimeout}
	f := fetcher.NewFetcher(hc, config.BodySizeLimit)
	c := crawler.NewCrawler(f, crawler.WithConcurrencyLimit(10))

	crawlCtrl := crawler.NewCrawlController(f, c)

	app.HandleFunc("/analyze", crawlCtrl.CrawlHandler)

	if err := http.ListenAndServe(":8080", app); err != nil {
		log.Fatal(err)
	}
}
