package grpcsrv

import (
	"context"
	"game/internal/service/grpc_service"
	"game/internal/utils"
	v1 "game/v1/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

type DoHiemGRPCServer struct {
	v1.UnimplementedDoHiemServiceServer
	grpc_service *grpc_service.DoHiemService
}

func NewDoHiemGRPCServer(svc *grpc_service.DoHiemService) *DoHiemGRPCServer {
	return &DoHiemGRPCServer{grpc_service: svc}
}

func (s *DoHiemGRPCServer) CreateDoHiem(ctx context.Context, req *v1.TaoDoHiemRequest) (*v1.DoHiem, error) {
	if req.TenDoHiem == "" {
		return nil, utils.NewError("ten_do_hiem là bắt buộc", utils.ErrCodeBadRequest)
	}
	if req.SoLuong == 0 {
		return nil, utils.NewError("so_luong là bắt buộc", utils.ErrCodeBadRequest)
	}
	if req.MauSac == "" {
		return nil, utils.NewError("mau_sac là bắt buộc", utils.ErrCodeBadRequest)
	}

	return s.grpc_service.CreateDoHiemService(ctx, req)
}

func (s *DoHiemGRPCServer) UpdateByNameDoHiem(ctx context.Context, req *v1.CapNhatTheoTenDoHiemRequest) (*v1.DoHiem, error) {
	if req.TenDoHiemNew == "" {
		return nil, utils.NewError("ten_do_hiem_new là bắt buộc", utils.ErrCodeBadRequest)
	}
	if req.TenDoHiem == "" {
		return nil, utils.NewError("ten_do_hiem là bắt buộc", utils.ErrCodeBadRequest)
	}
	if req.SoLuong == 0 {
		return nil, utils.NewError("so_luong là bắt buộc", utils.ErrCodeBadRequest)
	}
	if req.MauSac == "" {
		return nil, utils.NewError("mau_sac là bắt buộc", utils.ErrCodeBadRequest)
	}

	return s.grpc_service.UpdateByNameDoHiemService(ctx, req)
}

func (s *DoHiemGRPCServer) DeleteByNameDoHiem(ctx context.Context, req *v1.XoaTheoTenDoHiemRequest) (*v1.XoaTheoTenDoHiemResponse, error) {
	if req.TenDoHiem == "" {
		return nil, utils.NewError("ten_do_hiem là bắt buộc", utils.ErrCodeBadRequest)
	}
	return s.grpc_service.DeleteByNameDoHiemService(ctx, req)
}

func (s *DoHiemGRPCServer) GetAllDoHiem(ctx context.Context, req *emptypb.Empty) (*v1.DanhSachDoHiem, error) {
	return s.grpc_service.GetAllDoHiemService(ctx)
}

func (s *DoHiemGRPCServer) GetAllTenDoHiem(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachStringDoHiem, error) {
	return s.grpc_service.GetAllTenDoHiemService(ctx)
}

func (s *DoHiemGRPCServer) GetAllSoLuongDoHiem(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachIntDoHiem, error) {
	return s.grpc_service.GetAllSoLuongDoHiemService(ctx, nil)
}

func (s *DoHiemGRPCServer) GetAllMauSacDoHiem(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachStringDoHiem, error) {
	return s.grpc_service.GetAllMauSacDoHiemService(ctx)
}

func (s *DoHiemGRPCServer) GetAllSatThuongBonusDoHiem(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachDoubleDoHiem, error) {
	return s.grpc_service.GetAllSatThuongBonusDoHiemService(ctx, nil)
}
func (s *DoHiemGRPCServer) GetAllTocDoDanhBonusDoHiem(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachDoubleDoHiem, error) {
	return s.grpc_service.GetAllTocDoDanhBonusDoHiemService(ctx, nil)
}

func (s *DoHiemGRPCServer) GetAllCapBacDoHiem(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachIntDoHiem, error) {
	return s.grpc_service.GetAllCapBacDoHiemService(ctx, nil)
}

func (s *DoHiemGRPCServer) SearchDoHiem(ctx context.Context, req *v1.TimKiemDoHiemRequest) (*v1.DanhSachDoHiem, error) {
	items, total, err := s.grpc_service.SearchDoHiemService(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.DanhSachDoHiem{
		Items:     items,
		Total:     int32(total),
		ItemCount: int32(len(items)),
	}, nil
}
