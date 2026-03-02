package server

import (
	"net/http"
	"strings"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

// Hàm này có chức năng nếu các đi qua các path trong hàm thì đi qua gateway còn nếu không cùng path thì đi qua gin
func GatewayBuildH2CServer(grpcSrv *grpc.Server, gw http.Handler, ginHandler http.Handler) *http.Server {
	coreLogic := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Nếu là gRPC (HTTP/2 + grpc Content-Type)
		if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
			grpcSrv.ServeHTTP(w, r)
			return
		}

		// 2. Nếu là REST gRPC-Gateway
		path := r.URL.Path
		if strings.HasPrefix(path, "/vukhi/ent/") ||
			strings.HasPrefix(path, "/dohiem/ent/") ||
			strings.HasPrefix(path, "/he/ent/") ||
			strings.HasPrefix(path, "/loaivukhi/ent/") {
			gw.ServeHTTP(w, r)
			return
		}
		// 3. Mọi request khác đẩy về Gin
		ginHandler.ServeHTTP(w, r)
	})

	return &http.Server{
		Handler: h2c.NewHandler(coreLogic, &http2.Server{}),
	}
}
