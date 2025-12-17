// kafka/consumer_new.go
package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"game/internal/config"
	"game/internal/events"
	"game/pkg/logx"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

// ----- Reader helper -----
func taoReader(topic, group string) *kafka.Reader {
	cfg := config.DocCauHinh()
	brokers := tachBrokers(cfg.KafkaBrokers)
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  group,
		MinBytes: 1,
		MaxBytes: 10e6,
		MaxWait:  3 * time.Second,
	})
}

// gom meta Kafka chung cho 1 log
func kafkaKV(group string, m kafka.Message) map[string]any {
	kv := map[string]any{
		"event.dataset":        "kafka",
		"service.component":    "kafka-consumer",
		"kafka.consumer_group": group,
		"kafka.topic":          m.Topic,
		"kafka.partition":      m.Partition,
		"kafka.offset":         m.Offset,
		"kafka.timestamp":      m.Time.UTC().Format(time.RFC3339Nano),
		"kafka.msg_size":       len(m.Value),
	}
	if len(m.Key) > 0 {
		kv["kafka.key"] = string(m.Key)
	}
	return kv
}

// ----- CONSUMER 1: XÓA CACHE SEARCH -----
func ChayConsumerXoaCacheVuKhi(ctx context.Context, rdb *redis.Client, elog *logx.LoggerElastic) {
	const group = "vukhi-cache-service"
	cfg := config.DocCauHinh()
	r := taoReader(cfg.KafkaTopic, group)
	defer r.Close()

	elog.Info(ctx, "Consumer started", map[string]any{
		"event.dataset":        "kafka",
		"service.component":    "kafka-consumer",
		"kafka.consumer_group": group,
		"kafka.topic":          cfg.KafkaTopic,
	})

	for {
		tWait := time.Now()
		m, err := r.FetchMessage(ctx)
		waitMs := time.Since(tWait).Milliseconds()
		if err != nil {
			logx.KafkaGhiLogConsumeErr(ctx, elog, group, cfg.KafkaTopic, err)
			time.Sleep(2 * time.Second)
			continue
		}

		tProc := time.Now()

		var ev events.SuKienVuKhiCreated
		_ = json.Unmarshal(m.Value, &ev)

		pattern := "vukhi:search:v2:*"
		deleted, delErr := xoaTheoPattern(ctx, rdb, pattern)

		kv := kafkaKV(group, m)
		procMs := time.Since(tProc).Milliseconds()
		kv["kafka.wait_ms"] = waitMs
		kv["proc_ms"] = procMs
		kv["duration_ms"] = procMs
		kv["cache.pattern"] = pattern
		if ev.ID != 0 {
			kv["vukhi.id"] = ev.ID
		}

		if delErr != nil {
			kv["error.type"] = "Redis.Del"
			kv["error.message"] = delErr.Error()
			elog.Error(ctx, "Xóa cache cũ thất bại", kv)
			continue
		}

		kv["cache.deleted"] = deleted
		elog.Info(ctx, "Xóa cache cũ thành công", kv)

		// Commit message sau khi xử lý xong
		if err := r.CommitMessages(ctx, m); err != nil {
			elog.Error(ctx, "Commit Kafka message thất bại", map[string]any{
				"error":  err.Error(),
				"offset": m.Offset,
			})
		}
	}
}

// xóa theo pattern bằng SCAN
func xoaTheoPattern(ctx context.Context, rdb *redis.Client, pattern string) (int, error) {
	var cursor uint64
	var n int
	for {
		keys, next, err := rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return n, err
		}
		cursor = next

		if len(keys) > 0 {
			if err := rdb.Del(ctx, keys...).Err(); err != nil {
				return n, err
			}
			n += len(keys)
		}

		if cursor == 0 {
			break
		}
	}
	return n, nil
}

// ----- CONSUMER 2: AUDIT LOG -----
func ChayConsumerAuditVuKhi(ctx context.Context, elog *logx.LoggerElastic) {
	const group = "vukhi-audit-service"
	cfg := config.DocCauHinh()
	r := taoReader(cfg.KafkaTopic, group)
	defer r.Close()

	elog.Info(ctx, "Consumer started", map[string]any{
		"event.dataset":        "kafka",
		"service.component":    "kafka-consumer",
		"kafka.consumer_group": group,
		"kafka.topic":          cfg.KafkaTopic,
	})

	for {
		tWait := time.Now()
		m, err := r.FetchMessage(ctx)
		waitMs := time.Since(tWait).Milliseconds()
		if err != nil {
			logx.KafkaGhiLogConsumeErr(ctx, elog, group, cfg.KafkaTopic, err)
			time.Sleep(2 * time.Second)
			continue
		}

		tProc := time.Now()
		var ev events.SuKienVuKhiCreated
		_ = json.Unmarshal(m.Value, &ev)

		kv := kafkaKV(group, m)
		procMs := time.Since(tProc).Milliseconds()
		kv["kafka.wait_ms"] = waitMs
		kv["proc_ms"] = procMs
		kv["duration_ms"] = procMs

		if ev.ID != 0 {
			kv["vukhi.id"] = ev.ID
		}
		if strings.TrimSpace(ev.TenVuKhi) != "" {
			kv["vukhi.ten"] = ev.TenVuKhi
		}
		if ev.MaDoHiem != 0 {
			kv["vukhi.do_hiem"] = ev.MaDoHiem
		}

		elog.Info(ctx, "Ghi log audit sự kiện vu_khi_created", kv)

		// Commit message
		if err := r.CommitMessages(ctx, m); err != nil {
			elog.Error(ctx, "Commit Kafka message thất bại", map[string]any{
				"error":  err.Error(),
				"offset": m.Offset,
			})
		}
	}
}

// ----- CONSUMER 3: ĐỒNG BỘ / WARM CACHE -----
func ChayConsumerDongBoVuKhi(ctx context.Context, rdb *redis.Client, elog *logx.LoggerElastic) {
	const group = "vukhi-dongbo-service"
	cfg := config.DocCauHinh()
	r := taoReader(cfg.KafkaTopic, group)
	defer r.Close()

	elog.Info(ctx, "Consumer started", map[string]any{
		"event.dataset":        "kafka",
		"service.component":    "kafka-consumer",
		"kafka.consumer_group": group,
		"kafka.topic":          cfg.KafkaTopic,
	})

	for {
		tWait := time.Now()
		m, err := r.FetchMessage(ctx)
		waitMs := time.Since(tWait).Milliseconds()
		if err != nil {
			logx.KafkaGhiLogConsumeErr(ctx, elog, group, cfg.KafkaTopic, err)
			time.Sleep(2 * time.Second)
			continue
		}

		tProc := time.Now()
		var ev events.SuKienVuKhiCreated
		_ = json.Unmarshal(m.Value, &ev)

		key := "vukhi:last_created"
		val := fmt.Sprintf("%d|%s|%d", ev.ID, ev.TenVuKhi, ev.CreatedAt)

		kv := kafkaKV(group, m)
		procMs := time.Since(tProc).Milliseconds()
		kv["kafka.wait_ms"] = waitMs
		kv["proc_ms"] = procMs
		kv["duration_ms"] = procMs
		kv["redis.key"] = key
		kv["redis.value"] = val
		if ev.ID != 0 {
			kv["vukhi.id"] = ev.ID
		}

		if err := rdb.Set(ctx, key, val, 10*time.Minute).Err(); err != nil {
			kv["error.type"] = "Redis.Set"
			kv["error.message"] = err.Error()
			elog.Error(ctx, "cache cập nhật thất bại", kv)
			continue
		}

		elog.Info(ctx, "cache đã được cập nhật", kv)

		// Commit message
		if err := r.CommitMessages(ctx, m); err != nil {
			elog.Error(ctx, "Commit Kafka message thất bại", map[string]any{
				"error":  err.Error(),
				"offset": m.Offset,
			})
		}
	}
}
