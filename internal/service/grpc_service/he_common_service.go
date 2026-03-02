package grpc_service

import (
	"context"
	"log"
)

// invalidateHeCache thực hiện xóa hoặc làm mới các key cache liên quan đến Hệ trong Redis.
// Hàm này thường được gọi sau các thao tác ghi (Create, Update, Delete) để đảm bảo tính nhất quán dữ liệu.
// Sử dụng Pipeline để tối ưu hóa việc gửi nhiều lệnh đến Redis cùng lúc.
func (s *HeService) invalidateHeCache(ctx context.Context) {
	if s == nil || s.rdb == nil {
		return
	}

	pipe := s.rdb.Pipeline()
	pipe.Incr(ctx, "he:search:version")
	pipe.Del(ctx,
		"he:get_all_he",
		"he:get_all_ten_he",
		"he:get_all_mo_ta_he")

	if _, err := pipe.Exec(ctx); err != nil {
		log.Println("REDIS INVALIDATE ERROR:", err)
	}
}
