package interceptor

import (
	"context"
	v1 "game/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// map method gRPC -> slug cho phép
var HeMethodToSlug = map[string]string{
	"/v1.HeService/CreateHe":       "CreateHe",
	"/v1.HeService/UpdateByNameHe": "UpdateByNameHe",
	"/v1.HeService/DeleteByNameHe": "DeleteByNameHe",
	"/v1.HeService/GetAllHe":       "GetAllHe",
	"/v1.HeService/GetAllTenHe":    "GetAllTenHe",
	"/v1.HeService/GetAllMoTaHe":   "GetAllMoTaHe",
	"/v1.HeService/SearchHe":       "SearchHe",
}

// blockedResponseHe tạo response "mask" rỗng đúng kiểu protobuf
func blockedResponseHe(fullMethod string) (interface{}, bool) {
	switch fullMethod {
	case "/v1.HeService/GetAllHe", "/v1.HeService/SearchHe":
		return &v1.DanhSachHe{}, true

	case "/v1.HeService/GetAllTenHe", "/v1.HeService/GetAllMoTaHe":
		return &v1.DanhSachStringHe{}, true

	case "/v1.HeService/CreateHe", "/v1.HeService/UpdateByNameHe":
		return &v1.He{}, true

	case "/v1.HeService/DeleteByNameHe":
		return &v1.XoaTheoTenHeResponse{}, true
	}
	return nil, false
}

func HeAuthMaskInterceptorWithSlug(expectedSlug string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		pm := parseAndCheckMeta(md, expectedSlug)

		if !pm.BasicOK || !pm.TimeIsValid {
			if masked, ok := blockedResponseHe(info.FullMethod); ok {
				return masked, nil
			}
			return nil, status.Error(codes.PermissionDenied, "blocked")
		}

		return handler(ctx, req)
	}
}
