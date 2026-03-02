package interceptor

import (
	"context"
	v1 "game/v1/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// map method gRPC -> slug cho LoaiVuKhiService
var LoaiVuKhiMethodToSlug = map[string]string{
	"/v1.LoaiVuKhiService/CreateLoaiVuKhi":       "CreateLoaiVuKhi",
	"/v1.LoaiVuKhiService/UpdateByNameLoaiVuKhi": "UpdateByNameLoaiVuKhi",
	"/v1.LoaiVuKhiService/DeleteByNameLoaiVuKhi": "DeleteByNameLoaiVuKhi",

	"/v1.LoaiVuKhiService/GetAllLoaiVuKhi":     "GetAllLoaiVuKhi",
	"/v1.LoaiVuKhiService/GetAllTenLoaiVuKhi":  "GetAllTenLoaiVuKhi",
	"/v1.LoaiVuKhiService/GetAllMoTaLoaiVuKhi": "GetAllMoTaLoaiVuKhi",

	"/v1.LoaiVuKhiService/SearchLoaiVuKhi": "SearchLoaiVuKhi",
}

// maskedResponse trả response rỗng đúng proto, theo từng method
func blockedResponseLoaiVuKhi(fullMethod string) (interface{}, bool) {
	switch fullMethod {
	// danh sách LoaiVuKhi
	case "/v1.LoaiVuKhiService/GetAllLoaiVuKhi",
		"/v1.LoaiVuKhiService/SearchLoaiVuKhi":
		return &v1.DanhSachLoaiVuKhi{}, true

	// danh sách string
	case "/v1.LoaiVuKhiService/GetAllTenLoaiVuKhi",
		"/v1.LoaiVuKhiService/GetAllMoTaLoaiVuKhi":
		return &v1.DanhSachStringLoaiVuKhi{}, true

	// LoaiVuKhi đơn
	case "/v1.LoaiVuKhiService/CreateLoaiVuKhi",
		"/v1.LoaiVuKhiService/UpdateByNameLoaiVuKhi":
		return &v1.LoaiVuKhi{}, true

	// response xóa
	case "/v1.LoaiVuKhiService/DeleteByNameLoaiVuKhi":
		return &v1.XoaTheoTenLoaiVuKhiResponse{}, true
	}
	return nil, false
}

func LoaiVuKhiAuthMaskInterceptorWithSlug(expectedSlug string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		pm := parseAndCheckMeta(md, expectedSlug)

		if !pm.BasicOK || !pm.TimeIsValid {
			if masked, ok := blockedResponseLoaiVuKhi(info.FullMethod); ok {
				return masked, nil
			}
			return nil, status.Error(codes.PermissionDenied, "blocked")
		}

		return handler(ctx, req)
	}
}
