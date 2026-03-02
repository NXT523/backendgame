package server

import (
	"context"
	"log"
	"net/http"
)

// Chạy khởi tạo server
func RunServer(name, addr string, srv *http.Server) {
	go func() {
		log.Printf("🚀 %s chạy tại http://localhost%s", name, addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[%s] ListenAndServe lỗi: %v", name, err)
		}
	}()
}

func ShutdownAll(ctx context.Context, servers ...*http.Server) {
	for _, s := range servers {
		_ = s.Shutdown(ctx)
	}
}
