package fetcher

import "errors"

var PageNotFound = errors.New("page not found")

type Fetcher interface {
	Fetch(url string) (*FetchResult, error)
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
	URL       string
	Title     string
	HeaderMap map[string][]string
	Anchors   []Anchor
	Forms     []Form
}

type fakeFetcher map[string]*FetchResult

// Fetch
// Fetches the given url and returns a structured response
// if error returned it might be the reason the given url is not reachable
func (f fakeFetcher) Fetch(url string) (*FetchResult, error) {
	if res, ok := f[url]; ok {
		return res, nil
	}

	return nil, PageNotFound
}

var fakeForm1 = []Form{{
	Elements: []FormElement{
		{Name: "username", Type: "text"},
		{Name: "password", Type: "password"},
	},
	Method: "POST", Action: "/auth/login"},
}

func NewFetcher() Fetcher {
	return fakeFetcher{
		"https://crawler-test.com/mobile/separate_desktop_with_different_h1": {
			URL:       "https://crawler-test.com/mobile/separate_desktop_with_different_h1",
			Title:     "Desktop and Mobile pages with different H1 tags",
			HeaderMap: map[string][]string{"h1": {"desktop", "mobile"}},
			Anchors:   []Anchor{{External: true, URL: "https://google.com"}, {External: false, URL: "/faq"}},
			Forms:     fakeForm1,
		},
	}
}
