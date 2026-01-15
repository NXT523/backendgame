package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"game/ent"
	"game/ent/vukhi"
	v1 "game/v1"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

func StartVuKhiWorker(entClient *ent.Client, rdb *redis.Client, kafkaReader *kafka.Reader) {
	for {
		msg, err := kafkaReader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Lỗi đọc message:", err)
			continue
		}

		var req v1.CapNhatTheoTenVuKhiRequest
		if err := json.Unmarshal(msg.Value, &req); err != nil {
			fmt.Println("Unmarshal lỗi:", err)
			continue
		}

		processUpdate(req, entClient, rdb)
	}
}
func processUpdate(req v1.CapNhatTheoTenVuKhiRequest, entClient *ent.Client, rdb *redis.Client) {
	ctx := context.Background()
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return
	}

	lockKey := "lock:vukhi:" + name
	lockTTL := 5 * time.Second
	ok, err := rdb.SetNX(ctx, lockKey, "1", lockTTL).Result()
	if err != nil || !ok {
		fmt.Println("Vũ khí đang được cập nhật, retry sau")
		return
	}
	defer rdb.Del(ctx, lockKey)

	tx, err := entClient.Tx(ctx)
	if err != nil {
		fmt.Println("Tx lỗi:", err)
		return
	}

	cur, err := tx.VuKhi.Query().Where(vukhi.TenVuKhiEQ(name)).Only(ctx)
	if err != nil {
		tx.Rollback()
		fmt.Println("Không tìm thấy vũ khí:", err)
		return
	}

	_, err = tx.VuKhi.UpdateOneID(cur.ID).
		SetTenVuKhi(req.TenVuKhi).
		SetSatThuongCoBan(int(req.SatThuongCoBan)).
		SetTocDoDanh(req.TocDoDanh).
		SetTamDanh(int(req.TamDanh)).
		SetMoTa(req.MoTa).
		SetMaLoai(int(req.MaLoai)).
		SetMaDoHiem(int(req.MaDoHiem)).
		SetMaHe(int(req.MaHe)).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		fmt.Println("Update lỗi:", err)
		return
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("Commit lỗi:", err)
		return
	}

	fmt.Println("Update thành công:", req.Name)
}
