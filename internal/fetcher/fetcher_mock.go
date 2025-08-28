package fetcher

import (
	"context"
	"errors"
)

type fakeFetcher map[string]*FetchResult

func (f fakeFetcher) Ping(ctx context.Context, url string) (err error) {
	_, err = f.Fetch(ctx, url)
	return
}

func (f fakeFetcher) Fetch(ctx context.Context, url string) (*FetchResult, error) {
	if res, ok := f[url]; ok {
		return res, nil
	}

	return nil, errors.New("not found")
}

func NewFakeFetcher() Fetcher {
	return fakeFetcher{
		"https://crawler-test.com/mobile/separate_desktop_with_different_h1": {
			URL:          "https://crawler-test.com/mobile/separate_desktop_with_different_h1",
			Title:        "Desktop and Mobile pages with different H1 tags",
			HeaderMap:    map[string][]string{"h1": {"desktop", "mobile"}},
			Anchors:      []Anchor{{External: true, URL: "https://google.com"}, {External: false, URL: "/faq"}},
			HasLoginForm: true,
		},
	}
}
