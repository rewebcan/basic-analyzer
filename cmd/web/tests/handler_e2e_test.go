package tests

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rewebcan/url-fetcher-home24/internal/crawler"
	"github.com/rewebcan/url-fetcher-home24/internal/fetcher"
	"github.com/rewebcan/url-fetcher-home24/internal/util"
)

func TestCrawlHandler_E2E(t *testing.T) {
	// Setup fake website server
	fakeServer := setupFakeWebsite()
	defer fakeServer.Close()

	// Setup web application server
	app := setupWebApp()

	tests := []struct {
		name           string
		url            string
		expectedStatus int
		checkResponse  func(t *testing.T, resp *http.Response, body string)
	}{
		{
			name:           "valid URL with successful crawl",
			url:            fakeServer.URL + "/test-page",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp *http.Response, body string) {
				if !strings.Contains(body, "Test Page") {
					t.Errorf("Expected response to contain 'Test Page', got: %s", body)
				}
				if !strings.Contains(body, "https://external.com") {
					t.Errorf("Expected response to contain external link, got: %s", body)
				}
				if !strings.Contains(body, "/internal-link") {
					t.Errorf("Expected response to contain internal link, got: %s", body)
				}
				// Check that broken links are correctly marked as errors
				if !strings.Contains(body, "class=\"error\"") {
					t.Errorf("Expected broken links to be marked as errors, got: %s", body)
				}
				// Check that the title is correctly extracted and displayed
				if !strings.Contains(body, "<p>Title: Test Page</p>") {
					t.Errorf("Expected page title to be extracted and displayed, got: %s", body)
				}
				// Check that login form indicator is displayed (should be "No" since test page has no password fields)
				if !strings.Contains(body, "<p>Login Form:  No </p>") {
					t.Errorf("Expected login form indicator to show 'No' for page without login form, got: %s", body)
				}
				// Check that HTML version is correctly detected and displayed
				if !strings.Contains(body, "<p>HTML Version: HTML5</p>") {
					t.Errorf("Expected HTML version to be detected as HTML5, got: %s", body)
				}
			},
		},
		{
			name:           "empty URL validation",
			url:            "",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp *http.Response, body string) {
				if !strings.Contains(body, "url string can not be empty") {
					t.Errorf("Expected validation error for empty URL, got: %s", body)
				}
			},
		},
		{
			name:           "invalid URL scheme",
			url:            "ftp://example.com",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp *http.Response, body string) {
				if !strings.Contains(body, "unsupported URL scheme") {
					t.Errorf("Expected validation error for invalid scheme, got: %s", body)
				}
			},
		},
		{
			name:           "non-existent URL",
			url:            "http://non-existent-domain-that-should-fail-12345.com",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp *http.Response, body string) {
				// Should contain some error message
				if !strings.Contains(body, "error") && !strings.Contains(body, "Error") {
					t.Errorf("Expected error for non-existent URL, got: %s", body)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare form data
			formData := url.Values{}
			if tt.url != "" {
				formData.Set("url", tt.url)
			}

			// Create POST request
			req := httptest.NewRequest(http.MethodPost, "/analyze", strings.NewReader(formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// Create response recorder
			w := httptest.NewRecorder()

			// Execute request
			app.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Read response body
			body := w.Body.String()

			// Run custom response checks
			tt.checkResponse(t, w.Result(), body)
		})
	}
}

// setupFakeWebsite creates a test HTTP server that simulates a website to crawl
func setupFakeWebsite() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/test-page":
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			html := `
<!DOCTYPE html>
<html>
<head>
    <title>Test Page</title>
    <meta name="description" content="Test page for crawling">
</head>
<body>
    <h1>Test Page</h1>
    <p>This is a test page for the crawler.</p>
    <a href="/internal-link">Internal Link</a>
    <a href="https://external.com">External Link</a>
    <a href="/broken-link">Broken Link</a>
</body>
</html>`
			_, _ = w.Write([]byte(html))
		case "/internal-link":
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("<html><body><h1>Internal Page</h1></body></html>"))
		case "/broken-link":
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

// setupWebApp creates the web application with the crawl handler
func setupWebApp() http.Handler {
	// Get the absolute path to the views directory
	originalDir, _ := os.Getwd()
	projectRoot := filepath.Join(originalDir, "..", "..", "..")
	templatePath := filepath.Join(projectRoot, "views", "index.html")

	// Create configuration
	config := util.NewDefaultCrawlerConfig()

	// Create HTTP client
	hc := &http.Client{Timeout: config.CrawlerTimeout}

	// Create logger
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	// Create fetcher
	f := fetcher.NewFetcher(hc, logger, config.BodySizeLimit)

	// Create crawler
	c := crawler.NewCrawler(f, logger, crawler.WithConcurrencyLimit(10))

	// Create crawl controller with custom template path
	crawlCtrl := crawler.NewCrawlControllerWithTemplate(f, c, logger, templatePath)

	// Create router
	app := http.NewServeMux()
	app.HandleFunc("/analyze", crawlCtrl.CrawlHandler)

	return app
}

// TestCrawlHandler_GetRequest tests GET requests to the handler
func TestCrawlHandler_GetRequest(t *testing.T) {
	app := setupWebApp()

	req := httptest.NewRequest(http.MethodGet, "/analyze", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for GET request, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Home24 Basic Crawler") {
		t.Errorf("Expected GET response to contain page title, got: %s", body)
	}
}

// TestCrawlHandler_FormSubmission tests the complete form submission flow
func TestCrawlHandler_FormSubmission(t *testing.T) {
	// Setup fake website
	fakeServer := setupFakeWebsite()
	defer fakeServer.Close()

	app := setupWebApp()

	// Create form data
	formData := url.Values{}
	formData.Set("url", fakeServer.URL+"/test-page")

	// Create POST request
	req := httptest.NewRequest(http.MethodPost, "/analyze", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()

	// Verify the response contains crawled data
	expectedContent := []string{
		"Test Page",            // Title from crawled page
		"https://external.com", // External link
		"/internal-link",       // Internal link
		"Login Form:  No ",     // Login form indicator (should be "No" since test page has no password fields)
	}

	for _, content := range expectedContent {
		if !strings.Contains(body, content) {
			t.Errorf("Expected response to contain '%s', got: %s", content, body)
		}
	}
}
