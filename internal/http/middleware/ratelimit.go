package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	burstMultiplier = 2
)

// RateLimiter implements per-IP rate limiting
type RateLimiter struct {
	visitors map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     int
	burst    int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*rate.Limiter),
		rate:     requestsPerSecond,
		burst:    requestsPerSecond * burstMultiplier,
	}

	// Clean up old visitors every 5 minutes
	go rl.cleanupVisitors()

	return rl
}

func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(rate.Limit(rl.rate), rl.burst)
		rl.visitors[ip] = limiter
	}

	return limiter
}

func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		rl.visitors = make(map[string]*rate.Limiter)
		rl.mu.Unlock()
	}
}

// RateLimit middleware limits requests per IP
func (rl *RateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		// Try to get real IP from headers
		if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			ip = realIP
		} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			// Take first IP
			ip = forwarded
		}

		limiter := rl.getVisitor(ip)
		if !limiter.Allow() {
			http.Error(w, `{"error":{"code":"rate_limit_exceeded","message":"too many requests"}}`, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
