package grpc_service

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	v1 "game/v1/proto"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func normalizeLimit(l int32) int {
	n := int(l)
	if n <= 0 || n > 1000 {
		n = 30
	}
	return n
}

// parseCursor giải mã chuỗi cursor từ client gửi lên thành số nguyên (offset) để truy vấn DB.
// Hỗ trợ xử lý các trường hợp chuỗi rỗng hoặc template lỗi trước khi chuyển đổi.
func parseCursor(cur string) (int, error) {
	cur = strings.TrimSpace(cur)

	if cur == "" {
		return 0, nil
	}

	if strings.Contains(cur, "{{") {
		return 0, nil
	}

	v, err := strconv.Atoi(cur)
	if err != nil {
		return 0, fmt.Errorf("cursor không phải số")
	}

	if v < 0 {
		return 0, fmt.Errorf("cursor không được âm")
	}

	return v, nil
}

// invalidateVuKhiCache thực hiện xóa hoặc làm mới các key cache liên quan đến vũ Khí trong Redis.
// Hàm này thường được gọi sau các thao tác ghi (Create, Update, Delete) để đảm bảo tính nhất quán dữ liệu.
// Sử dụng Pipeline để tối ưu hóa việc gửi nhiều lệnh đến Redis cùng lúc.
func (s *VuKhiService) invalidateVuKhiCache(ctx context.Context) {
	pipe := s.rdb.Pipeline()

	pipe.Incr(ctx, "vukhi:search:version")
	pipe.Del(ctx,
		"vukhi:get_all_vukhi",
		"vukhi:get_all_ten_vukhi",
		"vukhi:get_all_sat_thuong_co_ban_vukhi",
		"vukhi:get_all_tam_danh_vukhi",
		"vukhi:get_all_toc_do_danh_vukhi",
	)

	if _, err := pipe.Exec(ctx); err != nil {
		log.Println("REDIS INVALIDATE ERROR:", err)
	}
}

// getStreamCacheKey tạo ra một mã khóa (key) định danh cho cache của luồng stream dữ liệu.
// Key này kết hợp giữa Global Version (để auto-expire) và mã băm SHA1 của chuỗi tìm kiếm q.
func (s *VuKhiService) getStreamCacheKey(ctx context.Context, q string) string {
	const kVuKhiVersionKey = "vukhi:version:global"
	ver, err := s.rdb.Get(ctx, kVuKhiVersionKey).Int()
	if err != nil {
		ver = 1
	}

	h := sha1.New()
	h.Write([]byte(q))
	hashStr := hex.EncodeToString(h.Sum(nil))[:10]

	return fmt.Sprintf("vukhi:stream:v%d:%s", ver, hashStr)
}

// sendInChunks chia nhỏ danh sách dữ liệu lớn từ cache thành từng gói nhỏ (mặc định 100 items).
// Việc này giúp tránh quá tải băng thông truyền tải và phù hợp với cơ chế gRPC Server Streaming.
func (s *VuKhiService) sendInChunks(cached *v1.DanhSachVuKhi, send func(*v1.DanhSachVuKhiStream) error) error {
	const chunkSize = 100
	for i := 0; i < len(cached.Items); i += chunkSize {
		end := i + chunkSize
		if end > len(cached.Items) {
			end = len(cached.Items)
		}
		if err := send(&v1.DanhSachVuKhiStream{
			Total: cached.Total,
			Items: cached.Items[i:end],
		}); err != nil {
			return err
		}
	}
	return nil
}

// saveCache thực hiện lưu trữ kết quả truy vấn vào Redis dưới dạng Protobuf.
// Dữ liệu được lưu kèm theo tổng số lượng và có thời gian sống (TTL) là 15 phút.
func (s *VuKhiService) saveCache(key string, items []*v1.VuKhi) {
	if len(items) == 0 {
		return
	}
	ctx := context.Background()
	resp := &v1.DanhSachVuKhi{
		Items: items,
		Total: int32(len(items)),
	}
	_ = s.cache.SetProto(ctx, key, resp, 15*time.Minute)
}

// startBatchWorker vận hành một worker chạy ngầm để gom nhóm các yêu cầu cập nhật.
//
// - Gom đủ BatchSize (100) yêu cầu thì xử lý một lần để tối ưu hóa Database.
// - Hoặc nếu sau BatchTimeout (20ms) mà chưa đủ số lượng, vẫn xử lý để đảm bảo độ trễ thấp.
func (s *VuKhiService) startBatchWorker() {
	const (
		BatchSize    = 100                   // Gom 100 request
		BatchTimeout = 20 * time.Millisecond // Hoặc chờ tối đa 20ms
	)

	jobs := make([]*updateJob, 0, BatchSize)
	ticker := time.NewTicker(BatchTimeout)
	defer ticker.Stop()

	for {
		select {
		case job := <-s.batchQueue:
			jobs = append(jobs, job)
			if len(jobs) >= BatchSize {
				s.processBatch(jobs)
				jobs = make([]*updateJob, 0, BatchSize)
			}

		case <-ticker.C:
			if len(jobs) > 0 {
				s.processBatch(jobs)
				jobs = make([]*updateJob, 0, BatchSize)
			}
		}
	}
}

// // processBatch thực thi việc cập nhật dữ liệu vũ khí theo lô (Batch Processing).
// //
// // Hàm này được các worker gọi để xử lý danh sách các job lấy từ batchQueue.
// func (s *VuKhiService) processBatch(jobs []*updateJob) {
// 	// 1. Sắp xếp để chống Deadlock
// 	// Nhờ có lệnh sort.Slice, cả hai worker đều sẽ cố gắng khóa món "Thanh Long Đao" đầu tiên.
// 	//
// 	// Database sẽ bắt Worker 2 phải xếp hàng đợi Worker 1 xong xuôi (Commit xong) mới cho phép Worker 2 bắt đầu lô hàng của nó.
// 	sort.Slice(jobs, func(i, j int) bool {
// 		return jobs[i].req.TenVuKhi < jobs[j].req.TenVuKhi
// 	})

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	// 2. Gom nhóm Transaction: Mở một Transaction duy nhất cho toàn bộ lô hàng
// 	tx, err := s.ent.Tx(ctx)
// 	if err != nil {
// 		for _, job := range jobs {
// 			job.resChan <- jobResult{
// 				Index: job.Index,
// 				Name:  job.req.TenVuKhi,
// 				Err:   status.Errorf(codes.Internal, "DB Bận: %v", err),
// 			}
// 		}
// 		return
// 	}

// 	results := make([]error, len(jobs))
// 	for i, job := range jobs {
// 		results[i] = s.grpc_repository.UpdateInTx(ctx, tx, job.req)
// 	}

// 	if err := tx.Commit(); err != nil {
// 		log.Printf("Batch Commit Thất bại: %v", err)
// 		for _, job := range jobs {
// 			job.resChan <- jobResult{
// 				Index: job.Index,
// 				Name:  job.req.TenVuKhi,
// 				Err:   status.Errorf(codes.Aborted, "Transaction Commit Thất bại"),
// 			}
// 		}
// 		return
// 	}

// 	for i, job := range jobs {
// 		job.resChan <- jobResult{
// 			Index: job.Index,
// 			Name:  job.req.TenVuKhi,
// 			Err:   results[i],
// 		}
// 	}

// 	go s.invalidateVuKhiCache(context.Background())
// }

// processBatch xử lý danh sách job từ queue
func (s *VuKhiService) processBatch(jobs []*updateJob) {
	// 1. Sắp xếp để chống Deadlock giữa các Worker
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].req.TenVuKhi < jobs[j].req.TenVuKhi
	})

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	for _, job := range jobs {
		// 2. Thực thi Transaction riêng cho từng món đồ (Ai làm nấy chịu)
		err := s.executeSingleUpdate(ctx, job)
		
		// 3. Phản hồi kết quả về cho luồng gRPC tương ứng
		job.resChan <- jobResult{
			Index: job.Index,
			TenVuKhiNew:  job.req.TenVuKhiNew,
			Err:   err,
		}
	}

	// Xóa cache một lần sau khi xong cả lô
	go s.invalidateVuKhiCache(context.Background())
}

func (s *VuKhiService) executeSingleUpdate(ctx context.Context, job *updateJob) error {
	tx, err := s.ent.Tx(ctx)
	if err != nil {
		return status.Errorf(codes.ResourceExhausted, "DB bận (TX Error)")
	}
	defer tx.Rollback()

	// Gọi Repository xử lý Optimistic Locking
	if err := s.grpc_repository.UpdateInTx(ctx, tx, job.req); err != nil {
		return err // Trả về lỗi Aborted (Version mismatch) hoặc Internal
	}

	return tx.Commit()
}