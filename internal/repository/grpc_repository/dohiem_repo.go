package grpc_repository

import (
	"context"
	"game/ent"
	"game/ent/dohiem"
	v1 "game/v1/proto"
	"strings"

	entsql "entgo.io/ent/dialect/sql"
)

type DoHiemRepo struct {
	ent *ent.Client
}

func NewDoHiemRepo(ent *ent.Client) *DoHiemRepo {
	return &DoHiemRepo{ent: ent}
}

// Create trong repository
func (r *DoHiemRepo) CreateDoHiemRepo(ctx context.Context, tx *ent.Tx, req *v1.TaoDoHiemRequest, name string) (*ent.DoHiem, error) {
	return tx.DoHiem.Create().
		SetTenDoHiem(name).
		SetSoLuong(int(req.SoLuong)).
		SetMauSac(req.MauSac).
		SetSatThuongBonus(req.SatThuongBonus).
		SetTocDoDanhBonus(req.TocDoDanhBonus).
		SetCapBac(int(req.CapBac)).
		Save(ctx)
}

func (r *DoHiemRepo) GetByIDDoHiemRepo(ctx context.Context, id int) (*ent.DoHiem, error) {
	return r.ent.DoHiem.
		Query().
		Where(dohiem.IDEQ(id)).
		Only(ctx)
}

func (r *DoHiemRepo) LockByNameDoHiemRepo(ctx context.Context, tx *ent.Tx, name string) (*ent.DoHiem, error) {
	return tx.DoHiem.Query().
		Where(dohiem.TenDoHiemEQ(name)).
		Select().
		Modify(func(s *entsql.Selector) {
			s.ForUpdate()
		}).
		Only(ctx)
}

func (r *DoHiemRepo) UpdateDoHiemRepo(ctx context.Context, tx *ent.Tx, id int, req *v1.CapNhatTheoTenDoHiemRequest) (*ent.DoHiem, error) {
	return tx.DoHiem.UpdateOneID(id).
		SetTenDoHiem(req.TenDoHiem).
		Save(ctx)
}

func (r *DoHiemRepo) DeleteByIDDoHiemRepo(ctx context.Context, tx *ent.Tx, id int) error {
	return tx.DoHiem.
		DeleteOneID(id).
		Exec(ctx)
}

func (r *DoHiemRepo) GetAllDoHiemRepo(ctx context.Context) ([]*ent.DoHiem, error) {
	return r.ent.DoHiem.
		Query().
		All(ctx)
}

func (r *DoHiemRepo) GetAllTenDoHiemRepo(ctx context.Context) ([]string, error) {
	return r.ent.DoHiem.
		Query().
		Select(dohiem.FieldTenDoHiem).
		Strings(ctx)
}

func (r *DoHiemRepo) GetAllSoLuongDoHiemRepo(ctx context.Context) ([]int32, error) {
	values, err := r.ent.DoHiem.
		Query().
		Select(dohiem.FieldSoLuong).
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

func (r *DoHiemRepo) GetAllMauSacDoHiemRepo(ctx context.Context) ([]string, error) {
	return r.ent.DoHiem.
		Query().
		Select(dohiem.FieldMauSac).
		Strings(ctx)
}

func (r *DoHiemRepo) GetAllSatThuongBonusDoHiemRepo(ctx context.Context) ([]float64, error) {
	return r.ent.DoHiem.
		Query().
		Select(dohiem.FieldSatThuongBonus).
		Float64s(ctx)
}

func (r *DoHiemRepo) GetAllTocDoDanhBonusDoHiemRepo(ctx context.Context) ([]float64, error) {
	return r.ent.DoHiem.
		Query().
		Select(dohiem.FieldTocDoDanhBonus).
		Float64s(ctx)
}

func (r *DoHiemRepo) GetAllCapBacDoHiemRepo(ctx context.Context) ([]int32, error) {
	values, err := r.ent.DoHiem.
		Query().
		Select(dohiem.FieldCapBac).
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

func (r *DoHiemRepo) SearchDoHiemRepo(ctx context.Context, req *v1.TimKiemDoHiemRequest, limit int, offset int) ([]*ent.DoHiem, int64, error) {
	q := r.ent.DoHiem.Query()

	if req.MaDoHiem != nil {
		q = q.Where(dohiem.IDEQ(int(req.MaDoHiem.Value)))
	} else {
		if req.MinMaDoHiem != nil {
			q = q.Where(dohiem.IDGTE(int(req.MinMaDoHiem.Value)))
		}
		if req.MaxMaDoHiem != nil {
			q = q.Where(dohiem.IDLTE(int(req.MaxMaDoHiem.Value)))
		}
	}
	if t := strings.TrimSpace(req.TenDoHiem); t != "" {
		q = q.Where(dohiem.TenDoHiemContainsFold(t))
	}
	if t := strings.TrimSpace(req.MauSac); t != "" {
		q = q.Where(dohiem.MauSacContainsFold(t))
	}
	if req.SoLuong != nil {
		q = q.Where(dohiem.SoLuongEQ(int(req.SoLuong.Value)))
	} else {
		if req.MinSoLuong != nil {
			q = q.Where(dohiem.SoLuongGTE(int(req.MinSoLuong.Value)))
		}
		if req.MaxSoLuong != nil {
			q = q.Where(dohiem.SoLuongLTE(int(req.MaxSoLuong.Value)))
		}
	}
	if req.SatThuongBonus != nil {
		q = q.Where(dohiem.SatThuongBonusEQ(float64(req.SatThuongBonus.Value)))
	} else {
		if req.MinSatThuongBonus != nil {
			q = q.Where(dohiem.SatThuongBonusGTE(float64(req.MinSatThuongBonus.Value)))
		}
		if req.MaxSatThuongBonus != nil {
			q = q.Where(dohiem.SatThuongBonusLTE(float64(req.MaxSatThuongBonus.Value)))
		}
	}
	if req.TocDoDanhBonus != nil {
		q = q.Where(dohiem.TocDoDanhBonusEQ(float64(req.TocDoDanhBonus.Value)))
	} else {
		if req.MinTocDoDanhBonus != nil {
			q = q.Where(dohiem.TocDoDanhBonusGTE(float64(req.MinTocDoDanhBonus.Value)))
		}
		if req.MaxTocDoDanhBonus != nil {
			q = q.Where(dohiem.TocDoDanhBonusLTE(float64(req.MaxTocDoDanhBonus.Value)))
		}
	}
	if req.CapBac != nil {
		q = q.Where(dohiem.CapBacEQ(int(req.CapBac.Value)))
	} else {
		if req.MinCapBac != nil {
			q = q.Where(dohiem.CapBacGTE(int(req.MinCapBac.Value)))
		}
		if req.MaxCapBac != nil {
			q = q.Where(dohiem.CapBacLTE(int(req.MaxCapBac.Value)))
		}
	}

	if req.ArrangeAsc != "" {
		switch req.ArrangeAsc {
		case "id":
			q = q.Order(ent.Asc(dohiem.FieldID))
		case "ten_do_hiem":
			q = q.Order(ent.Asc(dohiem.FieldTenDoHiem))
		case "so_luong":
			q = q.Order(ent.Asc(dohiem.FieldSoLuong))
		case "mau_sac":
			q = q.Order(ent.Asc(dohiem.FieldMauSac))
		case "sat_thuong_bonus":
			q = q.Order(ent.Asc(dohiem.FieldSatThuongBonus))
		case "toc_do_danh_bonus":
			q = q.Order(ent.Asc(dohiem.FieldTocDoDanhBonus))
		case "cap_bac":
			q = q.Order(ent.Asc(dohiem.FieldCapBac))
		}
	} else if req.ArrangeDesc != "" {
		switch req.ArrangeDesc {
		case "id":
			q = q.Order(ent.Desc(dohiem.FieldID))
		case "ten_do_hiem":
			q = q.Order(ent.Desc(dohiem.FieldTenDoHiem))
		case "so_luong":
			q = q.Order(ent.Desc(dohiem.FieldSoLuong))
		case "mau_sac":
			q = q.Order(ent.Desc(dohiem.FieldMauSac))
		case "sat_thuong_bonus":
			q = q.Order(ent.Desc(dohiem.FieldSatThuongBonus))
		case "toc_do_danh_bonus":
			q = q.Order(ent.Desc(dohiem.FieldTocDoDanhBonus))
		case "cap_bac":
			q = q.Order(ent.Desc(dohiem.FieldCapBac))
		}
	} else {
		q = q.Order(dohiem.ByID())
	}

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	items, err := q.
		Order(dohiem.ByID()).
		Limit(limit).
		Offset(offset).
		All(ctx)

	if err != nil {
		return nil, 0, err
	}

	return items, int64(total), nil
}
