package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logging(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &responseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(rw, r)

			duration := time.Since(start)
			reqID, _ := r.Context().Value(RequestIDKey).(string)

			logger.Info("http request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", rw.status),
				zap.Duration("latency", duration),
				zap.String("request_id", reqID),
			)
		})
	}
}