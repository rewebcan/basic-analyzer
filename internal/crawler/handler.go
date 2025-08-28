package crawler

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/rewebcan/url-fetcher-home24/internal/fetcher"
	"github.com/rewebcan/url-fetcher-home24/internal/util"
)

var ErrValidation = errors.New("validation error")

func NewCrawlRequestFromRequest(r *http.Request) *CrawlRequest {
	cr := &CrawlRequest{URL: r.FormValue("url")}

	return cr
}

type crawlController struct {
	f fetcher.Fetcher
	c Crawler
}

func NewCrawlController(f fetcher.Fetcher, c Crawler) *crawlController {
	return &crawlController{f, c}
}

func (ctrl *crawlController) CrawlHandler(w http.ResponseWriter, r *http.Request) {
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

		crawlResult, err = ctrl.c.Crawl(cr.URL)

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
		return fmt.Errorf("%w: url string can not be empty", ErrValidation)
	}

	normalizedURL, err := util.NormalizeURL(cr.URL)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrValidation, err.Error())
	}
	cr.URL = normalizedURL

	if err := util.IsValidHTTPURL(cr.URL); err != nil {
		return fmt.Errorf("%w: %s", ErrValidation, err.Error())
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
