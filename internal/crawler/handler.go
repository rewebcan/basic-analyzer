package crawler

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/rewebcan/url-fetcher-home24/internal/fetcher"
	"github.com/rewebcan/url-fetcher-home24/internal/util"
)

var ValidationErr = errors.New("validation error")

func NewCrawlRequestFromRequest(r *http.Request) *CrawlRequest {
	cr := &CrawlRequest{URL: r.FormValue("url")}

	return cr
}

func CrawlHandler(w http.ResponseWriter, r *http.Request) {
	var crawlResult *CrawlResult

	t := template.Must(
		template.New("index.html").
			Funcs(funcMap).
			ParseFiles("views/index.html"),
	)

	if r.Method == "POST" {
		cr := NewCrawlRequestFromRequest(r)

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

	normalizedURL, err := util.NormalizeURL(cr.URL)
	if err != nil {
		return fmt.Errorf("%w: %s", ValidationErr, err.Error())
	}
	cr.URL = normalizedURL

	if err := util.IsValidHTTPURL(cr.URL); err != nil {
		return fmt.Errorf("%w: %s", ValidationErr, err.Error())
	}

	return nil
}

var funcMap = template.FuncMap{
	"containsUrl": func(anchors []fetcher.Anchor, str string) bool {
		for _, a := range anchors {
			if a.URL == str {
				return true
			}
		}

		return false
	},
}
