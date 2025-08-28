package fetcher

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"golang.org/x/net/html"
)

type Fetcher interface {
	Fetch(url string) (*FetchResult, error)
	Ping(url string) error
}

type fetcher struct {
	httpClient    *http.Client
	bodySizeLimit int64
	logger        *slog.Logger
}

func (f fetcher) Ping(url string) error {
	resp, err := f.httpClient.Get(url)

	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.New("unexpected status code: " + resp.Status)
	}

	return nil
}

func NewFetcher(httpClient *http.Client, logger *slog.Logger, bodySizeLimit int64) Fetcher {
	return fetcher{httpClient: httpClient, logger: logger, bodySizeLimit: bodySizeLimit}
}

// Fetch
// Fetches the given url and returns a structured response
// if error returned it might be the reason the given url is not reachable
func (f fetcher) Fetch(url string) (*FetchResult, error) {
	f.logger.Info("Starting fetch", "url", url)

	resp, err := fetch(f.httpClient, url, f.bodySizeLimit)
	if err != nil {
		f.logger.Error("Failed to fetch URL", "url", url, "error", err.Error())
		return nil, err
	}

	defer resp.Close()

	f.logger.Info("Successfully fetched URL", "url", url)

	r := &FetchResult{}
	r.URL = url
	r.HeaderMap = make(map[string][]string)

	anchorMap := map[string]struct{}{}

	err = streamToken(resp, func(z *html.Tokenizer, tt html.TokenType, tok html.Token) error {
		if tt == html.ErrorToken {
			if z.Err() == io.EOF {
				return nil
			}

			return z.Err()
		}

		tag := tok.Data

		switch tag {
		case "a":
			a, ok := extractAnchor(tok)

			if !ok {
				break
			}

			if _, ok := anchorMap[a.URL]; !ok {
				anchorMap[a.URL] = struct{}{}
				r.Anchors = append(r.Anchors, a)
			}
		case "h1", "h2", "h3", "h4", "h5", "h6":
			extractHeaders(z, tok, r.HeaderMap)
		case "title":
			r.Title, _ = readTextValue(z)
		case "input":
			if attr, ok := findAttr(tok, "type"); ok && attr.Val == "password" {
				r.HasLoginForm = true
			}
		}

		return nil
	})

	if err != nil {
		f.logger.Error("Failed to parse HTML content", "url", url, "error", err.Error())
		return nil, err
	}

	f.logger.Info("Fetch completed successfully", "url", url, "title", r.Title, "anchors_found", len(r.Anchors), "has_login_form", r.HasLoginForm)

	return r, nil
}

type FormElement struct {
	Name string
	Type string
}

type Anchor struct {
	External bool
	URL      string
}

type Form struct {
	Elements []FormElement
	Action   string
	Method   string
}

type FetchResult struct {
	URL          string
	Title        string
	HeaderMap    map[string][]string
	Anchors      []Anchor
	HasLoginForm bool
}
