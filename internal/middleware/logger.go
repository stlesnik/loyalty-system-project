package middleware

import (
	"github.com/stlesnik/loyalty-system-project/internal/utils"
	"net/http"
	"time"
)

type (
	loggingResponseWriter struct {
		http.ResponseWriter
		status int
		size   int
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.status = statusCode
}

func WithLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.Log.Infow("-----Got request",
			"uri", r.RequestURI,
			"method", r.Method,
		)
		start := time.Now()

		lw := loggingResponseWriter{
			ResponseWriter: w,
		}
		next(&lw, r)

		duration := time.Since(start)

		utils.Log.Infow("-----Sent response",
			"status", lw.status,
			"size", lw.size,
			"duration", duration,
		)
	}
}
