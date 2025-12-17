package logx

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
)

type ctxKey string

const CtxKeyHTTPStart ctxKey = "httpStart"


func GinGhiLogTruyCap(l *LoggerElastic) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		level := "info"
		if status >= 500 {
			level = "error"
		} else if status >= 400 {
			level = "warn"
		}

		msg := "HTTP request handled"
		if status >= 400 {
			msg = "HTTP request handled with client error"
			if status >= 500 {
				msg = "HTTP request handled with server error"
			}
		}

		// ---- Lấy IP trần (không kèm :port). Gin đã có ClientIP() nên giữ nguyên. ----
		kv := map[string]any{
			"service.component":         "gateway",
			"event.dataset":             "http",
			"http.request.method":       c.Request.Method,
			"http.request.url":          c.Request.URL.Path,
			"http.response.status_code": status,
			"client.ip":                 c.ClientIP(),
			"duration_ms":               time.Since(start).Milliseconds(),
		}
		if tid := c.Request.Header.Get("X-Trace-Id"); tid != "" {
			kv["trace_id"] = tid
		}
		if rid := c.Request.Header.Get("X-Request-Id"); rid != "" {
			kv["request_id"] = rid
		}

		l.Ghi(c.Request.Context(), level, msg, kv)
	}
}

func HttpErrorHandler(elog *LoggerElastic) runtime.ErrorHandlerFunc {
	return func(
		ctx context.Context,
		mux *runtime.ServeMux,
		m runtime.Marshaler,
		w http.ResponseWriter,
		r *http.Request,
		err error,
	) {
		st := status.Convert(err)
		httpCode := runtime.HTTPStatusFromCode(st.Code())
		httpText := http.StatusText(httpCode) // ví dụ "Internal Server Error"

		// --- sanitize message: lấy phần trước dấu ":" nếu có ---
		rawMsg := st.Message()
		sanitized := rawMsg
		if i := strings.IndexByte(rawMsg, ':'); i > 0 {
			sanitized = strings.TrimSpace(rawMsg[:i])
		}
		if httpCode >= 500 && strings.TrimSpace(sanitized) == "" {
			sanitized = httpText
		}

		// 1) Trả JSON cho client HTTP
		w.Header().Set("Content-Type", m.ContentType(nil))
		w.WriteHeader(httpCode)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code":        int32(st.Code()),
			"error":       httpText,
			"http_status": httpCode,
			"message":     sanitized,
			"path":        r.URL.Path,
			"details":     []any{},
		})

		// 2) Log sang ES (kèm client.ip & duration_ms)
		if elog != nil {
			// --- lấy start-time mà wrapper đã gắn vào context ---
			var durMs int64
			if v := ctx.Value(CtxKeyHTTPStart); v != nil {
				if t, ok := v.(time.Time); ok {
					durMs = time.Since(t).Milliseconds()
				}
			}

			// --- bóc IP trần ---
			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				host, _, err := net.SplitHostPort(r.RemoteAddr)
				if err == nil && host != "" {
					ip = host
				} else {
					ip = r.RemoteAddr
				}
			}

			level := "warn"
			if httpCode >= 500 {
				level = "error"
			}
			kv := map[string]any{
				"service.component":         "gateway",
				"event.dataset":             "http",
				"http.request.method":       r.Method,
				"http.request.url":          r.URL.Path,
				"http.response.status_code": httpCode,
				"client.ip":                 ip,
				"duration_ms":               durMs,
				"http.request.host":         r.Host,

				// đúng yêu cầu:
				"error.message": sanitized, // "Tạo vũ khí lỗi"
				"error.type":    httpText,  // "Internal Server Error" | "Not Found" ...
			}
			if tid := r.Header.Get("X-Trace-Id"); tid != "" {
				kv["trace_id"] = tid
			}
			if rid := r.Header.Get("X-Request-Id"); rid != "" {
				kv["request_id"] = rid
			}

			elog.Ghi(ctx, level, "HTTP request failed", kv)
		}
	}
}
