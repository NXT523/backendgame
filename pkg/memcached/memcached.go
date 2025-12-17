package memcached

import (
	"context"
	"strings"
	"time"

	"game/internal/config"

	"github.com/bradfitz/gomemcache/memcache"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Tạo client Memcached
func NewMemcacheClient(_ context.Context, cfg config.CauHinh) (*memcache.Client, error) {
	addr := strings.TrimSpace(cfg.MemcachedAddr)
	if addr == "" {
		addr = "127.0.0.1:11211"
	}
	mc := memcache.New(addr)

	if cfg.MemcachedTimeout > 0 {
		mc.Timeout = cfg.MemcachedTimeout
	}
	if cfg.MemcachedMaxIdleConns > 0 {
		mc.MaxIdleConns = cfg.MemcachedMaxIdleConns
	}

	// Healthcheck (gomemcache không có Ping)
	key := cfg.MemcachedHealthKey
	if key == "" {
		key = "healthcheck"
	}
	ttl := int32(1)
	if cfg.MemcachedHealthTTL > 0 {
		ttl = int32(cfg.MemcachedHealthTTL / time.Second)
	}
	if err := mc.Set(&memcache.Item{Key: key, Value: []byte("ok"), Expiration: ttl}); err != nil {
		return nil, err
	}
	if _, err := mc.Get(key); err != nil {
		return nil, err
	}
	return mc, nil
}

// gomemcache không có Close(), nên để trống cho đồng bộ interface
func CloseMemcache(_ context.Context, _ *memcache.Client) error {
	return nil
}

func LayTrangCacheJSONMem[T proto.Message](mc *memcache.Client, key string, newT func() T) (T, bool) {
	var zero T
	if mc == nil {
		return zero, false
	}
	it, err := mc.Get(key)
	if err != nil || len(it.Value) == 0 {
		return zero, false
	}
	bs := it.Value
	if len(bs) == 0 || (bs[0] != '{' && bs[0] != '[') {
		return zero, false
	}

	v := newT()
	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(bs, v); err != nil {
		return zero, false
	}
	return v, true
}

func LuuTrangCacheJSONMem(mc *memcache.Client, key string, msg proto.Message, ttl time.Duration) {
	if mc == nil || msg == nil {
		return
	}
	bs, err := (protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}).Marshal(msg)
	if err != nil {
		return
	}

	exp := int32(0)
	if ttl > 0 {
		exp = int32(ttl / time.Second)
		if exp <= 0 {
			exp = 1
		}
	}

	_ = mc.Set(&memcache.Item{
		Key:        key,
		Value:      bs,
		Expiration: exp,
	})
}
