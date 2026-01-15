package proxy

import (
	"bytes"
	grpcsrv "game/internal/grpc"
	"io"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

// LEAST CONNECTIONS LOAD BALANCER

type LeastConnLoadBalancer struct {
	pool *BackendPool
}

func NewLeastConnLoadBalancer(pool *BackendPool) *LeastConnLoadBalancer {
	return &LeastConnLoadBalancer{pool: pool}
}

func (lb *LeastConnLoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const maxRetry = 2

	// 🔒 cache body để retry an toàn
	var body []byte
	if r.Body != nil {
		body, _ = io.ReadAll(r.Body)
		_ = r.Body.Close()
	}

	for attempt := 0; attempt <= maxRetry; attempt++ {
		backend := lb.pool.GetLeastLoadedBackend()
		if backend == nil {
			http.Error(w, "No backend available", http.StatusTooManyRequests)
			return
		}

		start := time.Now()
		atomic.AddInt64(&backend.Connections, 1)

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		backend.ReverseProxy.ServeHTTP(
			rw,
			cloneRequestWithBody(r, body),
		)

		atomic.AddInt64(&backend.Connections, -1)

		status := rw.statusCode
		elapsed := time.Since(start).Milliseconds()

		grpcsrv.AccessLogger.Printf(
			"[LC] %s %s -> %s | status=%d | %dms | conn=%d/%d",
			r.Method,
			r.URL.Path,
			backend.URL,
			status,
			elapsed,
			atomic.LoadInt64(&backend.Connections),
			backend.MaxConn,
		)

		// ✅ OK / client error → trả về
		if status < 500 {
			return
		}

		// ❌ backend lỗi → loại tạm
		log.Printf("[Failover][LC] %s returned %d", backend.URL, status)
		backend.SetAlive(false)
		go recoverBackend(backend)
	}

	http.Error(w, "All backends failed", http.StatusBadGateway)
}

// REQUEST CLONE

func cloneRequestWithBody(r *http.Request, body []byte) *http.Request {
	clone := r.Clone(r.Context())
	if body != nil {
		clone.Body = io.NopCloser(bytes.NewBuffer(body))
	}
	return clone
}

