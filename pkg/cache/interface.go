package cache

import (
	"context"
	"time"

	"google.golang.org/protobuf/proto"
)

type RedisCacheService interface {
	Get(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, key string, value any, ttl time.Duration) error

	GetJSON(ctx context.Context, key string, newT func() proto.Message) (proto.Message, bool)
	SetJSON(ctx context.Context, key string, msg proto.Message, ttl time.Duration) error

	// SetProto lưu protobuf message vào Redis dưới dạng binary (proto.Marshal).
	//
	// Cách này nhanh và gọn hơn JSON, phù hợp cache nội bộ backend.
	SetProto(ctx context.Context, key string, msg proto.Message, ttl time.Duration) error

	// GetProto lấy protobuf binary từ Redis và unmarshal vào message mới.
	//
	// Trả về (message, true) nếu cache hit, ngược lại (nil, false).
	GetProto(ctx context.Context, key string, newT func() proto.Message) (proto.Message, bool)

	Clear(ctx context.Context, pattern string) error
	Exists(ctx context.Context, key string) (bool, error)
}
