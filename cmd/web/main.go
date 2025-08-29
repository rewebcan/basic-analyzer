package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

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
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	l := slog.New(jsonHandler)
	f := fetcher.NewFetcher(hc, l, config.BodySizeLimit)
	c := crawler.NewCrawler(f, l, crawler.WithConcurrencyLimit(10))

	crawlCtrl := crawler.NewCrawlController(f, c, l)

	app.HandleFunc("/", crawlCtrl.CrawlHandler)

	fmt.Printf("Server is up and running at localhost:8080...\n")

	chain := util.Chain(
		util.RecoveryMiddleware(),
		util.LoggingMiddleware(),
		util.CORSMiddleware([]string{"*"}),
		util.RequestIDMiddleware(),
	)

	if err := http.ListenAndServe(":8080", chain(app)); err != nil {
		log.Fatal(err)
	}
}
