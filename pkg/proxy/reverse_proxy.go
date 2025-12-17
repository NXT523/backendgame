package proxy

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

// Backend đại diện cho một server backend
type Backend struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
	Connections  int64 // số connections hiện tại (cho least connections)
}

// SetAlive cập nhật trạng thái sống của backend
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()
}

// IsAlive kiểm tra backend có sống không
func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	alive := b.Alive
	b.mux.RUnlock()
	return alive
}

// ServerPool quản lý pool các backend servers
type ServerPool struct {
	backends []*Backend
	current  uint64 // index hiện tại cho round-robin
	mux      sync.RWMutex
}

// AddBackend thêm backend vào pool
func (s *ServerPool) AddBackend(backend *Backend) {
	s.mux.Lock()
	s.backends = append(s.backends, backend)
	s.mux.Unlock()
}

// GetNextPeerRoundRobin lấy backend tiếp theo theo thuật toán Round Robin
func (s *ServerPool) GetNextPeerRoundRobin() *Backend {
	s.mux.RLock()
	defer s.mux.RUnlock()

	if len(s.backends) == 0 {
		return nil
	}

	// Round-robin: tăng counter và lấy backend
	next := atomic.AddUint64(&s.current, 1)

	// Thử tìm backend sống, tối đa thử len(backends) lần
	for i := 0; i < len(s.backends); i++ {
		idx := int(next+uint64(i)) % len(s.backends)
		if s.backends[idx].IsAlive() {
			return s.backends[idx]
		}
	}
	return nil
}

// GetNextPeerLeastConnections lấy backend có ít connections nhất
func (s *ServerPool) GetNextPeerLeastConnections() *Backend {
	s.mux.RLock()
	defer s.mux.RUnlock()

	if len(s.backends) == 0 {
		return nil
	}

	var selected *Backend
	minConn := int64(^uint64(0) >> 1) // Max int64

	for _, backend := range s.backends {
		if backend.IsAlive() {
			conn := atomic.LoadInt64(&backend.Connections)
			if conn < minConn {
				minConn = conn
				selected = backend
			}
		}
	}
	return selected
}

// HealthCheck kiểm tra sức khỏe các backends
func (s *ServerPool) HealthCheck() {
	for _, b := range s.backends {
		alive := isBackendAlive(b.URL)
		b.SetAlive(alive)
		status := "UP"
		if !alive {
			status = "DOWN"
		}
		log.Printf("[HealthCheck] %s - %s", b.URL, status)
	}
}

// isBackendAlive ping backend để kiểm tra
func isBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	client := &http.Client{Timeout: timeout}

	checkURL := u.String() + "/health"
	req, err := http.NewRequest("GET", checkURL, nil)
	if err != nil {
		return false
	}

	// Đánh dấu đây là health-check nội bộ
	req.Header.Set("X-Health-Check", "true")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

// LoadBalancer là reverse proxy chính với load balancing
type LoadBalancer struct {
	serverPool *ServerPool
	algorithm  string // "round-robin" hoặc "least-connections"
}

// NewBackend tạo Backend từ url string và cấu hình reverse proxy
func NewBackend(urlStr string) (*Backend, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("[ProxyError] %s: %v", u, err)
		w.WriteHeader(http.StatusBadGateway)
	}
	return &Backend{URL: u, Alive: true, ReverseProxy: proxy}, nil
}

// NewServerPool tạo pool mới
func NewServerPool() *ServerPool {
	return &ServerPool{backends: make([]*Backend, 0)}
}

// NewLoadBalancer tạo LoadBalancer
func NewLoadBalancer(pool *ServerPool, algorithm string) *LoadBalancer {
	return &LoadBalancer{serverPool: pool, algorithm: algorithm}
}

// GetNextPeer lấy backend theo algorithm đã chọn
func (lb *LoadBalancer) GetNextPeer() *Backend {
	switch lb.algorithm {
	case "least-connections":
		return lb.serverPool.GetNextPeerLeastConnections()
	default: // round-robin
		return lb.serverPool.GetNextPeerRoundRobin()
	}
}

// ServeHTTP xử lý request và forward tới backend
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	attempts := 0
	maxAttempts := 3

	for attempts < maxAttempts {
		peer := lb.GetNextPeer()
		if peer == nil {
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}

		// Tăng connection count
		atomic.AddInt64(&peer.Connections, 1)
		defer atomic.AddInt64(&peer.Connections, -1)

		// Log request
		log.Printf("[%s] %s %s -> %s", lb.algorithm, r.Method, r.URL.Path, peer.URL)

		// Custom response writer để catch error
		rw := &responseWriter{ResponseWriter: w}
		peer.ReverseProxy.ServeHTTP(rw, r)

		// Nếu thành công, return
		if rw.statusCode < 500 || rw.statusCode == 0 {
			return
		}

		// Nếu lỗi 5xx, đánh dấu backend down tạm thời và retry
		log.Printf("[Error] Backend %s returned %d, retrying...", peer.URL, rw.statusCode)
		peer.SetAlive(false)
		attempts++
	}

	http.Error(w, "All backends failed", http.StatusBadGateway)
}

// responseWriter wrapper để track status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// StartHealthCheck chạy health check định kỳ (gọi trong goroutine)
func StartHealthCheck(pool *ServerPool, interval time.Duration, stopCh <-chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pool.HealthCheck()
		case <-stopCh:
			return
		}
	}
}

// GracefulShutdown helper để đóng context nếu cần
func GracefulShutdown(ctx context.Context) error {
	// placeholder nếu cần logic shutdown đặc thù ở tương lai
	<-ctx.Done()
	return ctx.Err()
}
