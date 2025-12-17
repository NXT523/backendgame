package logx

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func tachRPC(fullMethod string) (svc, m string) {
	fullMethod, _ = strings.CutPrefix(fullMethod, "/")
	svc, m, _ = strings.Cut(fullMethod, "/")
	return
}

func GRPCGhiLog(l *LoggerElastic) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)

		code := status.Code(err)
		level := "info"
		switch code {
		case codes.OK:
			level = "info"
		case codes.InvalidArgument, codes.NotFound, codes.AlreadyExists, codes.PermissionDenied, codes.Unauthenticated, codes.FailedPrecondition:
			level = "warn"
		default:
			if code != codes.OK {
				level = "error"
			}
		}

		svc, met := tachRPC(info.FullMethod)
		kv := map[string]any{
			"service.component": "grpc",
			"event.dataset":     "grpc",
			"grpc.service":      svc,
			"grpc.method":       met,
			"grpc.status_code":  code.String(),
			"duration_ms":       time.Since(start).Milliseconds(),
		}
		if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
			kv["grpc.peer.address"] = p.Addr.String()
		}
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if v := md.Get("x-trace-id"); len(v) > 0 {
				kv["trace_id"] = v[0]
			}
		}

		msg := "gRPC call handled"
		if err != nil {
			st := status.Convert(err)
			kv["error.message"] = st.Message()
			kv["error.type"] = "gRPC." + code.String()
			msg = "gRPC call finished with error"
		}

		l.Ghi(ctx, level, msg, kv)
		return resp, err
	}
}
