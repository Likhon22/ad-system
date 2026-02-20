package middleware

import (
	"net/http"
	"time"

	"github.com/likhon22/ad-system/auth-service/internal/utils"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {

	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)

}
func (mw *Middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		attrs := []any{
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.statusCode,
			"duration_ms", time.Since(start).Milliseconds(),
			"request_id", utils.RequestIDFromContext(r.Context()),
		}
		switch {
		case rw.statusCode >= 500:
			mw.logger.Error("request", attrs...)
		case rw.statusCode >= 400:
			mw.logger.Warn("request", attrs...)
		default:
			mw.logger.Info("request", attrs...)
		}
	})

}
