package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
)

func Timeout(duration time.Duration) func(http.Handler) http.Handler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//create context with timeout
			ctx, cancel := context.WithTimeout(r.Context(), duration)

			defer cancel()

			//replace request context with timeout context

			r = r.WithContext(ctx)

			//channel to signal handler completion

			done := make(chan struct{})

			go func() {
				next.ServeHTTP(w, r)
				close(done)
			}()

			select {
			case <- done:
				//Handler completed successfully
				log.Println("Handler completed successfully")
				return
			case <- ctx.Done():
				// Timeout exceeded
				log.Println("Handler didn't complete")
				http.Error(w, "Request timeout", http.StatusRequestTimeout)
			}
		})
	}
}