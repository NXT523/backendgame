package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	GatewayLogger = NewAsyncFileLogger("internal/logs/http_gateway.log")
	GinLogger     = NewAsyncFileLogger("internal/logs/http_gin.log")
	UnaryLogger   = NewAsyncFileLogger("internal/logs/grpc_unary.log")
	StreamLogger  = NewAsyncFileLogger("internal/logs/grpc_stream.log")
)

type HTTPLogData struct {
	Method        string
	Path          string
	Query         string
	ClientIP      string
	UserAgent     string
	Referer       string
	Protocol      string
	Host          string
	BackendHost   string
	RemoteAddr    string
	RequestURI    string
	ContentLength int64
	Headers       any
	RequestBody   any
	StatusCode    int
	ResponseBody  any
	DurationMs    int64
}

type HTTPGinResponseRecorder struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *HTTPGinResponseRecorder) Write(data []byte) (n int, err error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

type HTTPGatewayResponseRecorder struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func (w *HTTPGatewayResponseRecorder) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *HTTPGatewayResponseRecorder) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func HTTPLoggerMiddlewareGin() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		var body any
		raw, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(raw))
		if strings.HasPrefix(c.GetHeader("Content-Type"), "application/json") {
			_ = json.Unmarshal(raw, &body)
		}

		writer := &HTTPGinResponseRecorder{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = writer

		c.Next()

		ghiHTTPLog(GinLogger, HTTPLogData{
			Method:        c.Request.Method,
			Path:          c.Request.URL.Path,
			Query:         c.Request.URL.RawQuery,
			ClientIP:      c.ClientIP(),
			UserAgent:     c.Request.UserAgent(),
			Referer:       c.Request.Referer(),
			Protocol:      c.Request.Proto,
			Host:          c.Request.Host,
			RemoteAddr:    c.Request.RemoteAddr,
			RequestURI:    c.Request.RequestURI,
			ContentLength: int64(len(raw)),
			StatusCode:    c.Writer.Status(),
			RequestBody:   body,
			ResponseBody:  writer.body.String(),
			DurationMs:    time.Since(start).Milliseconds(),
		})
	}
}

func HTTPLoggerMiddlewareGateway(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var requestBody any
		var raw []byte
		if r.ContentLength > 0 && r.ContentLength < 1<<20 {
			raw, _ = io.ReadAll(r.Body)
		}
		r.Body = io.NopCloser(bytes.NewBuffer(raw))

		ct := r.Header.Get("Content-Type")
		if strings.HasPrefix(ct, "application/json") {
			_ = json.Unmarshal(raw, &requestBody)
		}

		rec := &HTTPGatewayResponseRecorder{
			ResponseWriter: w,
			status:         200,
			body:           bytes.NewBuffer(nil),
		}

		next.ServeHTTP(rec, r)

		var respBody any
		respRaw := rec.body.String()
		if strings.HasPrefix(rec.Header().Get("Content-Type"), "application/json") {
			_ = json.Unmarshal([]byte(respRaw), &respBody)
		} else {
			respBody = respRaw
		}

		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip, _, _ = net.SplitHostPort(r.RemoteAddr)
		}

		ghiHTTPLog(GatewayLogger, HTTPLogData{
			Method:        r.Method,
			Path:          r.URL.Path,
			Query:         r.URL.RawQuery,
			ClientIP:      ip,
			UserAgent:     r.UserAgent(),
			Referer:       r.Referer(),
			Protocol:      r.Proto,
			Host:          r.Host,
			BackendHost:   r.Header.Get("X-Backend-Host"),
			RemoteAddr:    r.RemoteAddr,
			RequestURI:    r.RequestURI,
			ContentLength: int64(len(raw)),
			Headers:       r.Header,
			RequestBody:   requestBody,
			StatusCode:    rec.status,
			ResponseBody:  respBody,
			DurationMs:    time.Since(start).Milliseconds(),
		})
	})
}

func ghiHTTPLog(logger zerolog.Logger, data HTTPLogData) {
	event := logger.Info()
	if data.StatusCode >= 500 {
		event = logger.Error()
	} else if data.StatusCode >= 400 {
		event = logger.Warn()
	}

	event.
		Str("method", data.Method).
		Str("path", data.Path).
		Str("query", data.Query).
		Str("client_ip", data.ClientIP).
		Str("user_agent", data.UserAgent).
		Str("referer", data.Referer).
		Str("protocol", data.Protocol).
		Str("host", data.Host).
		Str("backend_host", data.BackendHost).
		Str("remote_addr", data.RemoteAddr).
		Str("request_uri", data.RequestURI).
		Int64("content_length", data.ContentLength).
		Interface("headers", data.Headers).
		Interface("request_body", data.RequestBody).
		Int("status_code", data.StatusCode).
		Interface("response_body", data.ResponseBody).
		Int64("duration_ms", data.DurationMs).
		Msg("HTTP Request Log")
}

func GRPCUnaryLoggerMiddleware() grpc.UnaryServerInterceptor {

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		st, _ := status.FromError(err)
		md, _ := metadata.FromIncomingContext(ctx)

		event := UnaryLogger.Info()
		if err != nil {
			event = UnaryLogger.Error()
		}

		event.
			Str("grpc_type", "unary").
			Str("method", info.FullMethod).
			Int("grpc_code", int(st.Code())).
			Str("grpc_message", st.Message()).
			Interface("metadata", md).
			Interface("request", req).
			Interface("response", resp).
			Int64("duration_ms", duration.Milliseconds()).
			Msg("gRPC Unary Request")

		return resp, err
	}
}

func GRPCStreamLoggerMiddleware() grpc.StreamServerInterceptor {

	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {

		start := time.Now()
		err := handler(srv, ss)
		duration := time.Since(start)

		event := StreamLogger.Info()
		if err != nil {
			event = StreamLogger.Error()
		}

		event.
			Str("grpc_type", "stream").
			Str("method", info.FullMethod).
			Bool("is_client_stream", info.IsClientStream).
			Bool("is_server_stream", info.IsServerStream).
			Int64("duration_ms", duration.Milliseconds()).
			Err(err).
			Msg("gRPC Stream Request")

		return err
	}
}

func NewAsyncFileLogger(path string) zerolog.Logger {
	lWriter := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    100,
		MaxBackups: 10,
		Compress:   false,
	}

	// Tăng buffer lên 100.000 để chịu được "cú đấm" 10k req/s
	// Giảm thời gian flush xuống 5ms để xả log nhanh hơn
	wr := diode.NewWriter(lWriter, 100000, 5*time.Millisecond, func(dropped int) {
		fmt.Printf("CRITICAL: Logger dropped %d messages at %s\n", dropped, path)
	})

	return zerolog.New(wr).With().Timestamp().Logger()
}

func newKafkaLogger() zerolog.Logger {
	return zerolog.New(&lumberjack.Logger{
		Filename:   "internal/logs/kafka.log",
		MaxSize:    50,
		MaxBackups: 5,
		MaxAge:     7,
		Compress:   true,
	}).With().Timestamp().Str("component", "kafka-consumer").Logger()
}

func formatFileSize(size int64) string {
	switch {
	case size >= 1<<20:
		return fmt.Sprintf("%.2f MB", float64(size)/(1<<20))
	case size >= 1<<10:
		return fmt.Sprintf("%.2f KB", float64(size)/(1<<10))
	default:
		return fmt.Sprintf("%d B", size)
	}
}
