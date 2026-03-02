package grpc_service

import (
	"context"
	"log"
)

// invalidateLoaiVuKhiCache thực hiện xóa hoặc làm mới các key cache liên quan đến loại vũ khí trong Redis.
// Hàm này thường được gọi sau các thao tác ghi (Create, Update, Delete) để đảm bảo tính nhất quán dữ liệu.
// Sử dụng Pipeline để tối ưu hóa việc gửi nhiều lệnh đến Redis cùng lúc.
func (s *LoaiVuKhiService) invalidateLoaiVuKhiCache(ctx context.Context) {
	if s == nil || s.rdb == nil {
		return
	}

	pipe := s.rdb.Pipeline()
	pipe.Incr(ctx, "loaivukhi:search:version")
	pipe.Del(ctx, 
		"loaivukhi:get_all_loaivukhi",
		"loaivukhi:get_all_ten_loaivukhi",
		"loaivukhi:get_all_mo_ta_loaivukhi")

	if _, err := pipe.Exec(ctx); err != nil {
		log.Println("REDIS INVALIDATE ERROR:", err)
	}
}
