package interceptor

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ----------------- Cấu hình chung -----------------
const secret = "alo"

const (
	mdKeyTime   = "x-time"
	mdKeyString = "x-string"
)

// ----------------- ParsedMeta -----------------
type ParsedMeta struct {
	ServerTime  time.Time
	ClientTime  time.Time
	Slug        string
	TimeIsValid bool
	BasicOK     bool
}

// ----------------- Helpers -----------------
func firstMD(vals []string) string {
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}

func parseAndCheckMeta(md metadata.MD, expectedSlug string) ParsedMeta {
	var zero ParsedMeta

	tsStr := firstMD(md[mdKeyTime])
	xString := firstMD(md[mdKeyString])

	if tsStr == "" || xString == "" {
		// println("[DEBUG] Missing x-time or x-string")
		return zero
	}

	// parse timestamp
	unixSec, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		// println("[DEBUG] Invalid x-time:", tsStr)
		return zero
	}
	clientTS := time.Unix(unixSec, 0)

	// check ±5 phút
	if time.Since(clientTS) > 5*time.Minute || time.Until(clientTS) > 5*time.Minute {
		// println("[DEBUG] x-time out of range:", tsStr)
		return zero
	}

	// HMAC check: chỉ dùng x-time + expectedSlug
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(tsStr + expectedSlug))
	expectedXString := mac.Sum(nil)

	receivedMAC, err := base64.StdEncoding.DecodeString(xString)
	if err != nil {
		// println("[DEBUG] x-string HMAC mismatch, got:", xString, "expected:", expectedXString)
		return zero
	}
	if !hmac.Equal(receivedMAC, expectedXString) {
		return zero
	}

	// println("[DEBUG] Metadata OK: x-time =", tsStr, "x-string =", xString)

	return ParsedMeta{
		ServerTime:  time.Now(),
		ClientTime:  clientTS,
		Slug:        expectedSlug,
		TimeIsValid: true,
		BasicOK:     true,
	}
}

// CommonAuthInterceptor: cho tất cả service khác
func CommonAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		pm := parseAndCheckMeta(md, "")

		if !pm.BasicOK || !pm.TimeIsValid {
			// println("[DEBUG] Access blocked for method:", info.FullMethod)
			return nil, status.Error(codes.PermissionDenied, "blocked")
		}
		return handler(ctx, req)
	}
}

func MasterInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// VuKhiService
		if expectedSlug, ok := VuKhiMethodToSlug[info.FullMethod]; ok {
			return VuKhiAuthMaskInterceptorWithSlug(expectedSlug)(ctx, req, info, handler)
		}

		// HeService
		if expectedSlug, ok := HeMethodToSlug[info.FullMethod]; ok {
			return HeAuthMaskInterceptorWithSlug(expectedSlug)(ctx, req, info, handler)
		}

		// LoaiVuKhiService
		if expectedSlug, ok := LoaiVuKhiMethodToSlug[info.FullMethod]; ok {
			return LoaiVuKhiAuthMaskInterceptorWithSlug(expectedSlug)(ctx, req, info, handler)
		}

		// DoHiemService
		if expectedSlug, ok := DoHiemMethodToSlug[info.FullMethod]; ok {
			return DoHiemAuthMaskInterceptorWithSlug(expectedSlug)(ctx, req, info, handler)
		}
		// Các service khác
		return CommonAuthInterceptor()(ctx, req, info, handler)
	}
}
