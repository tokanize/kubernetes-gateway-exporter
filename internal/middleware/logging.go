package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// LoggingMiddleware logs incoming HTTP requests using slog.
func LoggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Use a custom ResponseWriter to capture the status code
		ww := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(ww, r)

		logger.Info("HTTP Request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.Int("status", ww.statusCode),
			slog.Duration("duration", time.Since(start)),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
