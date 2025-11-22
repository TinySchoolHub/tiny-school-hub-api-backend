package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	rps := 10
	rl := NewRateLimiter(rps)

	if rl == nil {
		t.Fatal("NewRateLimiter() returned nil")
	}
	if rl.rate != rps {
		t.Errorf("Rate = %v, want %v", rl.rate, rps)
	}
	if rl.burst != rps*2 {
		t.Errorf("Burst = %v, want %v", rl.burst, rps*2)
	}
	if rl.visitors == nil {
		t.Error("Visitors map is nil")
	}
}

func TestRateLimit_AllowRequests(t *testing.T) {
	rl := NewRateLimiter(10) // 10 requests per second

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := rl.RateLimit(handler)

	// Should allow first request
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Status code = %v, want %v", rr.Code, http.StatusOK)
	}
}

func TestRateLimit_ExceedLimit(t *testing.T) {
	rl := NewRateLimiter(1) // 1 request per second, burst of 2

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := rl.RateLimit(handler)

	ip := "192.168.1.1"

	// First two requests should succeed (burst)
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
		req.RemoteAddr = ip
		rr := httptest.NewRecorder()

		middleware.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Request %d: Status code = %v, want %v", i+1, rr.Code, http.StatusOK)
		}
	}

	// Third request should be rate limited
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.RemoteAddr = ip
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Status code = %v, want %v", rr.Code, http.StatusTooManyRequests)
	}
}

func TestRateLimit_DifferentIPs(t *testing.T) {
	rl := NewRateLimiter(1) // 1 request per second

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := rl.RateLimit(handler)

	// Different IPs should have separate rate limits
	ips := []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"}

	for _, ip := range ips {
		// Each IP should be allowed their burst
		for i := 0; i < 2; i++ {
			req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
			req.RemoteAddr = ip
			rr := httptest.NewRecorder()

			middleware.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("IP %s, Request %d: Status code = %v, want %v", ip, i+1, rr.Code, http.StatusOK)
			}
		}
	}
}

func TestRateLimit_XRealIP(t *testing.T) {
	rl := NewRateLimiter(1) // 1 request per second, burst of 2

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := rl.RateLimit(handler)

	realIP := "10.0.0.1"

	// Use X-Real-IP header
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
		req.RemoteAddr = "192.168.1.1" // Different from X-Real-IP
		req.Header.Set("X-Real-IP", realIP)
		rr := httptest.NewRecorder()

		middleware.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Request %d: Status code = %v, want %v", i+1, rr.Code, http.StatusOK)
		}
	}

	// Third request should be rate limited
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.RemoteAddr = "192.168.1.1"
	req.Header.Set("X-Real-IP", realIP)
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Status code = %v, want %v", rr.Code, http.StatusTooManyRequests)
	}
}

func TestRateLimit_XForwardedFor(t *testing.T) {
	rl := NewRateLimiter(1) // 1 request per second, burst of 2

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := rl.RateLimit(handler)

	forwardedIP := "203.0.113.1"

	// Use X-Forwarded-For header
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
		req.RemoteAddr = "192.168.1.1"
		req.Header.Set("X-Forwarded-For", forwardedIP+", 192.168.1.1")
		rr := httptest.NewRecorder()

		middleware.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Request %d: Status code = %v, want %v", i+1, rr.Code, http.StatusOK)
		}
	}

	// Third request should be rate limited
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.RemoteAddr = "192.168.1.1"
	req.Header.Set("X-Forwarded-For", forwardedIP+", 192.168.1.1")
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Status code = %v, want %v", rr.Code, http.StatusTooManyRequests)
	}
}

func TestRateLimit_Recovery(t *testing.T) {
	rl := NewRateLimiter(5) // 5 requests per second

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := rl.RateLimit(handler)

	ip := "192.168.1.1"

	// Exhaust burst (10 requests)
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
		req.RemoteAddr = ip
		rr := httptest.NewRecorder()
		middleware.ServeHTTP(rr, req)
	}

	// Should be rate limited
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.RemoteAddr = ip
	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Status code = %v, want %v", rr.Code, http.StatusTooManyRequests)
	}

	// Wait for rate limiter to recover (1 second / 5 rps = 200ms per request)
	time.Sleep(250 * time.Millisecond)

	// Should allow request after recovery
	req = httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.RemoteAddr = ip
	rr = httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("After recovery: Status code = %v, want %v", rr.Code, http.StatusOK)
	}
}

func TestGetVisitor(t *testing.T) {
	rl := NewRateLimiter(10)

	ip1 := "192.168.1.1"
	ip2 := "192.168.1.2"

	// Get limiter for IP1
	limiter1 := rl.getVisitor(ip1)
	if limiter1 == nil {
		t.Error("getVisitor() returned nil for ip1")
	}

	// Get limiter for IP1 again - should be same instance
	limiter1Again := rl.getVisitor(ip1)
	if limiter1 != limiter1Again {
		t.Error("getVisitor() should return same limiter instance for same IP")
	}

	// Get limiter for IP2 - should be different
	limiter2 := rl.getVisitor(ip2)
	if limiter2 == nil {
		t.Error("getVisitor() returned nil for ip2")
	}
	if limiter1 == limiter2 {
		t.Error("getVisitor() should return different limiters for different IPs")
	}
}

func BenchmarkRateLimit(b *testing.B) {
	rl := NewRateLimiter(1000)

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {})

	middleware := rl.RateLimit(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.RemoteAddr = "192.168.1.1"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		middleware.ServeHTTP(rr, req)
	}
}

func BenchmarkRateLimit_MultipleIPs(b *testing.B) {
	rl := NewRateLimiter(1000)

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {})

	middleware := rl.RateLimit(handler)

	ips := []string{
		"192.168.1.1",
		"192.168.1.2",
		"192.168.1.3",
		"192.168.1.4",
		"192.168.1.5",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
		req.RemoteAddr = ips[i%len(ips)]
		rr := httptest.NewRecorder()
		middleware.ServeHTTP(rr, req)
	}
}
