package grpc_repository

import (
	"context"
	"strings"

	"game/ent"
	"game/ent/he"
	v1 "game/v1/proto"

	entsql "entgo.io/ent/dialect/sql"
)

type HeRepo struct {
	ent *ent.Client
}

func NewHeRepo(ent *ent.Client) *HeRepo {
	return &HeRepo{ent: ent}
}

// Create trong repository
func (r *HeRepo) CreateHeRepo(ctx context.Context, tx *ent.Tx, req *v1.TaoHeRequest, name string) (*ent.He, error) {
	return tx.He.Create().
		SetTenHe(name).
		SetMoTa(req.MoTa).
		Save(ctx)
}

func (r *HeRepo) GetByIDHeRepo(ctx context.Context, id int) (*ent.He, error) {
	return r.ent.He.
		Query().
		Where(he.IDEQ(id)).
		Only(ctx)
}

func (r *HeRepo) LockByNameHeRepo(ctx context.Context, tx *ent.Tx, name string) (*ent.He, error) {
	return tx.He.Query().
		Where(he.TenHeEQ(name)).
		Select().
		Modify(func(s *entsql.Selector) {
			s.ForUpdate()
		}).
		Only(ctx)
}

func (r *HeRepo) UpdateHeRepo(ctx context.Context, tx *ent.Tx, id int, req *v1.CapNhatTheoTenHeRequest) (*ent.He, error) {
	return tx.He.UpdateOneID(id).
		SetTenHe(req.TenHe).
		SetMoTa(req.MoTa).
		Save(ctx)
}

func (r *HeRepo) DeleteByIDHeRepo(ctx context.Context, tx *ent.Tx, id int) error {
	return tx.He.
		DeleteOneID(id).
		Exec(ctx)
}

func (r *HeRepo) GetAllHeRepo(ctx context.Context) ([]*ent.He, error) {
	return r.ent.He.
		Query().
		All(ctx)
}

func (r *HeRepo) GetAllTenHeRepo(ctx context.Context) ([]string, error) {
	return r.ent.He.
		Query().
		Select(he.FieldTenHe).
		Strings(ctx)
}

func (r *HeRepo) GetAllMoTaHeRepo(ctx context.Context) ([]string, error) {
	return r.ent.He.
		Query().
		Select(he.FieldMoTa).
		Strings(ctx)
}

func (r *HeRepo) SearchHeRepo(ctx context.Context, req *v1.TimKiemHeRequest, limit int, offset int) ([]*ent.He, int64, error) {
	q := r.ent.He.Query()

	if req.MaHe != nil {
		q = q.Where(he.IDEQ(int(req.MaHe.Value)))
	} else {
		if req.MinMaHe != nil {
			q = q.Where(he.IDGTE(int(req.MinMaHe.Value)))
		}
		if req.MaxMaHe != nil {
			q = q.Where(he.IDLTE(int(req.MaxMaHe.Value)))
		}
	}
	if t := strings.TrimSpace(req.TenHe); t != "" {
		q = q.Where(he.TenHeContainsFold(t))
	}
	if t := strings.TrimSpace(req.MoTa); t != "" {
		q = q.Where(he.MoTaContainsFold(t))
	}
	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// sắp xếp tăng giảm theo từng trường
	if req.ArrangeAsc != "" {
		switch req.ArrangeAsc {
		case "id":
			q = q.Order(ent.Asc(he.FieldID))
		case "ten_he":
			q = q.Order(ent.Asc(he.FieldTenHe))
		case "mo_ta":
			q = q.Order(ent.Asc(he.FieldMoTa))
		}
	} else if req.ArrangeDesc != "" {
		switch req.ArrangeDesc {
		case "id":
			q = q.Order(ent.Desc(he.FieldID))
		case "ten_he":
			q = q.Order(ent.Desc(he.FieldTenHe))
		case "mo_ta":
			q = q.Order(ent.Desc(he.FieldMoTa))
		}
	} else {
		q = q.Order(he.ByID())
	}

	items, err := q.
		Order(he.ByID()).
		Limit(limit).
		Offset(offset).
		All(ctx)

	if err != nil {
		return nil, 0, err
	}

	return items, int64(total), nil
}
