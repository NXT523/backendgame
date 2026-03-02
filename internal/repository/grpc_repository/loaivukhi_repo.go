package grpc_repository

import (
	"context"
	"strings"

	"game/ent"
	"game/ent/loaivukhi"
	v1 "game/v1/proto"

	entsql "entgo.io/ent/dialect/sql"
)

type LoaiVuKhiRepo struct {
	ent *ent.Client
}

func NewLoaiVuKhiRepo(ent *ent.Client) *LoaiVuKhiRepo {
	return &LoaiVuKhiRepo{ent: ent}
}

func (r *LoaiVuKhiRepo) CreateLoaiVuKhiRepo(ctx context.Context, tx *ent.Tx, req *v1.TaoLoaiVuKhiRequest, name string) (*ent.LoaiVuKhi, error) {
	return tx.LoaiVuKhi.Create().
		SetTenLoaiVuKhi(name).
		SetMoTa(req.MoTa).
		Save(ctx)
}

func (r *LoaiVuKhiRepo) GetByIDLoaiVuKhiRepo(ctx context.Context, id int) (*ent.LoaiVuKhi, error) {
	return r.ent.LoaiVuKhi.
		Query().
		Where(loaivukhi.IDEQ(id)).
		Only(ctx)
}

func (r *LoaiVuKhiRepo) LockByNameLoaiVuKhiRepo(ctx context.Context, tx *ent.Tx, name string) (*ent.LoaiVuKhi, error) {
	return tx.LoaiVuKhi.Query().
		Where(loaivukhi.TenLoaiVuKhiEQ(name)).
		Select().
		Modify(func(s *entsql.Selector) {
			s.ForUpdate()
		}).
		Only(ctx)
}

func (r *LoaiVuKhiRepo) UpdateLoaiVuKhiRepo(ctx context.Context, tx *ent.Tx, id int, req *v1.CapNhatTheoTenLoaiVuKhiRequest) (*ent.LoaiVuKhi, error) {
	return tx.LoaiVuKhi.
		UpdateOneID(id).
		SetTenLoaiVuKhi(req.TenLoaiVuKhi).
		SetMoTa(req.MoTa).
		Save(ctx)
}

func (r *LoaiVuKhiRepo) DeleteByIDLoaiVuKhiRepo(ctx context.Context, tx *ent.Tx, id int) error {
	return tx.LoaiVuKhi.
		DeleteOneID(id).
		Exec(ctx)
}

func (r *LoaiVuKhiRepo) GetAllLoaiVuKhiRepo(ctx context.Context) ([]*ent.LoaiVuKhi, error) {
	return r.ent.LoaiVuKhi.
		Query().
		All(ctx)
}

func (r *LoaiVuKhiRepo) GetAllTenLoaiVuKhiRepo(ctx context.Context) ([]string, error) {
	return r.ent.LoaiVuKhi.
		Query().
		Select(loaivukhi.FieldTenLoaiVuKhi).
		Strings(ctx)
}

func (r *LoaiVuKhiRepo) GetAllMoTaLoaiVuKhiRepo(ctx context.Context) ([]string, error) {
	return r.ent.LoaiVuKhi.
		Query().
		Select(loaivukhi.FieldMoTa).
		Strings(ctx)
}

func (r *LoaiVuKhiRepo) SearchLoaiVuKhiRepo(ctx context.Context, req *v1.TimKiemLoaiVuKhiRequest, limit int, offset int) ([]*ent.LoaiVuKhi, int64, error) {
	q := r.ent.LoaiVuKhi.Query()

	if req.MaLoaiVuKhi != nil {
		q = q.Where(loaivukhi.IDEQ(int(req.MaLoaiVuKhi.Value)))
	} else {
		if req.MinMaLoaiVuKhi != nil {
			q = q.Where(loaivukhi.IDGTE(int(req.MinMaLoaiVuKhi.Value)))
		}
		if req.MaxMaLoaiVuKhi != nil {
			q = q.Where(loaivukhi.IDLTE(int(req.MaxMaLoaiVuKhi.Value)))
		}
	}

	if t := strings.TrimSpace(req.TenLoaiVuKhi); t != "" {
		q = q.Where(loaivukhi.TenLoaiVuKhiContainsFold(t))
	}

	if t := strings.TrimSpace(req.MoTa); t != "" {
		q = q.Where(loaivukhi.MoTaContainsFold(t))
	}

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	if req.ArrangeAsc != "" {
		switch req.ArrangeAsc {
		case "ma_loai_vu_khi":
			q = q.Order(ent.Asc(loaivukhi.FieldID))
		case "ten_loai_vu_khi":
			q = q.Order(ent.Asc(loaivukhi.FieldTenLoaiVuKhi))
		case "mo_ta":
			q = q.Order(ent.Asc(loaivukhi.FieldMoTa))
		}
	} else if req.ArrangeDesc != "" {
		switch req.ArrangeDesc {
		case "ma_loai_vu_khi":
			q = q.Order(ent.Desc(loaivukhi.FieldID))
		case "ten_loai_vu_khi":
			q = q.Order(ent.Desc(loaivukhi.FieldTenLoaiVuKhi))
		case "mo_ta":
			q = q.Order(ent.Desc(loaivukhi.FieldMoTa))
		}
	} else {
		q = q.Order(loaivukhi.ByID())
	}

	items, err := q.
		Limit(limit).
		Offset(offset).
		All(ctx)

	if err != nil {
		return nil, 0, err
	}

	return items, int64(total), nil
}
