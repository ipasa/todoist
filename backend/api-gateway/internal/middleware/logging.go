package middleware

import (
	"net/http"
	"time"

	"github.com/todoist/backend/pkg/logger"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

// Logging middleware logs HTTP requests
func Logging(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rw, r)

			duration := time.Since(start)

			log.WithFields(map[string]interface{}{
				"method":   r.Method,
				"path":     r.URL.Path,
				"status":   rw.status,
				"duration": duration.Milliseconds(),
				"ip":       r.RemoteAddr,
			}).Info("HTTP request")
		})
	}
}
