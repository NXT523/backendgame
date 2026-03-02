// kafka/consumer_new.go
package kafka

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"game/internal/utils"
	"game/internal/events"
	"game/pkg/logx"

	"github.com/segmentio/kafka-go"
)

// ----- Reader helper -----
func taoReader(topic, group string) *kafka.Reader {
	brokers := utils.GetenvStrings("KAFKA_BROKERS", nil)
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

func ChayConsumerAuditVuKhi(ctx context.Context, elog *logx.LoggerElastic) {
	const group = "vukhi-audit-service"
	r := taoReader(utils.GetenvString("KAFKA_TOPIC_CREATE", ""), group)
	defer r.Close()

	elog.Info(ctx, "Consumer started", map[string]any{
		"event.dataset":        "kafka",
		"service.component":    "kafka-consumer",
		"kafka.consumer_group": group,
		"kafka.topic":          utils.GetenvString("KAFKA_TOPIC_CREATE", ""),
	})

	for {
		tWait := time.Now()
		m, err := r.FetchMessage(ctx)
		waitMs := time.Since(tWait).Milliseconds()
		if err != nil {
			logx.KafkaGhiLogConsumeErr(ctx, elog, group, utils.GetenvString("KAFKA_TOPIC_CREATE", ""), err)
			time.Sleep(2 * time.Second)
			continue
		}

		tProc := time.Now()
		var ev events.SuKienVuKhiCreated
		if err := json.Unmarshal(m.Value, &ev); err != nil {
			elog.Error(ctx, "Unmarshal message thất bại", map[string]any{
				"error":  err.Error(),
				"offset": m.Offset,
			})
			_ = r.CommitMessages(ctx, m)
			continue
		}

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

		if err := r.CommitMessages(ctx, m); err != nil {
			elog.Error(ctx, "Commit Kafka message thất bại", map[string]any{
				"error":  err.Error(),
				"offset": m.Offset,
			})
		}
	}
}
