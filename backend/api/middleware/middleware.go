package middleware

import (
	"log"
	"net/http"
	"time"

	"warehouse-robots/backend/config"
)

// Middleware represents a function that wraps an http.Handler
type Middleware func(http.Handler) http.Handler

// loggingResponseWriter wraps http.ResponseWriter to capture status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// Chain applies middlewares to a handler in the order they are provided
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// JSONMiddleware sets the Content-Type header to application/json for all responses
func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware handles Cross-Origin Resource Sharing based on configuration.
// We need to setup here especially when my local frontend connects to backend
func CORSMiddleware(cfg *config.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", cfg.CORS.AllowedOrigins)
			w.Header().Set("Access-Control-Allow-Methods", cfg.CORS.AllowedMethods)
			w.Header().Set("Access-Control-Allow-Headers", cfg.CORS.AllowedHeaders)

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// LoggingMiddleware logs HTTP requests with timing information. For performance.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap ResponseWriter to capture status code
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		log.Printf("Started %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		log.Printf("Completed %s %s %d in %v", r.Method, r.URL.Path, lrw.statusCode, duration)
	})
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
