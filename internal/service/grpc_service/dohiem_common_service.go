package grpc_service

import (
	"context"
	"log"
)

// invalidateDoHiemCache thực hiện xóa hoặc làm mới các key cache liên quan đến độ hiếm trong Redis.
// Hàm này thường được gọi sau các thao tác ghi (Create, Update, Delete) để đảm bảo tính nhất quán dữ liệu.
// Sử dụng Pipeline để tối ưu hóa việc gửi nhiều lệnh đến Redis cùng lúc.
func (s *DoHiemService) invalidateDoHiemCache(ctx context.Context) {
	if s == nil || s.rdb == nil {
		return
	}

	pipe := s.rdb.Pipeline()
	pipe.Incr(ctx, "dohiem:search:version")
	pipe.Del(ctx,
		"dohiem:get_all_dohiem",
		"dohiem:get_all_ten_dohiem",
		"dohiem:get_all_so_luong_dohiem",
		"dohiem:get_all_mau_sac_dohiem",
		"dohiem:get_all_sat_thuong_bonus_dohiem",
		"dohiem:get_all_toc_do_danh_bonus_dohiem",
		"dohiem:get_all_cap_bac_dohiem")

	if _, err := pipe.Exec(ctx); err != nil {
		log.Println("REDIS INVALIDATE ERROR:", err)
	}
}
