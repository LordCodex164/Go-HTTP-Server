package middleware

import (
	"net/http"
	"sync"
	"time"
	"golang.org/x/time/rate"
)

//implement rate limiting
//ratelimiter holds rate limiters for each IP

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

func NewRateLimiter (r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:      r,
		burst:     b,
	}
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	//if the limiter does not exist for the ip address, create a new limiter
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[ip] = limiter
		return limiter
	}

	return limiter
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get client IP
		ip := r.RemoteAddr

		// Get limiter for this IP
		limiter := rl.getLimiter(ip)

		// Check if request is allowed
		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) Cleanup() {
	ticker := time.NewTicker(5 * time.Minute)

	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		defer rl.mu.Unlock()

	}

}