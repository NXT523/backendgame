package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"game/internal/events"
	"game/pkg/logx"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type EventType string

const (
	EventVuKhiCreated EventType = "vukhi.created"
	EventVuKhiUpdated EventType = "vukhi.updated"
)

func SenderWorker(
	workerCtx context.Context,
	rdb *redis.Client,
	eventType EventType,
	outboxKey string,
	writer *kafka.Writer,
	elog *logx.LoggerElastic,
) {
	processingKey := outboxKey + ":processing"

	for {
		select {
		case <-workerCtx.Done():
			return
		default:
		}

		// 1️⃣ Lấy message từ đầu outbox (block 5s nếu rỗng)
		res, err := rdb.BLPop(context.Background(), 5*time.Second, outboxKey).Result()
		if err == redis.Nil {
			continue // không có message, loop lại
		}
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		raw := res[1] // BLPop trả về [key, value]

		// 2️⃣ Push vào cuối processing queue để giữ thứ tự
		if _, err := rdb.RPush(context.Background(), processingKey, raw).Result(); err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		var key []byte
		var id string

		// 3️⃣ Parse message
		switch eventType {
		case EventVuKhiCreated:
			var ev events.SuKienVuKhiCreated
			if err := json.Unmarshal([]byte(raw), &ev); err != nil {
				rdb.LRem(context.Background(), processingKey, 1, raw)
				continue
			}
			id = fmt.Sprintf("%d", ev.ID)
			key = []byte(id)

		case EventVuKhiUpdated:
			var ev events.SuKienVuKhiUpdated
			if err := json.Unmarshal([]byte(raw), &ev); err != nil {
				rdb.LRem(context.Background(), processingKey, 1, raw)
				continue
			}
			id = fmt.Sprintf("%d", ev.MaVuKhi)
			key = []byte(id)
		}

		// 4️⃣ Gửi Kafka
		kctx, cancel := context.WithTimeout(workerCtx, 5*time.Second)
		err = writer.WriteMessages(kctx, kafka.Message{
			Key:   key,
			Value: []byte(raw),
		})
		cancel()

		if err != nil {
			// 5️⃣ Nếu fail → requeue trở lại outbox
			pipe := rdb.TxPipeline()
			pipe.LRem(context.Background(), processingKey, 1, raw)
			pipe.RPush(context.Background(), outboxKey, raw)
			_, _ = pipe.Exec(context.Background())

			elog.Error(workerCtx, "Kafka thất bại, requeue tin nhắn", map[string]any{
				"eventType": eventType,
				"id":        id,
				"error":     err.Error(),
				"value":     raw,
			})
			time.Sleep(2 * time.Second)
			continue
		}

		// 6️⃣ Kafka thành công → remove khỏi processing
		rdb.LRem(context.Background(), processingKey, 1, raw)

		elog.Info(workerCtx, "Kafka sent", map[string]any{
			"eventType": eventType,
			"id":        id,
		})
	}
}
