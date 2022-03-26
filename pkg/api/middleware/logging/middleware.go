package logging

import (
	"log"
	"net/http"
	"time"
)

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusResponseWriter) WriteHeader(status int) {
	r.statusCode = status
	r.ResponseWriter.WriteHeader(status)
}

// NewStatusResponseWriter returns pointer to a new statusResponseWriter object
func NewStatusResponseWriter(w http.ResponseWriter) *statusResponseWriter {
	return &statusResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}
func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		sw := NewStatusResponseWriter(w)

		defer func() {
			log.Printf(
				"[%s] [%v] [%d] %s %s %s",
				req.Method,
				time.Since(start),
				sw.statusCode,
				req.Host,
				req.URL.Path,
				req.URL.RawQuery,
			)
		}()

		next.ServeHTTP(sw, req)
	})
}
