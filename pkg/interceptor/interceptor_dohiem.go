package interceptor

import (
	"context"

	v1 "game/v1/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// map method gRPC -> slug cho phép
var DoHiemMethodToSlug = map[string]string{
	"/v1.DoHiemService/CreateDoHiem":       "CreateDoHiem",
	"/v1.DoHiemService/UpdateByNameDoHiem": "UpdateByNameDoHiem",
	"/v1.DoHiemService/DeleteByNameDoHiem": "DeleteByNameDoHiem",

	"/v1.DoHiemService/GetAllDoHiem":          "GetAllDoHiem",
	"/v1.DoHiemService/GetAllTenDoHiem":       "GetAllTenDoHiem",
	"/v1.DoHiemService/GetAllSoLuongDoHiem":   "GetAllSoLuongDoHiem",
	"/v1.DoHiemService/GetAllMauSacDoHiem":    "GetAllMauSacDoHiem",
	"/v1.DoHiemService/GetAllSatThuongDoHiem": "GetAllSatThuongDoHiem",
	"/v1.DoHiemService/GetAllTocDoDoHiem":     "GetAllTocDoDoHiem",

	"/v1.DoHiemService/SearchDoHiem": "SearchDoHiem",
}

// blockedResponseDoHiem trả response rỗng đúng protobuf
func blockedResponseDoHiem(fullMethod string) (interface{}, bool) {
	switch fullMethod {
	case "/v1.DoHiemService/GetAllDoHiem", "/v1.DoHiemService/SearchDoHiem":
		return &v1.DanhSachDoHiem{}, true

	case "/v1.DoHiemService/GetAllTenDoHiem", "/v1.DoHiemService/GetAllMauSacDoHiem":
		return &v1.DanhSachStringDoHiem{}, true

	case "/v1.DoHiemService/GetAllSoLuongDoHiem":
		return &v1.DanhSachIntDoHiem{}, true

	case "/v1.DoHiemService/GetAllSatThuongDoHiem", "/v1.DoHiemService/GetAllTocDoDoHiem":
		return &v1.DanhSachDoubleDoHiem{}, true

	case "/v1.DoHiemService/CreateDoHiem", "/v1.DoHiemService/UpdateByNameDoHiem":
		return &v1.DoHiem{}, true

	case "/v1.DoHiemService/DeleteByNameDoHiem":
		return &v1.XoaTheoTenDoHiemResponse{}, true
	}
	return nil, false
}

func DoHiemAuthMaskInterceptorWithSlug(expectedSlug string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		pm := parseAndCheckMeta(md, expectedSlug)

		if !pm.BasicOK || !pm.TimeIsValid {
			if masked, ok := blockedResponseDoHiem(info.FullMethod); ok {
				return masked, nil
			}
			return nil, status.Error(codes.PermissionDenied, "blocked")
		}

		return handler(ctx, req)
	}
}
