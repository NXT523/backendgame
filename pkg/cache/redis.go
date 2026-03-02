package cache

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type redisCacheService struct {
	rdb *redis.Client
}

func NewRedisCacheService(rdb *redis.Client) RedisCacheService {
	return &redisCacheService{
		rdb: rdb,
	}
}

func (cs *redisCacheService) Get(
	ctx context.Context,
	key string,
	dest any,
) error {
	data, err := cs.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (cs *redisCacheService) Set(
	ctx context.Context,
	key string,
	value any,
	ttl time.Duration,
) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return cs.rdb.Set(ctx, key, data, ttl).Err()
}

func (cs *redisCacheService) Exists(
	ctx context.Context,
	key string,
) (bool, error) {
	n, err := cs.rdb.Exists(ctx, key).Result()
	return n > 0, err
}

func (cs *redisCacheService) Clear(
	ctx context.Context,
	pattern string,
) error {
	var cursor uint64
	for {
		keys, next, err := cs.rdb.Scan(ctx, cursor, pattern, 50).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			_ = cs.rdb.Del(ctx, keys...).Err()
		}
		cursor = next
		if cursor == 0 {
			return nil
		}
	}
}

func CloseRedis(_ context.Context, rdb *redis.Client) error {
	if rdb == nil {
		return nil
	}
	return rdb.Close()
}

// SetJSON lưu protobuf message vào Redis dưới dạng JSON.
func (cs *redisCacheService) SetJSON(ctx context.Context, key string, msg proto.Message, ttl time.Duration) error {
	if msg == nil {
		return nil
	}

	bs, err := (protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: false,
	}).Marshal(msg)
	if err != nil {
		return err
	}

	return cs.rdb.Set(ctx, key, bs, ttl).Err()
}

// GetJSON lấy dữ liệu JSON từ Redis và unmarshal vào protobuf message mới.
func (cs *redisCacheService) GetJSON(ctx context.Context, key string, newT func() proto.Message) (proto.Message, bool) {

	bs, err := cs.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}

	bs = bytes.TrimSpace(bs)
	if len(bs) == 0 || (bs[0] != '{' && bs[0] != '[') {
		return nil, false
	}

	v := newT()
	if err := (protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}).Unmarshal(bs, v); err != nil {
		return nil, false
	}

	return v, true
}

// SetProto lưu protobuf message vào Redis dưới dạng binary (proto.Marshal).
// Cách này nhanh và gọn hơn JSON, phù hợp cache nội bộ backend.
func (cs *redisCacheService) SetProto(ctx context.Context, key string, msg proto.Message, ttl time.Duration) error {
	bs, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	return cs.rdb.Set(ctx, key, bs, ttl).Err()
}

// GetProto lấy protobuf binary từ Redis và unmarshal vào message mới.
// Trả về (message, true) nếu cache hit, ngược lại (nil, false).
func (cs *redisCacheService) GetProto(ctx context.Context, key string, newT func() proto.Message) (proto.Message, bool) {
	bs, err := cs.rdb.Get(ctx, key).Bytes()
	if err != nil || len(bs) == 0 {
		return nil, false
	}

	v := newT()
	if err := proto.Unmarshal(bs, v); err != nil {
		return nil, false
	}

	return v, true
}
