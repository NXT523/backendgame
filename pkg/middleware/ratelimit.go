package middleware

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

/* ================= CORE ================= */

type client struct {
	tokens      int
	lastRefill  time.Time
	bannedUntil time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	clients  map[string]*client
	rate     int
	interval time.Duration
}

func TaoRateLimiter(rate int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients:  make(map[string]*client),
		rate:     rate,
		interval: interval,
	}

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			rl.cleanup()
		}
	}()

	return rl
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for k, c := range rl.clients {
		if now.Sub(c.lastRefill) > 10*rl.interval {
			delete(rl.clients, k)
		}
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	now := time.Now()

	rl.mu.Lock()
	defer rl.mu.Unlock()

	c, ok := rl.clients[key]
	if !ok {
		rl.clients[key] = &client{
			tokens:     rl.rate - 1,
			lastRefill: now,
		}
		return true
	}

	/* 🚫 ĐANG BỊ BAN */
	if now.Before(c.bannedUntil) {
		return false
	}

	/* 🔓 HẾT BAN → RESET */
	if !c.bannedUntil.IsZero() && now.After(c.bannedUntil) {
		c.bannedUntil = time.Time{}
		c.tokens = rl.rate
		c.lastRefill = now
	}

	/* ♻️ REFILL TOKEN */
	if now.Sub(c.lastRefill) >= rl.interval {
		c.tokens = rl.rate
		c.lastRefill = now
	}

	/* 🚨 HẾT TOKEN → BAN 5 PHÚT */
	if c.tokens <= 0 {
		c.bannedUntil = now.Add(5 * time.Minute)
		log.Printf("[RATE LIMIT] ban %s until %s\n",
			key,
			c.bannedUntil.Format(time.RFC3339),
		)
		return false
	}

	c.tokens--
	return true
}

/* ================= HTTP ================= */

func layIPHTTP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func (rl *RateLimiter) RateLimitHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := layIPHTTP(r)

		if !rl.Allow(ip) {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

/* ================= gRPC ================= */

func layIPGrpc(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if xff := md.Get("x-forwarded-for"); len(xff) > 0 {
			parts := strings.Split(xff[0], ",")
			return strings.TrimSpace(parts[0])
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		if addr, ok := p.Addr.(*net.TCPAddr); ok {
			return addr.IP.String()
		}
	}

	return "unknown"
}

func RateLimitGrpc(rl *RateLimiter) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		ip := layIPGrpc(ctx)

		if !rl.Allow(ip) {
			return nil, status.Error(
				codes.ResourceExhausted,
				"Too Many Requests",
			)
		}

		return handler(ctx, req)
	}
}
