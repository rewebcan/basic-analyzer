package util

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Middleware represents a function that wraps an http.Handler
type Middleware func(http.Handler) http.Handler

// Chain applies multiple middleware in order (left to right)
// The first middleware in the slice will be the outermost layer
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		// Apply middlewares in reverse order so the first one becomes outermost
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// Apply applies a single middleware to a handler
func Apply(middleware Middleware, handler http.Handler) http.Handler {
	return middleware(handler)
}

// ApplyChain applies a chain of middleware to a handler
func ApplyChain(handler http.Handler, middlewares ...Middleware) http.Handler {
	return Chain(middlewares...)(handler)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(p []byte) (int, error) {
	if r.status == 0 {
		// Default status if WriteHeader was not explicitly called
		r.status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(p)
	r.bytes += n
	return n, err
}

// LoggingMiddleware creates a middleware that logs HTTP requests
func LoggingMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w}
			next.ServeHTTP(rec, r)
			duration := time.Since(start)
			log.Printf("%s %s %d %dB %s", r.Method, r.URL.Path, rec.status, rec.bytes, duration)
		})
	}
}

// CORSMiddleware creates a middleware that handles CORS headers
func CORSMiddleware(allowedOrigins []string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			for _, allowedOrigin := range allowedOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RecoveryMiddleware creates a middleware that recovers from panics
func RecoveryMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("Panic recovered: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// RequestIDMiddleware creates a middleware that adds a request ID to each request
func RequestIDMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = generateRequestID()
			}
			w.Header().Set("X-Request-ID", requestID)
			next.ServeHTTP(w, r)
		})
	}
}

// TimeoutMiddleware creates a middleware that adds a timeout to requests
func TimeoutMiddleware(timeout time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
