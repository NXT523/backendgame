package kafka

import (
	"context"
	"encoding/json"
	"time"

	"game/ent"
	"game/internal/events"
	"game/pkg/logx"
)

func ChayConsumerVuKhiUpdate(ctx context.Context, entClient *ent.Client, elog *logx.LoggerElastic) {
	const (
		group = "vukhi-update-group"
		topic = "vu_khi_updated"
	)

	r := taoReader(topic, group)
	defer r.Close()

	for {
		// 1️⃣ Kafka fetch theo OFFSET → đảm bảo thứ tự
		m, err := r.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			time.Sleep(time.Second)
			continue
		}

		// 2️⃣ Parse message
		var ev events.SuKienVuKhiUpdated
		if err := json.Unmarshal(m.Value, &ev); err != nil {
			_ = r.CommitMessages(ctx, m) // message hỏng → bỏ
			continue
		}

		// 4️⃣ Check giá trị trước khi update
		if ev.SatThuongCoBan < 60 {
			// Bỏ qua message, vẫn commit để Kafka không retry
			_ = r.CommitMessages(ctx, m)
			continue
		}

		// 3️⃣ Update DB
		_, err = entClient.VuKhi.
			UpdateOneID(ev.MaVuKhi).
			SetSatThuongCoBan(ev.SatThuongCoBan).
			SetTocDoDanh(ev.TocDoDanh).
			SetTamDanh(ev.TamDanh).
			SetMoTa(ev.MoTa).
			SetMaLoaiVuKhi(ev.MaLoai).
			SetMaDoHiem(ev.MaDoHiem).
			SetMaHe(ev.MaHe).
			Save(ctx)

		if err != nil {
			// không commit → Kafka thử lại → vẫn giữ order
			elog.Error(ctx, "Update DB lỗi", map[string]any{
				"id":     ev.MaVuKhi,
				"offset": m.Offset,
				"error":  err.Error(),
			})
			continue
		}

		// 4️⃣ Commit SAU khi update xong
		_ = r.CommitMessages(ctx, m)
	}
}
