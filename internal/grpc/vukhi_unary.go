package grpcsrv

import (
	"context"
	"strings"

	grpc_service "game/internal/service/grpc_service"
	"game/internal/utils"
	v1 "game/v1/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

type VuKhiGRPCServer struct {
	v1.UnimplementedVuKhiServiceServer
	grpc_service *grpc_service.VuKhiService
}

func NewVuKhiGRPCServer(svc *grpc_service.VuKhiService) *VuKhiGRPCServer {
	return &VuKhiGRPCServer{grpc_service: svc}
}

func (s *VuKhiGRPCServer) CreateVuKhi(ctx context.Context, req *v1.TaoVuKhiRequest) (*v1.VuKhi, error) {
	if req.TenVuKhi == "" {
		return nil, utils.NewError("ten_vu_khi là bắt buộc", utils.ErrCodeBadRequest)
	}
	if req.SatThuongCoBan == 0 {
		return nil, utils.NewError("sat_thuong_co_ban là bắt buộc", utils.ErrCodeBadRequest)
	}
	if req.TocDoDanh == 0 {
		return nil, utils.NewError("toc_do_danh là bắt buộc", utils.ErrCodeBadRequest)
	}
	if req.TamDanh == 0 {
		return nil, utils.NewError("tam_danh là bắt buộc", utils.ErrCodeBadRequest)
	}
	if req.MoTa == "" {
		req.MoTa = "Mô tả"
	}
	if req.MaHe == 0 {
		return nil, utils.NewError("ma_he là bắt buộc", utils.ErrCodeBadRequest)
	}
	if req.MaDoHiem == 0 {
		return nil, utils.NewError("ma_do_hiem là bắt buộc", utils.ErrCodeBadRequest)
	}
	if req.MaLoaiVuKhi == 0 {
		return nil, utils.NewError("ma_loai_vu_khi là bắt buộc", utils.ErrCodeBadRequest)
	}

	return s.grpc_service.CreateVuKhiService(ctx, req)
}

func (s *VuKhiGRPCServer) UpdateByNameVuKhi(ctx context.Context, req *v1.CapNhatTheoTenVuKhiRequest) (*v1.VuKhi, error) {
	if req.TenVuKhiNew == "" {
		return nil, utils.NewError("ten_vu_khi_new là bắt buộc", utils.ErrCodeBadRequest)
	}
	if strings.TrimSpace(req.TenVuKhi) == "" {
		return nil, utils.NewError("ten_vu_khi (path) là bắt buộc", utils.ErrCodeBadRequest)
	}
	if req.Version <= 0 {
		return nil, utils.NewError("version phải > 0", utils.ErrCodeBadRequest)
	}

	return s.grpc_service.UpdateByNameVuKhiService(ctx, req)
}

func (s *VuKhiGRPCServer) DeleteByNameVuKhi(ctx context.Context, req *v1.XoaTheoTenVuKhiRequest) (*v1.XoaTheoTenVuKhiResponse, error) {
	if req.TenVuKhi == "" {
		return nil, utils.NewError("ten_vu_khi là bắt buộc", utils.ErrCodeBadRequest)
	}
	return s.grpc_service.DeleteByNameVuKhiService(ctx, req)
}

func (s *VuKhiGRPCServer) GetAllVuKhi(ctx context.Context, req *v1.LayTatCaRequest) (*v1.DanhSachVuKhi, error) {
	return s.grpc_service.GetAllVuKhiService(ctx)
}

func (s *VuKhiGRPCServer) GetAllTenVuKhi(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachStringVuKhi, error) {
	return s.grpc_service.GetAllTenVuKhiService(ctx)
}

func (s *VuKhiGRPCServer) GetAllSatThuongCoBan(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachIntVuKhi, error) {
	return s.grpc_service.GetAllSatThuongCoBanVuKhiService(ctx, nil)
}

func (s *VuKhiGRPCServer) GetAllTamDanh(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachIntVuKhi, error) {
	return s.grpc_service.GetAllTamDanhVuKhiService(ctx, nil)
}

func (s *VuKhiGRPCServer) GetAllTocDoDanh(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachDoubleVuKhi, error) {
	return s.grpc_service.GetAllTocDoDanhVuKhiService(ctx, nil)
}

func (s *VuKhiGRPCServer) SearchVuKhi(ctx context.Context, req *v1.TimKiemVuKhiRequest) (*v1.DanhSachVuKhi, error) {
	items, total, err := s.grpc_service.SearchVuKhiService(ctx, req)
	if err != nil {
		return nil, err
	}

	return &v1.DanhSachVuKhi{
		Items:     items,
		Total:     int32(total),
		ItemCount: int32(len(items)),
	}, nil
}

func (s *VuKhiGRPCServer) LayVersion(ctx context.Context, req *v1.LayVersionRequest) (*v1.LayVersionResponse, error) {
	if req.TenVuKhi == "" {
		return nil, utils.NewError("ten_vu_khi là bắt buộc", utils.ErrCodeBadRequest)
	}
	return s.grpc_service.GetVersionVuKhiService(ctx, req)
}
