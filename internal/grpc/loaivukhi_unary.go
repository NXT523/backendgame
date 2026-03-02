package grpcsrv

import (
	"context"
	"game/internal/service/grpc_service"
	"game/internal/utils"
	v1 "game/v1/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

type LoaiVuKhiGRPCServer struct {
	v1.UnimplementedLoaiVuKhiServiceServer
	grpc_service *grpc_service.LoaiVuKhiService
}

func NewLoaiVuKhiGRPCServer(svc *grpc_service.LoaiVuKhiService) *LoaiVuKhiGRPCServer {
	return &LoaiVuKhiGRPCServer{grpc_service: svc}
}

func (s *LoaiVuKhiGRPCServer) CreateLoaiVuKhi(ctx context.Context, req *v1.TaoLoaiVuKhiRequest) (*v1.LoaiVuKhi, error) {
	if req.TenLoaiVuKhi == "" {
		return nil, utils.NewError("ten_loai_vu_khi là bắt buộc", utils.ErrCodeBadRequest)
	}
	if req.MoTa == "" {
		req.MoTa = "Mô tả"
	}
	return s.grpc_service.CreateLoaiVuKhiService(ctx, req)
}

func (s *LoaiVuKhiGRPCServer) UpdateByNameLoaiVuKhi(ctx context.Context, req *v1.CapNhatTheoTenLoaiVuKhiRequest) (*v1.LoaiVuKhi, error) {
	if req.TenLoaiVuKhiNew == "" {
		return nil, utils.NewError("ten_loai_vu_khi_new là bắt buộc", utils.ErrCodeBadRequest)
	}
	if req.TenLoaiVuKhi == "" || req.TenLoaiVuKhiNew == "" {
		return nil, utils.NewError("ten_loai_vu_khi không hợp lệ", utils.ErrCodeBadRequest)
	}
	return s.grpc_service.UpdateByNameLoaiVuKhiService(ctx, req)
}

func (s *LoaiVuKhiGRPCServer) DeleteByNameLoaiVuKhi(ctx context.Context, req *v1.XoaTheoTenLoaiVuKhiRequest) (*v1.XoaTheoTenLoaiVuKhiResponse, error) {
	if req.TenLoaiVuKhi == "" {
		return nil, utils.NewError("ten_loai_vu_khi là bắt buộc", utils.ErrCodeBadRequest)
	}
	return s.grpc_service.DeleteByNameLoaiVuKhiService(ctx, req)
}

func (s *LoaiVuKhiGRPCServer) GetAllLoaiVuKhi(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachLoaiVuKhi, error) {
	return s.grpc_service.GetAllLoaiVuKhiService(ctx)
}

func (s *LoaiVuKhiGRPCServer) GetAllTenLoaiVuKhi(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachStringLoaiVuKhi, error) {
	return s.grpc_service.GetAllTenLoaiVuKhiService(ctx)
}

func (s *LoaiVuKhiGRPCServer) GetAllMoTaLoaiVuKhi(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachStringLoaiVuKhi, error) {
	return s.grpc_service.GetAllMoTaLoaiVuKhiService(ctx)
}

func (s *LoaiVuKhiGRPCServer) SearchLoaiVuKhi(ctx context.Context, req *v1.TimKiemLoaiVuKhiRequest) (*v1.DanhSachLoaiVuKhi, error) {
	items, total, err := s.grpc_service.SearchLoaiVuKhiService(ctx, req)
	if err != nil {
		return nil, err
	}

	return &v1.DanhSachLoaiVuKhi{
		Items:     items,
		Total:     int32(total),
		ItemCount: int32(len(items)),
	}, nil
}
