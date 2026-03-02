package grpcsrv

import (
	"context"

	"game/internal/service/grpc_service"
	"game/internal/utils"
	v1 "game/v1/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

type HeGRPCServer struct {
	v1.UnimplementedHeServiceServer
	grpc_service *grpc_service.HeService
}

func NewHeGRPCServer(svc *grpc_service.HeService) *HeGRPCServer {
	return &HeGRPCServer{grpc_service: svc}
}

func (s *HeGRPCServer) CreateHe(ctx context.Context, req *v1.TaoHeRequest) (*v1.He, error) {
	if req.TenHe == "" {
		return nil, utils.NewError("ten_he là bắt buộc", utils.ErrCodeBadRequest)
	}
	if req.MoTa == "" {
		req.MoTa = "Mô tả"
	}
	return s.grpc_service.CreateHeService(ctx, req)
}

func (s *HeGRPCServer) UpdateByNameHe(ctx context.Context, req *v1.CapNhatTheoTenHeRequest) (*v1.He, error) {
	if req.TenHeNew == "" {
		return nil, utils.NewError("ten_he_new là bắt buộc", utils.ErrCodeBadRequest)
	}
	if req.TenHe == "" {
		return nil, utils.NewError("ten_he là bắt buộc", utils.ErrCodeBadRequest)
	}

	return s.grpc_service.UpdateByNameHeService(ctx, req)
}

func (s *HeGRPCServer) DeleteByNameHe(ctx context.Context, req *v1.XoaTheoTenHeRequest) (*v1.XoaTheoTenHeResponse, error) {
	if req.TenHe == "" {
		return nil, utils.NewError("ten_he là bắt buộc", utils.ErrCodeBadRequest)
	}
	return s.grpc_service.DeleteByNameHeService(ctx, req)
}

func (s *HeGRPCServer) GetAllHe(ctx context.Context, req *emptypb.Empty) (*v1.DanhSachHe, error) {
	return s.grpc_service.GetAllHeService(ctx)
}

func (s *HeGRPCServer) GetAllTenHe(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachStringHe, error) {
	return s.grpc_service.GetAllTenHeService(ctx)
}

func (s *HeGRPCServer) GetAllMoTaHe(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachStringHe, error) {
	return s.grpc_service.GetAllMoTaHeService(ctx)
}

func (s *HeGRPCServer) SearchHe(ctx context.Context, req *v1.TimKiemHeRequest) (*v1.DanhSachHe, error) {
	items, total, err := s.grpc_service.SearchHeService(ctx, req)
	if err != nil {
		return nil, err
	}

	return &v1.DanhSachHe{
		Items:     items,
		Total:     int32(total),
		ItemCount: int32(len(items)),
	}, nil
}
