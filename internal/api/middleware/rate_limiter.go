package middleware

import (
	"net/http"
	"sync"
	"time"
)

type client struct {
	tokens   int
	lastSeen time.Time
}

type RateLimiter struct {
	mu        sync.Mutex
	clients   map[string]*client
	rate      int           // tokens added per interval
	interval  time.Duration // refill interval
	maxTokens int           // bucket capacity
}

func NewRateLimiter(rate int, interval time.Duration, maxTokens int) *RateLimiter {
	rl := &RateLimiter{
		clients:   make(map[string]*client),
		rate:      rate,
		interval:  interval,
		maxTokens: maxTokens,
	}

	// Cleanup stale entries every 5 minutes
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			rl.mu.Lock()
			for ip, c := range rl.clients {
				if time.Since(c.lastSeen) > 10*time.Minute {
					delete(rl.clients, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()

	return rl
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			ip = xff
		}

		rl.mu.Lock()
		c, exists := rl.clients[ip]
		now := time.Now()

		if !exists {
			rl.clients[ip] = &client{
				tokens:   rl.maxTokens - 1,
				lastSeen: now,
			}
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		elapsed := now.Sub(c.lastSeen)
		c.lastSeen = now

		refill := int(elapsed/rl.interval) * rl.rate
		if refill > 0 {
			c.tokens = c.tokens + refill
			if c.tokens > rl.maxTokens {
				c.tokens = rl.maxTokens
			}
		}

		if c.tokens <= 0 {
			rl.mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"status":"error","message":"Too many requests. Please try again later."}`))
			return
		}

		c.tokens--
		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
