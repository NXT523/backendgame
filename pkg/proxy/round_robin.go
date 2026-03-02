package proxy

// import (
// 	"log"
// 	"net/http"
// 	"sync/atomic"
// )

// type RoundRobin struct {
// 	backends []*BackendNode
// 	idx      uint64
// }

// func NewRoundRobin(pool *BackendPool) *RoundRobin {
// 	return &RoundRobin{backends: pool.All()}
// }

// func (rr *RoundRobin) next() *BackendNode {
// 	n := len(rr.backends)
// 	for i := 0; i < n; i++ {
// 		idx := int(atomic.AddUint64(&rr.idx, 1) % uint64(n))
// 		b := rr.backends[idx]
// 		if b.IsAlive() && !b.Overloaded() {
// 			return b
// 		}
// 	}
// 	return nil
// }

// func (rr *RoundRobin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	b := rr.next()
// 	if b == nil {
// 		http.Error(w, "no upstream", http.StatusServiceUnavailable)
// 		return
// 	}

// 	b.Inc()
// 	defer b.Dec()

// 	log.Printf(
// 		"[RR-PROXY] %s %s | Client: %s -> Backend: %s",
// 		r.Method,
// 		r.URL.Path,
// 		r.RemoteAddr,
// 		b.URL.String(),
// 	)
// 	b.Proxy.ServeHTTP(w, r)
// }
