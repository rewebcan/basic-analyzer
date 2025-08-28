package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

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

	app.HandleFunc("/analyze", crawlCtrl.CrawlHandler)

	fmt.Println("Server is up and running...")

	if err := http.ListenAndServe(":8080", loggingMiddleware(app)); err != nil {
		log.Fatal(err)
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(p []byte) (int, error) {
	if r.status == 0 {
		// Default status if WriteHeader was not explicitly called
		r.status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(p)
	r.bytes += n
	return n, err
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w}
		next.ServeHTTP(rec, r)
		duration := time.Since(start)
		log.Printf("%s %s %d %dB %s", r.Method, r.URL.Path, rec.status, rec.bytes, duration)
	})
}
