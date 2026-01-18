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
				defer close(done)

				// RECOVERY MUST BE HERE
				defer func() {
					if err := recover(); err != nil {
						log.Printf("PANIC in handler: %v", err)
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					}
				}()

				next.ServeHTTP(w, r.WithContext(ctx))
			}()

			select {
			case <-done:
				return
			case <-ctx.Done():
				http.Error(w, "Request Timeout", http.StatusGatewayTimeout)
				return
			}
		})
	}
}