// pkg/logx/kafka_log.go
package logx

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

func KafkaGhiLogProduceOK(ctx context.Context, elog *LoggerElastic, topic, key string, valueLen int, waitMs, ackMs int64) {
	if elog == nil {
		return
	}
	elog.Info(ctx, "Kafka produce ok", map[string]any{
		"event.dataset":     "kafka",
		"service.component": "kafka-producer",
		"kafka.topic":       topic,
		"kafka.key":         key,
		"kafka.msg_size":    valueLen,
		"kafka.wait_ms":     waitMs,
		"kafka.ack_ms":      ackMs,
		"duration_ms":       ackMs,
	})
}

func KafkaGhiLogProduceErr(ctx context.Context, elog *LoggerElastic, topic, key string, valueLen int, waitMs int64, err error) {
	if elog == nil {
		return
	}
	elog.Error(ctx, "Kafka produce failed", map[string]any{
		"event.dataset":     "kafka",
		"service.component": "kafka-producer",
		"kafka.topic":       topic,
		"kafka.key":         key,
		"kafka.msg_size":    valueLen,
		"kafka.wait_ms":     waitMs,
		"duration_ms":       0,
		"error.type":        "Kafka.Produce",
		"error.message":     err.Error(),
	})
}

func KafkaGhiLogConsumeOK(ctx context.Context, l *LoggerElastic, group string, m kafka.Message, procTook time.Duration) {
	if l == nil {
		return
	}
	kv := map[string]any{
		"event.dataset":        "kafka",
		"service.component":    "kafka-consumer",
		"kafka.consumer_group": group,
		"kafka.topic":          m.Topic,
		"kafka.partition":      m.Partition,
		"kafka.offset":         m.Offset,
		"kafka.key.size":       len(m.Key),
		"kafka.value.size":     len(m.Value),
		"duration_ms":          procTook.Milliseconds(), // thời gian xử lý nghiệp vụ
	}
	// nếu bạn nhét trace_id/request_id vào header message thì kéo ra luôn
	for _, h := range m.Headers {
		if h.Key == "trace_id" {
			kv["trace_id"] = string(h.Value)
		}
		if h.Key == "request_id" {
			kv["request_id"] = string(h.Value)
		}
	}
	l.Ghi(ctx, "info", "Kafka consumed message", kv)
}

func KafkaGhiLogConsumeErr(ctx context.Context, l *LoggerElastic, group, topic string, err error) {
	if l == nil {
		return
	}
	l.Ghi(ctx, "warn", "Kafka consume failed", map[string]any{
		"event.dataset":        "kafka",
		"service.component":    "kafka-consumer",
		"kafka.consumer_group": group,
		"kafka.topic":          topic,
		"error.message":        err.Error(),
		"error.type":           "Kafka.Consume",
	})
}
