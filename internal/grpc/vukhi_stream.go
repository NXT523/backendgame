package grpcsrv

import (
	v1 "game/v1/proto"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StreamServerSideGetVuKhi
// Server-side streaming: client gửi 1 request, server stream danh sách vũ khí
// 1. Stream Server-Side (Get All)
func (s *VuKhiGRPCServer) StreamServerSideGetVuKhi(req *v1.LayTatCaRequest, stream v1.VuKhiService_StreamServerSideGetVuKhiServer) error {
	return s.grpc_service.StreamServerSideVuKhiService(
		stream.Context(),
		req.Q,
		func(res *v1.DanhSachVuKhiStream) error {
			return stream.Send(res)
		},
	)
}

// StreamClientSideUpdateVuKhi
// Client-side streaming: client gửi nhiều request, server trả 1 kết quả
// 2. Stream Client-Side (Batch Update)
func (s *VuKhiGRPCServer) StreamClientSideUpdateVuKhi(stream v1.VuKhiService_StreamClientSideUpdateVuKhiServer) error {
	var updates []*v1.CapNhatTheoTenVuKhiRequest

	// Nhận toàn bộ request từ stream
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "recv error: %v", err)
		}

		if req.TenVuKhiNew != "" && req.Version > 0 {
			updates = append(updates, req)
		}
	}

	if len(updates) == 0 {
		return stream.SendAndClose(&v1.BatchUpdateResult{Total: 0})
	}

	// Gửi xuống Service xử lý
	result, err := s.grpc_service.StreamClientSideUpdateVuKhiService(stream.Context(), updates)
	if err != nil {
		return err
	}

	return stream.SendAndClose(result)
}

// StreamBidirectionalUpdateVuKhi
// Bidirectional streaming: cập nhật và phản hồi realtime từng item
// 3. Stream Bidirectional (Realtime Update)
func (s *VuKhiGRPCServer) StreamBidirectionalUpdateVuKhi(stream v1.VuKhiService_StreamBidirectionalUpdateVuKhiServer) error {
    for {
        req, err := stream.Recv()
        if err == io.EOF { return nil }
        if err != nil { return err }

        // 1. KIỂM TRA VALIDATION (Tên không được để trống)
        var validationErr string
        if req.TenVuKhi == "" {
            validationErr = "ten_vu_khi (tên hiện tại) không được để trống"
        } else if req.TenVuKhiNew == "" {
            validationErr = "ten_vu_khi_new (tên mới) không được để trống"
        }

        // Nếu có lỗi validation, gửi phản hồi lỗi ngay lập tức và tiếp tục vòng lặp
        if validationErr != "" {
            resp := &v1.UpdateItemResult{
                Success:      false,
                ErrorMessage: &validationErr, 
            }
            if err := stream.Send(resp); err != nil {
                return err
            }
            continue // Chuyển sang nhận request tiếp theo trong stream
        }

        newVer, errUpdate := s.grpc_service.StreamBidirectionalUpdateVuKhiService(stream.Context(), req)

        resp := &v1.UpdateItemResult{
            TenVuKhi: req.TenVuKhi,
            TenVuKhiNew: req.TenVuKhiNew,
            Success:     errUpdate == nil,
        }

        if errUpdate != nil {
            errMsg := errUpdate.Error()
            resp.ErrorMessage = &errMsg
        } else {
            v := int32(newVer)
            resp.NewVersion = &v
        }

        if err := stream.Send(resp); err != nil {
            return err
        }
    }
}