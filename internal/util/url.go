package util

import (
	"fmt"
	"net/url"
	"strings"
)

// NormalizeURL normalizes a URL by:
// - Converting the scheme and host to lowercase.
// - Removing default ports (80 for http, 443 for https).
// - Removing trailing slashes (unless it's the root path).
// - Sorting query parameters (not implemented for simplicity in this example).
func NormalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	// Convert scheme and host to lowercase
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)

	// Remove default ports
	if (u.Scheme == "http" && u.Port() == "80") || (u.Scheme == "https" && u.Port() == "443") {
		u.Host = strings.TrimSuffix(u.Host, ":"+u.Port())
	}

	// Remove trailing slash unless it's the root path
	if u.Path != "/" {
		u.Path = strings.TrimSuffix(u.Path, "/")
	}

	// For simplicity, query parameters are not sorted here, but could be added.

	return u.String(), nil
}

// IsValidHTTPURL checks if a URL has a valid HTTP or HTTPS scheme.
func IsValidHTTPURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("unsupported URL scheme: %s, only http and https are allowed", u.Scheme)
	}

	if u.Host == "" {
		return fmt.Errorf("URL is missing a host")
	}

	return nil
}
