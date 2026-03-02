package proxy

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"golang.org/x/net/http2"
)

type ReverseProxy interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func NewReverseProxy(u *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(u)

	// ===== HTTP/2 over h2c =====
	h2 := &http2.Transport{
		AllowHTTP: true,
		DialTLSContext: func(ctx context.Context, netw, addr string, _ *tls.Config) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, netw, addr)
		},
	}

	t := &http.Transport{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 500,
		IdleConnTimeout:     90 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	t.RegisterProtocol("http", h2)
	proxy.Transport = t
	proxy.FlushInterval = -1

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		origHost := req.Host
		req.Header.Set("X-Backend-Host", u.Host)
		req.Header.Set("X-Forwarded-Host", origHost)
		req.Header.Set("X-Forwarded-Proto", "http")
		req.Header.Set("X-Forwarded-For", req.RemoteAddr)
	}

	return proxy
}
