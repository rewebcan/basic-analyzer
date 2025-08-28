package crawler

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/rewebcan/url-fetcher-home24/internal/fetcher"
	"github.com/rewebcan/url-fetcher-home24/internal/util"
)

var ErrValidation = errors.New("validation error")

func NewCrawlRequestFromRequest(r *http.Request) *CrawlRequest {
	cr := &CrawlRequest{URL: r.FormValue("url")}

	return cr
}

type crawlController struct {
	f            fetcher.Fetcher
	c            Crawler
	logger       *slog.Logger
	templatePath string
}

func NewCrawlController(f fetcher.Fetcher, c Crawler, l *slog.Logger) *crawlController {
	return &crawlController{f: f, c: c, logger: l, templatePath: "views/index.html"}
}

func NewCrawlControllerWithTemplate(f fetcher.Fetcher, c Crawler, l *slog.Logger, templatePath string) *crawlController {
	return &crawlController{f: f, c: c, logger: l, templatePath: templatePath}
}

func (ctrl *crawlController) CrawlHandler(w http.ResponseWriter, r *http.Request) {
	var crawlResult *CrawlResult

	ctrl.logger.Info("CrawlHandler started", "method", r.Method, "url", r.URL.Path)

	t := template.Must(
		template.New("index.html").
			Funcs(funcMap).
			ParseFiles(ctrl.templatePath),
	)

	if r.Method == "POST" {
		cr := NewCrawlRequestFromRequest(r)

		if err := cr.Validate(); err != nil {
			ctrl.logger.Warn("Crawl request validation failed", "error", err.Error(), "url", cr.URL)
			_ = t.Execute(w, CrawlPageResponse{
				CrawlResult: nil,
				Errors:      []string{err.Error()},
			})
			return
		}

		ctrl.logger.Info("Starting crawl", "url", cr.URL)

		var err error

		// Create context with timeout to prevent long-running operations
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		crawlResult, err = ctrl.c.Crawl(ctx, cr.URL)

		if err != nil {
			ctrl.logger.Error("Crawl failed", "error", err.Error(), "url", cr.URL)
			_ = t.Execute(w, CrawlPageResponse{
				CrawlResult: nil,
				Errors:      []string{err.Error()},
			})
			return
		}

		ctrl.logger.Info("Crawl completed successfully", "url", cr.URL, "anchors_found", len(crawlResult.Anchors), "failed_urls", len(crawlResult.FailedURLs))
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
