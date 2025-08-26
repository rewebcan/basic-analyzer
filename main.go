package main

import (
	"encoding/json"
	"fmt"
	"github.com/rewebcan/url-fetcher-home24/crawler"
	"github.com/rewebcan/url-fetcher-home24/fetcher"
	"log"
	"net/http"
)

var url = "https://crawler-test.com/mobile/separate_desktop_with_different_h1"

func main() {
	r, err := crawler.Crawl(url, fetcher.NewFetcher(http.DefaultClient))

	if err != nil {
		log.Fatal(err)
	}

	str, _ := json.MarshalIndent(r, "", "\t")

	fmt.Println(string(str))
}
