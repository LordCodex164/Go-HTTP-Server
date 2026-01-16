package middleware

import (
	"context"
	"net/http"
	"github.com/google/uuid"
)

type contextKey string

const reqIDKey contextKey = "requestID"

// RequestID middleware adds a unique ID to each request
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate unique ID
		requestID := uuid.New().String()

		// Add to response header (for client)
		w.Header().Set("X-Request-ID", requestID)

		// Add to context (for use in handlers)
		ctx := context.WithValue(r.Context(), reqIDKey, requestID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// GetRequestID extracts request ID from context
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(reqIDKey).(string); ok {
		return id
	}
	return ""
}