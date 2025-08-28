package fetcher

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

var ErrBadStatus = errors.New("bad status")

type limitedBody struct {
	io.Reader
	closer io.Closer
}

func (l limitedBody) Close() error {
	return l.closer.Close()
}

func fetch(httpClient *http.Client, rawUrl string, bodySizeLimit int64) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, rawUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		io.Copy(io.Discard, io.LimitReader(resp.Body, 8<<10))
		resp.Body.Close()

		return nil, ErrBadStatus
	}

	return limitedBody{
		Reader: io.LimitReader(resp.Body, bodySizeLimit),
		closer: resp.Body,
	}, nil
}

func streamToken(reader io.Reader, handler func(z *html.Tokenizer, tokenType html.TokenType, tok html.Token) error) error {
	z := html.NewTokenizer(reader)

	for {
		tt := z.Next()
		zn := z.Token()

		if tt == html.ErrorToken {
			if z.Err() == io.EOF {
				return nil
			}
		}

		switch tt {
		case html.StartTagToken, html.SelfClosingTagToken:
			if err := handler(z, tt, zn); err != nil {
				return err
			}

			break
		}
	}
}

func extractAnchor(tok html.Token) (Anchor, bool) {
	attr, ok := findAttr(tok, "href")

	if !ok {
		return Anchor{}, false
	}

	if attr.Val == "#" || attr.Val == "javascript::" || attr.Val == "/" {
		return Anchor{}, false
	}

	urlStr := attr.Val
	isUrlInternal := strings.HasPrefix(urlStr, "http://") || strings.HasPrefix(urlStr, "https://")

	return Anchor{
		URL:      urlStr,
		External: isUrlInternal,
	}, true
}

func extractHeaders(z *html.Tokenizer, tok html.Token, hm map[string][]string) {
	tv, _ := readTextValue(z)

	titles, ok := hm[tok.Data]

	if !ok {
		hm[tok.Data] = []string{tv}
	}

	hm[tok.Data] = append(titles, tv)
}

func findAttr(token html.Token, key string) (html.Attribute, bool) {
	for _, attr := range token.Attr {
		if attr.Key == key {
			return attr, true
		}
	}

	return html.Attribute{}, false
}

func readTextValue(z *html.Tokenizer) (string, error) {
	var (
		depth = 1
		buf   bytes.Buffer
	)

	for {
		switch z.Next() {
		case html.ErrorToken:
			return strings.TrimSpace(buf.String()), z.Err()

		case html.TextToken:
			buf.Write(z.Text())

		case html.StartTagToken, html.SelfClosingTagToken:
			depth++

			t := z.Token()
			if t.Data == "br" {
				buf.WriteByte(' ')
			}

		case html.EndTagToken:
			depth--

			if depth == 0 {
				return collapseWS(buf.String()), nil
			}
		}
	}
}

func collapseWS(s string) string {
	s = strings.TrimSpace(s)
	return strings.Join(strings.Fields(s), " ")
}
