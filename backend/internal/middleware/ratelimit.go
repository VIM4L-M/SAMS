package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type tokenBucket struct {
	tokens   float64
	capacity float64
	rate     float64 // tokens per second
	lastFill time.Time
	mu       sync.Mutex
}

func (b *tokenBucket) allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.lastFill).Seconds()
	b.tokens = min(b.capacity, b.tokens+elapsed*b.rate)
	b.lastFill = now

	if b.tokens >= 1 {
		b.tokens--
		return true
	}
	return false
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

type RateLimiter struct {
	buckets  map[string]*tokenBucket
	mu       sync.Mutex
	capacity float64
	rate     float64
}

func NewRateLimiter(requestsPerSecond float64, burst float64) *RateLimiter {
	rl := &RateLimiter{
		buckets:  make(map[string]*tokenBucket),
		capacity: burst,
		rate:     requestsPerSecond,
	}
	// Periodic cleanup of stale buckets
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		for ip, b := range rl.buckets {
			b.mu.Lock()
			if time.Since(b.lastFill) > 10*time.Minute {
				delete(rl.buckets, ip)
			}
			b.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) getBucket(ip string) *tokenBucket {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	b, ok := rl.buckets[ip]
	if !ok {
		b = &tokenBucket{
			tokens:   rl.capacity,
			capacity: rl.capacity,
			rate:     rl.rate,
			lastFill: time.Now(),
		}
		rl.buckets[ip] = b
	}
	return b
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		// Respect X-Forwarded-For only in trusted deployments; omit for now
		if !rl.getBucket(ip).allow() {
			http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
