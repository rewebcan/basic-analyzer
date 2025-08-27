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
					<title>Test 1</title>
				</head>
				<body>
					<h1>Title 1</h1>
					<form action="/" method="post">
						<input type="text" name="username" />
						<input type="password" name="password" />
					</form>
					<a href="https://google.com">Page 1</a>
					<a href="/faq">Page 2</a>
				</body>
			</html>`
		w.Write([]byte(html))
	}))

	defer server.Close()

	f := NewFetcher(server.Client())

	result, err := f.Fetch(server.URL)

	if err != nil {
		t.Errorf("Fetcher returned an error: %s", err)
	}

	assert.NotNil(t, result)
	assert.Nil(t, err)

	assert.Equal(t, "Test 1", result.Title)
	assert.Equal(t, true, result.Anchors[0].External)
	assert.Equal(t, false, result.Anchors[1].External)
	assert.Equal(t, []string{"Title 1"}, result.HeaderMap["h1"])
	assert.Equal(t, true, result.HasLoginForm)
}
