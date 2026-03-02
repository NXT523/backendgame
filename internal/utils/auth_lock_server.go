package utils

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	kLockDefaultTTL = 5 * time.Second
	kLockRetryCount = 1
	kLockSleepTime  = 100 * time.Millisecond
)

type RedisLock struct {
	rds *redis.Client
}

func NewRedisLock(rds *redis.Client) *RedisLock {
	return &RedisLock{rds: rds}
}

// WithDistributedLock
// - Distributed lock bằng Redis
// - Có fencing token
// - Có watchdog gia hạn TTL
// - Cleanup an toàn kể cả khi ctx bị cancel
func (l *RedisLock) WithDistributedLock(
	ctx context.Context,
	key string,
	ttl time.Duration,
	action func(fencingToken string) error,
) error {

	if l == nil || l.rds == nil {
		return status.Error(codes.Internal, "RedisLock chưa được khởi tạo")
	}

	if ttl <= 0 {
		ttl = kLockDefaultTTL
	}

	// 1. Tạo token định danh duy nhất
	token := uuid.New().String()

	// 2. Acquire lock (retry)
	acquired := false
	for i := 0; i < kLockRetryCount; i++ {
		ok, err := l.rds.SetNX(ctx, key, token, ttl).Result()
		if err != nil {
			return status.Errorf(codes.Internal, "Redis error: %v", err)
		}
		if ok {
			acquired = true
			break
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(kLockSleepTime):
		}
	}

	if !acquired {
		return status.Error(codes.ResourceExhausted, "Hệ thống đang bận, vui lòng thử lại sau")
	}

	// 3. Watchdog gia hạn TTL (context độc lập)
	watchdogCtx, cancelWatchdog := context.WithCancel(context.Background())
	doneCh := make(chan struct{})

	go func() {
		defer close(doneCh)

		ticker := time.NewTicker(ttl / 3)
		defer ticker.Stop()

		refreshScript := `
			if redis.call("GET", KEYS[1]) == ARGV[1] then
				return redis.call("PEXPIRE", KEYS[1], ARGV[2])
			else
				return 0
			end
		`

		for {
			select {
			case <-ticker.C:
				_, _ = l.rds.Eval(
					watchdogCtx,
					refreshScript,
					[]string{key},
					token,
					int(ttl/time.Millisecond),
				).Result()
			case <-watchdogCtx.Done():
				return
			}
		}
	}()

	// 4. Cleanup an toàn
	defer func() {
		// Stop watchdog
		cancelWatchdog()
		<-doneCh

		// Xóa lock bằng context riêng
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		deleteScript := `
			if redis.call("GET", KEYS[1]) == ARGV[1] then
				return redis.call("DEL", KEYS[1])
			else
				return 0
			end
		`
		_, _ = l.rds.Eval(cleanupCtx, deleteScript, []string{key}, token).Result()
	}()

	// 5. Thực thi logic chính (truyền fencing token)
	return action(token)
}
