package fetcher

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		html := `
			<!DOCTYPE html>
			<html>
				<head>
					<title></title>
				</head>
				<body>
					<h1>Title 1<h2>
					<form action="/" method="post">
						
					</form>
					<a href="https://google.com">Page 1</a>
					<a href="/faq">Page 2</a>
				</body>
			</html>`
		w.Write([]byte(html))
	}))

	defer server.Close()

	f := NewFetcher()

	result, err := f.Fetch("https://crawler-test.com/mobile/separate_desktop_with_different_h1")

	if err != nil {
		t.Errorf("Fetcher returned an error: %s", err)
	}

	assert.NotNil(t, result)
	assert.Nil(t, err)
}
