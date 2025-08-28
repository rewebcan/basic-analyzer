package fetcher

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
		_, _ = w.Write([]byte(html))
	}))

	defer server.Close()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	f := NewFetcher(server.Client(), logger, 10<<20)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := f.Fetch(ctx, server.URL)

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
