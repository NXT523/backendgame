package grpc_service

import (
	"context"
	mappers "game/internal/mapping/ent"
	v1 "game/v1/proto"
	"sort"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func (s *VuKhiService) StreamServerSideVuKhiService(ctx context.Context, q string, send func(*v1.DanhSachVuKhiStream) error) error {
	cacheKey := s.getStreamCacheKey(ctx, q)
	if data, err := s.rdb.Get(ctx, cacheKey).Bytes(); err == nil {
		var cached v1.DanhSachVuKhi
		if proto.Unmarshal(data, &cached) == nil {
			return s.sendInChunks(&cached, send)
		}
	}

	const batchSize = 500
	var lastID int
	var allItems []*v1.VuKhi

	for {
		list, err := s.grpc_repository.GetBatchByKeyset(ctx, lastID, batchSize, q)
		if err != nil {
			return err
		}
		if len(list) == 0 {
			break
		}

		// Convert Ent -> Proto
		pbItems := make([]*v1.VuKhi, 0, len(list))
		for _, vk := range list {
			pb := mappers.VuKhiEntToPB(vk)
			pbItems = append(pbItems, pb)
			lastID = vk.ID
		}

		allItems = append(allItems, pbItems...)

		// Stream trả về Client ngay lập tức
		if err := send(&v1.DanhSachVuKhiStream{
			Total: int32(len(allItems)), // Lưu ý: Total này chỉ là total tạm thời đã load
			Items: pbItems,
		}); err != nil {
			return err
		}
	}

	s.saveCache(cacheKey, allItems)
	return nil
}

// 2. Stream Client-Side (Batch Update - Write)
func (s *VuKhiService) StreamClientSideUpdateVuKhiService(ctx context.Context, updates []*v1.CapNhatTheoTenVuKhiRequest) (*v1.BatchUpdateResult, error) {
	if len(updates) == 0 {
		return &v1.BatchUpdateResult{Total: 0}, nil
	}

	res := &v1.BatchUpdateResult{
		Total:   int32(len(updates)),
		Details: make([]*v1.BatchItemDetail, 0),
	}

	resChan := make(chan jobResult, len(updates))

	for i, req := range updates {
		s.batchQueue <- &updateJob{
			Index:   i,
			req:     req,
			resChan: resChan,
		}
	}

	for i := 0; i < len(updates); i++ {
		result := <-resChan
		detail := &v1.BatchItemDetail{
			Index:    int32(result.Index),
			TenVuKhi: result.TenVuKhiNew,
		}
		if result.Err == nil {
			res.Success++
			detail.Success = true
			detail.Message = "OK"
		} else {
			res.Failed++
			detail.Success = false
			detail.Message = result.Err.Error()
		}
		res.Details = append(res.Details, detail)
	}

	sort.Slice(res.Details, func(i, j int) bool {
		return res.Details[i].Index < res.Details[j].Index
	})

	return res, nil
}

func (s *VuKhiService) StreamBidirectionalUpdateVuKhiService(ctx context.Context, req *v1.CapNhatTheoTenVuKhiRequest) (int32, error) {
	resChan := make(chan jobResult, 1)
	job := &updateJob{
		req:     req,
		resChan: resChan,
	}

	select {
	case s.batchQueue <- job:
		result := <-resChan

		if result.Err != nil {
			return 0, result.Err
		}

		return req.Version + 1, nil

	case <-ctx.Done():
		return 0, status.Error(codes.Canceled, "Client ngắt kết nối hoặc timeout")

	case <-time.After(30 * time.Second):
		return 0, status.Error(codes.ResourceExhausted, "Server quá tải, hàng đợi đầy")
	}
}
