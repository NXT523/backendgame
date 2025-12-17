package redisx

import (
	"bytes"
	"context"
	"crypto/tls"
	"time"

	"game/internal/config"

	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func NewRedisClient(ctx context.Context, cfg config.CauHinh) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:         cfg.RedisAddr,
		Username:     cfg.RedisUsername,
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDB,
		Protocol:     cfg.RedisProtocol,
		DialTimeout:  cfg.RedisDialTimeout,
		ReadTimeout:  cfg.RedisReadTimeout,
		WriteTimeout: cfg.RedisWriteTimeout,
	}

	if cfg.RedisTLS {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: cfg.RedisTLSSkipVerify,
		}
	}

	rdb := redis.NewClient(opts)

	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, err
	}
	return rdb, nil
}

func CloseRedis(_ context.Context, rdb *redis.Client) error {
	if rdb == nil {
		return nil
	}
	return rdb.Close()
}

func LayTrangCacheJSON[T proto.Message](ctx context.Context, rdb *redis.Client, key string, newT func() T) (T, bool) {
	var zero T
	if rdb == nil {
		return zero, false
	}
	bs, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		// Key không tồn tại
		if err == redis.Nil {
			return zero, false
		}
		// WRONGTYPE hoặc lỗi Redis khác -> coi như miss cache
		return zero, false
	}

	// Guard nhanh + TrimSpace để chắc chắn là JSON
	bs = bytes.TrimSpace(bs)
	if len(bs) == 0 || (bs[0] != '{' && bs[0] != '[') {
		return zero, false
	}

	v := newT()
	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(bs, v); err != nil {
		return zero, false
	}
	return v, true
}

func LuuTrangCacheJSON(ctx context.Context, rdb *redis.Client, key string, msg proto.Message, ttl time.Duration) {
	if rdb == nil || msg == nil {
		return
	}
	bs, err := (protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}).Marshal(msg)
	if err != nil {
		return
	}
	_ = rdb.Set(ctx, key, bs, ttl).Err()
}
