package memcached

import (
	"context"
	"strings"
	"time"

	"game/internal/utils"

	"github.com/bradfitz/gomemcache/memcache"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Tạo client Memcached
func NewMemcacheClient(_ context.Context) (*memcache.Client, error) {
	addr := strings.TrimSpace(utils.GetenvString("MEMCACHED_ADDR", "127.0.0.1:11211"))
	mc := memcache.New(addr)
	timeout := utils.GetenvDuration("MEMCACHED_TIMEOUT", 200*time.Millisecond)
	if timeout > 0 {
		mc.Timeout = timeout
	}
	maxIdleConns := utils.GetenvInt("MEMCACHED_MAX_IDLE_CONNS", 100)
	if maxIdleConns > 0 {
		mc.MaxIdleConns = maxIdleConns
	}

	key := utils.GetenvString("MEMCACHED_HEALTH_KEY", "healthcheck")
	ttlDur := utils.GetenvDuration("MEMCACHED_HEALTH_TTL", 1*time.Second)

	ttl := int32(ttlDur / time.Second)

	if err := mc.Set(&memcache.Item{
		Key:        key,
		Value:      []byte("ok"),
		Expiration: ttl,
	}); err != nil {
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
