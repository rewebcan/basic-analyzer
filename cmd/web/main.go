package main

import (
	"github.com/rewebcan/url-fetcher-home24/internal/crawler"
	"log"
	"net/http"
)

var url = "https://crawler-test.com/mobile/separate_desktop_with_different_h1"

func main() {
	app := http.NewServeMux()

	app.HandleFunc("/analyze", crawler.CrawlHandler)

	if err := http.ListenAndServe(":8080", app); err != nil {
		log.Fatal(err)
	}
}
