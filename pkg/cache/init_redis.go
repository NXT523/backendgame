package cache

import (
	"context"
	"crypto/tls"
	"game/internal/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:         utils.GetenvString("REDIS_ADDR", "localhost:6379"),
		Username:     utils.GetenvString("REDIS_USERNAME", ""),
		Password:     utils.GetenvString("REDIS_PASSWORD", ""),
		DB:           utils.GetenvInt("REDIS_DB", 0),
		Protocol:     utils.GetenvInt("REDIS_PROTOCOL", 2),
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	if utils.GetenvBool("REDIS_TLS", false) {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: utils.GetenvBool("REDIS_TLS_SKIP_VERIFY", false),
		}
	}

	rdb := redis.NewClient(opts)

	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, err
	}

	return rdb, nil
}
