package crawler

import (
	"errors"
	"fmt"
	"github.com/rewebcan/url-fetcher-home24/internal/fetcher"
	"html/template"
	"net/http"
)

var ValidationErr = errors.New("validation error")

func NewCrawlRequestFromRequest(r *http.Request) *CrawlRequest {
	cr := &CrawlRequest{URL: r.FormValue("url")}

	return cr
}

func CrawlHandler(w http.ResponseWriter, r *http.Request) {
	var crawlResult *CrawlResult

	t := template.Must(template.ParseFiles("views/index.html"))

	if r.Method == "POST" {
		cr := NewCrawlRequestFromRequest(r)

		fmt.Printf("Request received: %s\n", cr)

		if err := cr.Validate(); err != nil {
			_ = t.Execute(w, CrawlPageResponse{
				CrawlResult: nil,
				Errors:      []string{err.Error()},
			})
			return
		}

		var err error

		crawlResult, err = Crawl(cr.URL, fetcher.NewFetcher(http.DefaultClient))

		if err != nil {
			_ = t.Execute(w, CrawlPageResponse{
				CrawlResult: nil,
				Errors:      []string{err.Error()},
			})
			return
		}
	}

	_ = t.Execute(w, CrawlPageResponse{
		CrawlResult: crawlResult,
		Errors:      nil,
	})
}

type CrawlPageResponse struct {
	CrawlResult *CrawlResult
	Errors      []string
}

type CrawlRequest struct {
	URL string
}

func (cr *CrawlRequest) Validate() error {
	if cr.URL == "" {
		return fmt.Errorf("%w: url string can not be empty", ValidationErr)
	}

	return nil
}
