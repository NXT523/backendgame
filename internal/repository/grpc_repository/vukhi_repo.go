package grpc_repository

import (
	"context"
	"strings"

	"game/ent"
	"game/ent/dohiem"
	"game/ent/vukhi"
	v1 "game/v1/proto"

	entsql "entgo.io/ent/dialect/sql"
)

type VuKhiRepo struct {
	ent *ent.Client
}

func NewVuKhiRepo(ent *ent.Client) *VuKhiRepo {
	return &VuKhiRepo{ent: ent}
}

// Create trong repository
func (r *VuKhiRepo) CreateVuKhiRepo(ctx context.Context, tx *ent.Tx, req *v1.TaoVuKhiRequest, name string) (*ent.VuKhi, error) {
	return tx.VuKhi.Create().
		SetTenVuKhi(name).
		SetSatThuongCoBan(int(req.SatThuongCoBan)).
		SetTocDoDanh(req.TocDoDanh).
		SetTamDanh(int(req.TamDanh)).
		SetMoTa(req.MoTa).
		SetMaLoaiVuKhi(int(req.MaLoaiVuKhi)).
		SetMaDoHiem(int(req.MaDoHiem)).
		SetMaHe(int(req.MaHe)).
		SetVersion(1).
		Save(ctx)
}

func (r *VuKhiRepo) GetByIDVuKhiRepo(ctx context.Context, id int) (*ent.VuKhi, error) {
	return r.ent.VuKhi.
		Query().
		Where(vukhi.IDEQ(id)).
		WithLoai().
		WithDoHiem().
		WithHe().
		Only(ctx)
}

func (r *VuKhiRepo) LockByNameVuKhiRepo(ctx context.Context, tx *ent.Tx, name string) (*ent.VuKhi, error) {
	return tx.VuKhi.Query().
		Where(vukhi.TenVuKhiEQ(name)).
		Select().
		Modify(func(s *entsql.Selector) {
			s.ForUpdate()
		}).
		Only(ctx)
}

func (r *VuKhiRepo) UpdateVuKhiRepo(ctx context.Context, tx *ent.Tx, id int, req *v1.CapNhatTheoTenVuKhiRequest, version int) (*ent.VuKhi, error) {
	return tx.VuKhi.UpdateOneID(id).
		SetTenVuKhi(req.TenVuKhiNew).
		SetSatThuongCoBan(int(req.SatThuongCoBan)).
		SetTocDoDanh(req.TocDoDanh).
		SetTamDanh(int(req.TamDanh)).
		SetMoTa(req.MoTa).
		SetMaLoaiVuKhi(int(req.MaLoaiVuKhi)).
		SetMaDoHiem(int(req.MaDoHiem)).
		SetMaHe(int(req.MaHe)).
		SetVersion(version).
		Save(ctx)
}

func (r *VuKhiRepo) DeleteByIDVuKhiRepo(ctx context.Context, tx *ent.Tx, id int) error {
	return tx.VuKhi.
		DeleteOneID(id).
		Exec(ctx)
}

func (r *VuKhiRepo) GetAllVuKhiRepo(ctx context.Context) ([]*ent.VuKhi, error) {
	return r.ent.VuKhi.
		Query().
		WithLoai().
		WithDoHiem().
		WithHe().
		All(ctx)
}

func (r *VuKhiRepo) GetAllTenVuKhiRepo(ctx context.Context) ([]string, error) {
	return r.ent.VuKhi.
		Query().
		Select(vukhi.FieldTenVuKhi).
		Strings(ctx)
}

func (r *VuKhiRepo) GetAllSatThuongCoBanVuKhiRepo(ctx context.Context) ([]int32, error) {
	values, err := r.ent.VuKhi.
		Query().
		Select(vukhi.FieldSatThuongCoBan).
		Ints(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]int32, 0, len(values))
	for _, v := range values {
		out = append(out, int32(v))
	}
	return out, nil
}

func (r *VuKhiRepo) GetAllTocDoDanhVuKhiRepo(ctx context.Context) ([]float64, error) {
	return r.ent.VuKhi.
		Query().
		Select(vukhi.FieldTocDoDanh).
		Float64s(ctx)
}

func (r *VuKhiRepo) GetAllTamDanhVuKhiRepo(ctx context.Context) ([]int32, error) {
	values, err := r.ent.VuKhi.
		Query().
		Select(vukhi.FieldTamDanh).
		Ints(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]int32, 0, len(values))
	for _, v := range values {
		out = append(out, int32(v))
	}
	return out, nil
}

func (r *VuKhiRepo) SearchVuKhiRepo(ctx context.Context, req *v1.TimKiemVuKhiRequest, limit int, offset int) ([]*ent.VuKhi, int64, error) {
	q := r.ent.VuKhi.Query()

	//  MaVuKhi
	if req.MaVuKhi != nil {
		q = q.Where(vukhi.IDEQ(int(req.MaVuKhi.Value)))
	} else {
		if req.MinMaVuKhi != nil {
			q = q.Where(vukhi.IDGTE(int(req.MinMaVuKhi.Value)))
		}
		if req.MaxMaVuKhi != nil {
			q = q.Where(vukhi.IDLTE(int(req.MaxMaVuKhi.Value)))
		}
	}

	//  TenVuKhi
	if t := strings.TrimSpace(req.TenVuKhi); t != "" {
		q = q.Where(vukhi.TenVuKhiContainsFold(t))
	}

	//  FK
	if req.MaLoaiVuKhi != nil {
		q = q.Where(vukhi.MaLoaiVuKhiEQ(int(req.MaLoaiVuKhi.Value)))
	}
	if req.MaHe != nil {
		q = q.Where(vukhi.MaHeEQ(int(req.MaHe.Value)))
	}
	if req.MaDoHiem != nil {
		q = q.Where(vukhi.MaDoHiemEQ(int(req.MaDoHiem.Value)))
	}

	//  SatThuongCoBan
	if req.SatThuongCoBan != nil {
		q = q.Where(vukhi.SatThuongCoBanEQ(int(req.SatThuongCoBan.Value)))
	} else {
		if req.MinSatThuongCoBan != nil {
			q = q.Where(vukhi.SatThuongCoBanGTE(int(req.MinSatThuongCoBan.Value)))
		}
		if req.MaxSatThuongCoBan != nil {
			q = q.Where(vukhi.SatThuongCoBanLTE(int(req.MaxSatThuongCoBan.Value)))
		}
	}

	//  TocDoDanh
	if req.TocDoDanh != nil {
		q = q.Where(vukhi.TocDoDanhEQ(req.TocDoDanh.Value))
	} else {
		if req.MinTocDoDanh != nil {
			q = q.Where(vukhi.TocDoDanhGTE(req.MinTocDoDanh.Value))
		}
		if req.MaxTocDoDanh != nil {
			q = q.Where(vukhi.TocDoDanhLTE(req.MaxTocDoDanh.Value))
		}
	}

	//  TamDanh
	if req.TamDanh != nil {
		q = q.Where(vukhi.TamDanhEQ(int(req.TamDanh.Value)))
	} else {
		if req.MinTamDanh != nil {
			q = q.Where(vukhi.TamDanhGTE(int(req.MinTamDanh.Value)))
		}
		if req.MaxTamDanh != nil {
			q = q.Where(vukhi.TamDanhLTE(int(req.MaxTamDanh.Value)))
		}
	}

	//  MauSac
	if t := strings.TrimSpace(req.MauSac); t != "" {
		q = q.Where(vukhi.HasDoHiemWith(
			dohiem.MauSacContainsFold(t),
		))
	} else {
		if req.MinMauSac != "" {
			q = q.Where(vukhi.HasDoHiemWith(
				dohiem.MauSacGTE(req.MinMauSac),
			))
		}
		if req.MaxMauSac != "" {
			q = q.Where(vukhi.HasDoHiemWith(
				dohiem.MauSacLTE(req.MaxMauSac),
			))
		}
	}

	//  SoLuong
	if req.SoLuong != nil {
		q = q.Where(vukhi.HasDoHiemWith(
			dohiem.SoLuongEQ(int(req.SoLuong.Value)),
		))
	} else {
		if req.MinSoLuong != nil {
			q = q.Where(vukhi.HasDoHiemWith(
				dohiem.SoLuongGTE(int(req.MinSoLuong.Value)),
			))
		}
		if req.MaxSoLuong != nil {
			q = q.Where(vukhi.HasDoHiemWith(
				dohiem.SoLuongLTE(int(req.MaxSoLuong.Value)),
			))
		}
	}

	//  SatThuongBonus
	if req.SatThuongBonus != nil {
		q = q.Where(vukhi.HasDoHiemWith(
			dohiem.SatThuongBonusEQ(req.SatThuongBonus.Value),
		))
	} else {
		if req.MinSatThuongBonus != nil {
			q = q.Where(vukhi.HasDoHiemWith(
				dohiem.SatThuongBonusGTE(req.MinSatThuongBonus.Value),
			))
		}
		if req.MaxSatThuongBonus != nil {
			q = q.Where(vukhi.HasDoHiemWith(
				dohiem.SatThuongBonusLTE(req.MaxSatThuongBonus.Value),
			))
		}
	}

	//  TocDoDanhBonus
	if req.TocDoDanhBonus != nil {
		q = q.Where(vukhi.HasDoHiemWith(
			dohiem.TocDoDanhBonusEQ(req.TocDoDanhBonus.Value),
		))
	} else {
		if req.MinTocDoDanhBonus != nil {
			q = q.Where(vukhi.HasDoHiemWith(
				dohiem.TocDoDanhBonusGTE(req.MinTocDoDanhBonus.Value),
			))
		}
		if req.MaxTocDoDanhBonus != nil {
			q = q.Where(vukhi.HasDoHiemWith(
				dohiem.TocDoDanhBonusLTE(req.MaxTocDoDanhBonus.Value),
			))
		}
	}

	//  CapBac
	if req.CapBac != nil {
		q = q.Where(vukhi.HasDoHiemWith(
			dohiem.CapBacEQ(int(req.CapBac.Value)),
		))
	} else {
		if req.MinCapBac != nil {
			q = q.Where(vukhi.HasDoHiemWith(
				dohiem.CapBacGTE(int(req.MinCapBac.Value)),
			))
		}
		if req.MaxCapBac != nil {
			q = q.Where(vukhi.HasDoHiemWith(
				dohiem.CapBacLTE(int(req.MaxCapBac.Value)),
			))
		}
	}

	//  Sắp xếp tăng giảm
	if req.ArrangeAsc != "" {
		switch req.ArrangeAsc {
		case "id", "ma_vu_khi":
			q = q.Order(ent.Asc(vukhi.FieldID))

		case "ten_vu_khi":
			q = q.Order(ent.Asc(vukhi.FieldTenVuKhi))

		case "sat_thuong_co_ban":
			q = q.Order(ent.Asc(vukhi.FieldSatThuongCoBan))

		case "toc_do_danh":
			q = q.Order(ent.Asc(vukhi.FieldTocDoDanh))

		case "tam_danh":
			q = q.Order(ent.Asc(vukhi.FieldTamDanh))

		case "ma_loai_vu_khi":
			q = q.Order(ent.Asc(vukhi.FieldMaLoaiVuKhi))

		case "ma_he":
			q = q.Order(ent.Asc(vukhi.FieldMaHe))

		case "ma_do_hiem":
			q = q.Order(ent.Asc(vukhi.FieldMaDoHiem))

		case "version":
			q = q.Order(ent.Asc(vukhi.FieldVersion))

		//  (JOIN)
		case "ten_do_hiem":
			q = q.Order(ent.Asc(dohiem.FieldTenDoHiem))

		case "so_luong":
			q = q.Order(ent.Asc(dohiem.FieldSoLuong))

		case "sat_thuong_bonus":
			q = q.Order(ent.Asc(dohiem.FieldSatThuongBonus))

		case "toc_do_danh_bonus":
			q = q.Order(ent.Asc(dohiem.FieldTocDoDanhBonus))

		case "mau_sac":
			q = q.Order(ent.Asc(dohiem.FieldMauSac))

		case "cap_bac":
			q = q.Order(ent.Asc(dohiem.FieldCapBac))
		}
	} else if req.ArrangeDesc != "" {
		switch req.ArrangeDesc {
		case "id", "ma_vu_khi":
			q = q.Order(ent.Desc(vukhi.FieldID))

		case "ten_vu_khi":
			q = q.Order(ent.Desc(vukhi.FieldTenVuKhi))

		case "sat_thuong_co_ban":
			q = q.Order(ent.Desc(vukhi.FieldSatThuongCoBan))

		case "toc_do_danh":
			q = q.Order(ent.Desc(vukhi.FieldTocDoDanh))

		case "tam_danh":
			q = q.Order(ent.Desc(vukhi.FieldTamDanh))

		case "ma_loai_vu_khi":
			q = q.Order(ent.Desc(vukhi.FieldMaLoaiVuKhi))

		case "ma_he":
			q = q.Order(ent.Desc(vukhi.FieldMaHe))

		case "ma_do_hiem":
			q = q.Order(ent.Desc(vukhi.FieldMaDoHiem))

		case "version":
			q = q.Order(ent.Desc(vukhi.FieldVersion))

		//  DoHiem (JOIN)
		case "ten_do_hiem":
			q = q.Order(ent.Desc(dohiem.FieldTenDoHiem))

		case "so_luong":
			q = q.Order(ent.Desc(dohiem.FieldSoLuong))

		case "sat_thuong_bonus":
			q = q.Order(ent.Desc(dohiem.FieldSatThuongBonus))

		case "toc_do_danh_bonus":
			q = q.Order(ent.Desc(dohiem.FieldTocDoDanhBonus))

		case "mau_sac":
			q = q.Order(ent.Desc(dohiem.FieldMauSac))

		case "cap_bac":
			q = q.Order(ent.Desc(dohiem.FieldCapBac))
		}
	} else {
		q = q.Order(vukhi.ByID())
	}

	//  Count
	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	//  Query
	items, err := q.
		WithLoai().
		WithHe().
		WithDoHiem().
		Limit(limit).
		Offset(offset).
		All(ctx)

	if err != nil {
		return nil, 0, err
	}

	return items, int64(total), nil
}

func (r *VuKhiRepo) GetVersionByNameVuKhiRepo(
	ctx context.Context,
	name string,
) (*ent.VuKhi, error) {

	return r.ent.VuKhi.
		Query().
		Where(vukhi.TenVuKhiEQ(name)).
		Only(ctx)
}
