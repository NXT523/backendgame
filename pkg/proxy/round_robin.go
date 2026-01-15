package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

// UPSTREAM

type Upstream struct {
	URL        *url.URL
	Proxy      *httputil.ReverseProxy
	alive      int32
	activeConn int64
	maxConn    int64
}

func NewUpstream(rawURL string, maxConn int64) (*Upstream, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	up := &Upstream{
		URL:     u,
		alive:   1,
		maxConn: maxConn,
	}

	proxy := httputil.NewSingleHostReverseProxy(u)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = u.Host
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		atomic.StoreInt32(&up.alive, 0)
		log.Printf("[RR] %s error: %v", u, err)
		w.WriteHeader(http.StatusBadGateway)
	}

	up.Proxy = proxy
	return up, nil
}

func (u *Upstream) IsAlive() bool {
	return atomic.LoadInt32(&u.alive) == 1
}

func (u *Upstream) CanAccept() bool {
	if !u.IsAlive() {
		return false
	}
	if u.maxConn <= 0 {
		return true
	}
	return atomic.LoadInt64(&u.activeConn) < u.maxConn
}

// REGISTRY

type UpstreamRegistry struct {
	mu        sync.RWMutex
	upstreams []*Upstream
}

func NewUpstreamRegistry() *UpstreamRegistry {
	return &UpstreamRegistry{}
}

func (r *UpstreamRegistry) Register(up *Upstream) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.upstreams = append(r.upstreams, up)
}

func (r *UpstreamRegistry) List() []*Upstream {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]*Upstream(nil), r.upstreams...)
}

// ROUND ROBIN

type RoundRobin struct {
	registry *UpstreamRegistry
	idx      uint64
}

func NewRoundRobin(reg *UpstreamRegistry) *RoundRobin {
	return &RoundRobin{registry: reg}
}

func (rr *RoundRobin) next() *Upstream {
	ups := rr.registry.List()
	n := len(ups)
	if n == 0 {
		return nil
	}

	for i := 0; i < n; i++ {
		idx := int(atomic.AddUint64(&rr.idx, 1) % uint64(n))
		up := ups[idx]
		if up.CanAccept() {
			return up
		}
	}
	return nil
}

// HTTP HANDLER

func (rr *RoundRobin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	up := rr.next()
	if up == nil {
		http.Error(w, "no healthy upstream", http.StatusServiceUnavailable)
		return
	}

	atomic.AddInt64(&up.activeConn, 1)
	defer atomic.AddInt64(&up.activeConn, -1)

	start := time.Now()
	log.Printf("[RR] %s -> %s", r.URL.Path, up.URL)

	up.Proxy.ServeHTTP(w, r)

	log.Printf("[RR] done %s (%v)", up.URL, time.Since(start))
}
