package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/LordCodex164/httpserver/internal/metrics"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//current date
		start := time.Now()

		// Get request ID from context
		requestID := GetRequestID(r.Context())

		wrapped := NewResponseWriter(w)

		//call the next handler
		next.ServeHTTP(wrapped, r)

		log.Println(">>>r", r.Response)

		//calculate latency
		latency := time.Since(start)

		metrics.GetInstance().RecordRequest(
			wrapped.StatusCode(),
			latency,
			r.URL.Path,
		)

		//log after handler 
		log.Printf(
			"[%s] %s %s %s %d %v",
			requestID,
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			wrapped.StatusCode(),
			time.Since(start),
		)
	})
}