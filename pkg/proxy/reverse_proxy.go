package proxy

import (
	"context"
	grpcsrv "game/internal/grpc"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)


// BACKEND NODE


type BackendNode struct {
	URL          *url.URL
	ReverseProxy *httputil.ReverseProxy

	alive bool
	mux   sync.RWMutex

	Connections int64
	MaxConn     int64
}

func (b *BackendNode) SetAlive(v bool) {
	b.mux.Lock()
	b.alive = v
	b.mux.Unlock()
}

func (b *BackendNode) IsAlive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.alive
}

func (b *BackendNode) IsOverloaded() bool {
	return atomic.LoadInt64(&b.Connections) >= b.MaxConn
}


// BACKEND POOL


type BackendPool struct {
	backends []*BackendNode
	mux      sync.RWMutex
}

func NewBackendPool() *BackendPool {
	return &BackendPool{}
}

func (p *BackendPool) AddBackend(b *BackendNode) {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.backends = append(p.backends, b)
}

func (p *BackendPool) GetLeastLoadedBackend() *BackendNode {
	p.mux.RLock()
	defer p.mux.RUnlock()

	var selected *BackendNode
	minConn := int64(^uint64(0) >> 1)

	for _, b := range p.backends {
		if !b.IsAlive() || b.IsOverloaded() {
			continue
		}
		conn := atomic.LoadInt64(&b.Connections)
		if conn < minConn {
			minConn = conn
			selected = b
		}
	}
	return selected
}


// SMART LOAD BALANCER (LC + Retry)


type SmartLoadBalancer struct {
	pool *BackendPool
}

func NewSmartLoadBalancer(pool *BackendPool) *SmartLoadBalancer {
	return &SmartLoadBalancer{pool: pool}
}

func (lb *SmartLoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const maxRetry = 2

	var body []byte
	if r.Body != nil {
		body, _ = io.ReadAll(r.Body)
		_ = r.Body.Close()
	}

	for attempt := 0; attempt <= maxRetry; attempt++ {
		backend := lb.pool.GetLeastLoadedBackend()
		if backend == nil {
			http.Error(w, "All backends overloaded", http.StatusTooManyRequests)
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
			"[LC-SMART] %s %s -> %s | status=%d | %dms | conn=%d/%d",
			r.Method,
			r.URL.Path,
			backend.URL,
			status,
			elapsed,
			atomic.LoadInt64(&backend.Connections),
			backend.MaxConn,
		)

		if status < 500 {
			return
		}

		log.Printf("[Failover] %s returned %d", backend.URL, status)
		backend.SetAlive(false)
		go recoverBackend(backend)
	}

	http.Error(w, "All backends failed", http.StatusBadGateway)
}


// RESPONSE WRAPPER


type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}


// BACKEND FACTORY


func NewBackendNode(rawURL string, maxConn int64) (*BackendNode, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("[ProxyError] %s: %v", u, err)
		w.WriteHeader(http.StatusBadGateway)
	}

	return &BackendNode{
		URL:          u,
		ReverseProxy: proxy,
		alive:        true,
		MaxConn:      maxConn,
	}, nil
}


// HEALTH CHECK + RECOVERY


func recoverBackend(b *BackendNode) {
	time.Sleep(5 * time.Second)

	if isBackendAlive(b.URL) {
		b.SetAlive(true)
		log.Printf("[Recovery] %s back online", b.URL)
		return
	}

	log.Printf("[Recovery] %s still DOWN", b.URL)
}

func isBackendAlive(u *url.URL) bool {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(u.String() + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}


// GRACEFUL SHUTDOWN


func GracefulShutdown(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}
