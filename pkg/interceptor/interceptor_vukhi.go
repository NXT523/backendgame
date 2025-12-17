package interceptor

import (
	"context"
	v1 "game/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ----------------- Map method -> slug -----------------
var VuKhiMethodToSlug = map[string]string{
	"/v1.VuKhiService/Search":               "Search",
	"/v1.VuKhiService/GetAllVuKhi":          "GetAllVuKhi",
	"/v1.VuKhiService/GetAllTenVuKhi":       "GetAllTenVuKhi",
	"/v1.VuKhiService/GetAllSatThuongCoBan": "GetAllSatThuongCoBan",
	"/v1.VuKhiService/GetAllTocDoDanh":      "GetAllTocDoDanh",
	"/v1.VuKhiService/GetAllTamDanh":        "GetAllTamDanh",
	"/v1.VuKhiService/CreateVuKhi":          "CreateVuKhi",
	"/v1.VuKhiService/UpdateByNameVuKhi":    "UpdateByNameVuKhi",
	"/v1.VuKhiService/DeleteByNameVuKhi":    "DeleteByNameVuKhi",
}

// ----------------- Masked response -----------------
func BlockedResponseVuKhi(fullMethod string) (interface{}, bool) {
	switch fullMethod {
	case "/v1.VuKhiService/Search",
		"/v1.VuKhiService/GetAllVuKhi":
		return &v1.DanhSachVuKhi{}, true

	case "/v1.VuKhiService/GetAllTenVuKhi":
		return &v1.DanhSachStringVuKhi{}, true

	case "/v1.VuKhiService/GetAllSatThuongCoBan",
		"/v1.VuKhiService/GetAllTamDanh":
		return &v1.DanhSachIntVuKhi{}, true

	case "/v1.VuKhiService/GetAllTocDoDanh":
		return &v1.DanhSachDoubleVuKhi{}, true

	case "/v1.VuKhiService/CreateVuKhi",
		"/v1.VuKhiService/UpdateByNameVuKhi":
		return &v1.VuKhi{}, true

	case "/v1.VuKhiService/DeleteByNameVuKhi":
		return &v1.XoaTheoTenVuKhiResponse{}, true
	}
	return nil, false
}

func VuKhiAuthMaskInterceptorWithSlug(expectedSlug string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		pm := parseAndCheckMeta(md, expectedSlug)

		if !pm.BasicOK || !pm.TimeIsValid {
			if masked, ok := BlockedResponseVuKhi(info.FullMethod); ok {
				return masked, nil
			}
			return nil, status.Error(codes.PermissionDenied, "blocked")
		}

		return handler(ctx, req)
	}
}
