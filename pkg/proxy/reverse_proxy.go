package proxy

import (
	"log"
	"net/http"
	"net/url"

	"game/internal/utils"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// Helper để tạo Handler hỗ trợ H2C
func createH2CHandler(handler http.Handler) http.Handler {
	h2s := &http2.Server{}
	return h2c.NewHandler(handler, h2s)
}

func NewLeastConnServer(name string, backends []string, maxConn int64) *http.Server {
	pool := NewBackendPool()
	for _, raw := range backends {
		u, err := url.Parse(raw)
		if err != nil {
			log.Printf("[%s] Error url: %s", name, raw)
			continue
		}
		proxy := NewReverseProxy(u)
		pool.Add(NewBackendNode(u, proxy, maxConn))
		log.Printf("[%s] +Backend: %s", name, u)
	}

	return &http.Server{
		Addr: utils.GetenvString("PROXY_PORT_LEAST_CONNECTIONS", ":9000"),
		// Dùng h2c để nhận gRPC từ Client
		Handler: createH2CHandler(NewLeastConnLB(pool)),
	}
}

func NewRoundRobinServer(name string, backends []string, maxConn int64) *http.Server {
	pool := NewBackendPool()
	for _, raw := range backends {
		u, err := url.Parse(raw)
		if err != nil {
			log.Printf("[%s] Error url: %s", name, raw)
			continue
		}
		proxy := NewReverseProxy(u)
		pool.Add(NewBackendNode(u, proxy, maxConn))
		log.Printf("[%s] +Backend: %s", name, u)
	}

	// RoundRobin cũng phải hỗ trợ h2c
	return &http.Server{
		Addr:    utils.GetenvString("PROXY_PORT_ROUND_ROBIN", ":9010"),
		Handler: createH2CHandler(NewRoundRobin(pool)),
	}
}
