package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"game/internal/events"
	"game/pkg/logx"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

// SenderWorker đọc từ Redis outbox (BLPOP) và gửi tới Kafka.
// Worker sẽ return khi ctx cancelled.
func SenderWorker(ctx context.Context, rdb *redis.Client, writer *kafka.Writer, elog *logx.LoggerElastic) {
	const outboxKey = "outbox:vukhi"
	const dlqKey = "outbox:vukhi:dlq"

	// Log khởi động bằng elog (bình thường). Nếu elog gặp lỗi, những thông báo shutdown sẽ dùng log.Printf.
	elog.Info(ctx, "SenderWorker đã bắt đầu", map[string]any{"outbox": outboxKey})

	for {
		// BLPop sẽ unblock khi ctx bị cancel
		res, err := rdb.BLPop(ctx, 0*time.Second, outboxKey).Result()
		if err != nil {
			// Nếu ctx bị hủy -> thoát ngay, dùng standard logger để tránh block nếu elog gặp vấn đề
			if ctx.Err() != nil {
				log.Println("SenderWorker: ngữ cảnh đã bị hủy, đang thoát")
				return
			}
			// Lỗi Redis khác: log qua elog với field error.message
			elog.Error(ctx, "Lỗi BLPop", map[string]any{"error.message": err.Error()})
			time.Sleep(time.Second)
			continue
		}
		if len(res) < 2 {
			elog.Error(ctx, "BLPop kết quả bất ngờ", map[string]any{"res": res})
			continue
		}

		raw := []byte(res[1])

		// Try to unmarshal to extract ID (để đặt Key), fallback empty key
		var ev events.SuKienVuKhiCreated
		var key []byte
		if err := json.Unmarshal(raw, &ev); err == nil {
			key = []byte(fmt.Sprintf("%d", ev.ID))
		} else {
			// Unmarshal lỗi -> log và tiếp tục gửi với empty key
			elog.Error(ctx, "Mục hộp thư đi không sắp xếp không thành công, gửi không có khóa", map[string]any{"error.message": err.Error()})
		}

		// retry loop with exponential backoff
		maxRetry := 10
		backoff := time.Second
		var errKafka error
		for i := 0; i < maxRetry; i++ {
			ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
			errKafka = writer.WriteMessages(ctx2, kafka.Message{
				Key:   key,
				Value: raw,
			})
			cancel()

			if errKafka == nil {
				elog.Info(ctx, "Đã gửi mục hộp thư đi tới kafka", map[string]any{"attempt": i + 1, "id": ev.ID})
				break
			}

			// if ctx cancelled, break (dùng standard logger để không phụ thuộc elog)
			if ctx.Err() != nil {
				log.Println("SenderWorker: ngữ cảnh bị hủy trong quá trình gửi:", ctx.Err())
				break
			}

			elog.Error(ctx, "Gửi Kafka không thành công, đang thử lại", map[string]any{"attempt": i + 1, "error.message": errKafka.Error()})
			time.Sleep(backoff)
			backoff *= 2
		}

		if errKafka != nil {
			elog.Error(ctx, "Không gửi được hộp thư đi sau khi thử lại, chuyển sang DLQ", map[string]any{"error.message": errKafka.Error(), "id": ev.ID})
			if err := rdb.RPush(ctx, dlqKey, raw).Err(); err != nil {
				elog.Error(ctx, "Push to DLQ failed", map[string]any{"error.message": err.Error()})
			}
			// don't requeue automatically to avoid tight loops
		}
	}
}
