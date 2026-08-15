package middleware

import (
	"net/http"
	"sync"
	"time"
)

type visitor struct {
	count   int
	resetAt time.Time
}

// RateLimiter is an in-memory per-remote-address sliding-window limiter.
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	limit    int
	window   time.Duration
}

// NewRateLimiter builds a limiter that allows limit requests per window per
// address and starts its background cleanup loop.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		limit:    limit,
		window:   window,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(rl.window)
		rl.mu.Lock()
		now := time.Now()
		for ip, v := range rl.visitors {
			if now.After(v.resetAt) {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		rl.mu.Lock()
		v, exists := rl.visitors[ip]
		now := time.Now()

		if !exists || now.After(v.resetAt) {
			rl.visitors[ip] = &visitor{count: 1, resetAt: now.Add(rl.window)}
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		v.count++
		if v.count > rl.limit {
			rl.mu.Unlock()
			w.Header().Set("Retry-After", v.resetAt.Format(time.RFC1123))
			http.Error(w, `{"error":{"code":"rate_limited","message":"too many requests"}}`, http.StatusTooManyRequests)
			return
		}

		rl.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
