package fetcher

import (
	"errors"
	"golang.org/x/net/html"
	"io"
	"net/http"
)

var PageNotFound = errors.New("page not found")

type Fetcher interface {
	Fetch(url string) (*FetchResult, error)
}

type fetcher struct {
	httpClient *http.Client
}

func NewFetcher(httpClient *http.Client) Fetcher {
	return fetcher{httpClient}
}

// Fetch
// Fetches the given url and returns a structured response
// if error returned it might be the reason the given url is not reachable
func (f fetcher) Fetch(url string) (*FetchResult, error) {
	resp, err := fetch(f.httpClient, url)
	if err != nil {
		return nil, err
	}

	defer resp.Close()

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

			break
		case "h1", "h2", "h3", "h4", "h5", "h6":
			extractHeaders(z, tok, r.HeaderMap)
			break
		case "title":
			r.Title, _ = readTextValue(z)
			break
		case "input":
			if attr, ok := findAttr(tok, "type"); ok && attr.Val == "password" {
				r.HasLoginForm = true
			}
			break
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

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
