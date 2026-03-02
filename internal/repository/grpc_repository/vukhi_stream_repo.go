package grpc_repository

import (
	"context"
	"game/ent"
	"game/ent/vukhi"
	v1 "game/v1/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetBatchByKeyset lấy danh sách vũ khí theo lô dùng Keyset Pagination (ID > lastID)
func (r *VuKhiRepo) GetBatchByKeyset(ctx context.Context, lastID, limit int, q string) ([]*ent.VuKhi, error) {
	qb := r.ent.VuKhi.Query().
		WithDoHiem().WithHe().WithLoai().
		Where(vukhi.IDGT(lastID)).
		Order(ent.Asc(vukhi.FieldID)).
		Limit(limit)

	if q != "" {
		qb = qb.Where(vukhi.TenVuKhiHasPrefix(q))
	}
	return qb.All(ctx)
}

// UpdateInTx: Thực thi lệnh Update chuẩn xác với Optimistic Lock
func (r *VuKhiRepo) UpdateInTx(ctx context.Context, tx *ent.Tx, req *v1.CapNhatTheoTenVuKhiRequest) error {
	n, err := tx.VuKhi.
		Update().
		Where(
			vukhi.TenVuKhiEQ(req.TenVuKhi),
			vukhi.VersionEQ(int(req.Version)),
		).
		SetTenVuKhi(req.TenVuKhiNew).
		SetSatThuongCoBan(int(req.SatThuongCoBan)).
		SetTocDoDanh(req.TocDoDanh).
		SetTamDanh(int(req.TamDanh)).
		SetMoTa(req.MoTa).
		SetMaLoaiVuKhi(int(req.MaLoaiVuKhi)).
		SetMaDoHiem(int(req.MaDoHiem)).
		SetMaHe(int(req.MaHe)).
		SetVersion(int(req.Version) + 1).
		Save(ctx)

	if err != nil {
		return status.Errorf(codes.Internal, "Lỗi hệ thống DB: %v", err)
	}

	if n == 0 {
		exists, errExist := tx.VuKhi.Query().
			Where(vukhi.TenVuKhiEQ(req.TenVuKhi)).
			Exist(ctx)

		if errExist != nil {
			return status.Errorf(codes.Internal, "Lỗi kiểm tra dữ liệu: %v", errExist)
		}

		if !exists {
			return status.Errorf(codes.NotFound, "Vật phẩm '%s' không tồn tại", req.TenVuKhi)
		}

		return status.Errorf(codes.Aborted, "Version %d không khớp với Version hiện tại", req.Version)
	}

	return nil
}
