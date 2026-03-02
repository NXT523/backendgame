package proxy

import (
	"log"
	"net/http"
	"sync/atomic"
)

// LeastConnLB: Load Balancer chọn node ít kết nối nhất
type LeastConnLB struct {
	pool        *BackendPool
	totalReqOps uint64 
}

func NewLeastConnLB(pool *BackendPool) *LeastConnLB {
	return &LeastConnLB{
		pool: pool,
	}
}

func (lb *LeastConnLB) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. Chọn Backend tốt nhất hiện tại
	backend := lb.pool.LeastConn()
	if backend == nil {
		// Nếu không có backend nào sống, trả về 503
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	// 2. Tăng Counter (Concurrency tracking)
	// Việc này giúp thuật toán LeastConn biết node nào đang bận
	backend.Inc()
	defer backend.Dec()

	// 3. Log Sampling (Chỉ in 1% số lượng log để tối ưu hiệu năng)
	// Atomic Add trả về giá trị mới, ta dùng nó để chia dư
	currentOps := atomic.AddUint64(&lb.totalReqOps, 1)
	if currentOps%100 == 0 {
		log.Printf("[PROXY] %s -> %s (Active Conns: %d)",
			r.URL.Path,
			backend.URL.Host,
			atomic.LoadInt64(&backend.connections),
		)
	}

	// 4. Wrapper ResponseWriter
	// Mục đích: Để bắt được Status Code mà backend trả về (để log lỗi nếu có)
	rw := &responseWriter{
		ResponseWriter: w,
		status:         http.StatusOK, // Mặc định là 200 OK
	}

	// 5. Forward Request
	// Lưu ý: Backend Proxy đã được cấu hình FlushInterval = -1 và h2c
	// nên nó sẽ tự động stream data mà không cần buffer.
	backend.Proxy.ServeHTTP(rw, r)

	// 6. Log lỗi nếu Backend trả về 5xx (Lỗi Server)
	// Giúp debug xem server nào đang bị lỗi
	if rw.status >= 500 {
		log.Printf("[PROXY-ERR] Backend %s responded %d for %s",
			backend.URL.Host,
			rw.status,
			r.URL.Path,
		)
	}
}

// =================================================================
// Custom ResponseWriter
// BẮT BUỘC phải implement http.Flusher thì gRPC Streaming mới chạy mượt
// =================================================================
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

// Flush đẩy dữ liệu từ buffer xuống client ngay lập tức.
// Nếu không có hàm này, gRPC sẽ bị lag (do data bị om lại chờ đầy buffer).
func (w *responseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}
