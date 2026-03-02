package proxy

import (
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
)

type BackendNode struct {
	URL         *url.URL
	Proxy       ReverseProxy
	connections int64
	maxConn     int64
	alive       bool
	mu          sync.RWMutex
}

func NewBackendNode(u *url.URL, proxy ReverseProxy, maxConn int64) *BackendNode {
	return &BackendNode{
		URL:     u,
		Proxy:   proxy,
		alive:   true,
		maxConn: maxConn,
	}
}

func (b *BackendNode) IsAlive() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.alive
}

func (b *BackendNode) SetAlive(v bool) {
	b.mu.Lock()
	b.alive = v
	b.mu.Unlock()
}

func (b *BackendNode) Overloaded() bool {
	if b.maxConn <= 0 {
		return false
	}
	return atomic.LoadInt64(&b.connections) >= b.maxConn
}

func (b *BackendNode) Inc() { atomic.AddInt64(&b.connections, 1) }
func (b *BackendNode) Dec() { atomic.AddInt64(&b.connections, -1) }

// =======================

type BackendPool struct {
	list []*BackendNode
	mu   sync.RWMutex
}

func NewBackendPool() *BackendPool {
	return &BackendPool{}
}

func (p *BackendPool) Add(b *BackendNode) {
	p.mu.Lock()
	p.list = append(p.list, b)
	p.mu.Unlock()
}

// Logic Least Connections
func (p *BackendPool) LeastConn() *BackendNode {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var pick *BackendNode
	// Mẹo: Khởi tạo min bằng MaxInt64
	min := int64(^uint64(0) >> 1)

	for _, b := range p.list {
		if !b.IsAlive() || b.Overloaded() {
			continue
		}
		c := atomic.LoadInt64(&b.connections)
		if c < min {
			min = c
			pick = b
		}
	}
	return pick
}

// Logic Round Robin (Thêm vào nếu chưa có file round_robin.go)
type RoundRobin struct {
	pool    *BackendPool
	current uint64
}

func NewRoundRobin(pool *BackendPool) *RoundRobin {
	return &RoundRobin{pool: pool}
}

func (rr *RoundRobin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rr.pool.mu.RLock()
	defer rr.pool.mu.RUnlock()

	backends := rr.pool.list
	if len(backends) == 0 {
		http.Error(w, "No backends", http.StatusServiceUnavailable)
		return
	}

	// Chọn round robin đơn giản
	next := atomic.AddUint64(&rr.current, 1)
	idx := next % uint64(len(backends))
	node := backends[idx]

	node.Inc()
	defer node.Dec()

	// Dùng lại responseWriter có Flush
	rw := &responseWriter{ResponseWriter: w, status: 200}
	node.Proxy.ServeHTTP(rw, r)
}
