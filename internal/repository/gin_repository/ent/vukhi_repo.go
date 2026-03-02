package gin_repository

import (
	"context"
	"errors"
	"game/ent"
	"game/ent/dohiem"
	"game/ent/vukhi"
	mappers "game/internal/mapping/ent"
	"game/internal/models"
	"strings"
)

// 1. ENT
type VuKhiRepositoryEnt struct {
	vukhis []models.VuKhis
	client *ent.Client
}

func NewVuKhiRepositoryEnt(client *ent.Client) VuKhiRepoEnt {
	return &VuKhiRepositoryEnt{
		vukhis: make([]models.VuKhis, 0),
		client: client,
	}
}

// Repository: Tạo vũ khí
func (r *VuKhiRepositoryEnt) CreateEnt(ctx context.Context, vk *models.VuKhis) (*models.VuKhis, error) {
	createVuKhi, err := r.client.VuKhi.Create().
		SetTenVuKhi(vk.TenVuKhi).
		SetSatThuongCoBan(vk.SatThuongCoBan).
		SetMaLoaiVuKhi(vk.MaLoaiVuKhi).
		SetMaDoHiem(vk.MaDoHiem).
		SetMaHe(vk.MaHe).
		SetTocDoDanh(vk.TocDoDanh).
		SetTamDanh(vk.TamDanh).
		SetMoTa(vk.MoTa).
		SetVersion(1).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return mappers.VuKhiToModel(createVuKhi), nil
}

func (r *VuKhiRepositoryEnt) UpdateByNameEnt(ctx context.Context, name string, in models.VuKhiUpdateInput) (*models.VuKhis, error) {

	cur, err := r.client.VuKhi.
		Query().
		Where(vukhi.TenVuKhiEQ(name)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	if in.Version > 0 {
		if int(cur.Version) != in.Version {
			return nil, errors.New("data has been modified")
		}
	}

	up := r.client.VuKhi.UpdateOneID(cur.ID)

	if in.TenVuKhi != "" {
		up.SetTenVuKhi(in.TenVuKhi)
	}
	if in.SatThuongCoBan != 0 {
		up.SetSatThuongCoBan(in.SatThuongCoBan)
	}
	if in.TocDoDanh != 0 {
		up.SetTocDoDanh(in.TocDoDanh)
	}
	if in.TamDanh != 0 {
		up.SetTamDanh(in.TamDanh)
	}
	if in.MoTa != "" {
		up.SetMoTa(in.MoTa)
	}
	if in.MaLoaiVuKhi != 0 {
		up.SetMaLoaiVuKhi(in.MaLoaiVuKhi)
	}
	if in.MaDoHiem != 0 {
		up.SetMaDoHiem(in.MaDoHiem)
	}
	if in.MaHe != 0 {
		up.SetMaHe(in.MaHe)
	}

	res, err := up.
		SetVersion(cur.Version + 1).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return mappers.VuKhiToModel(res), nil
}

func (r *VuKhiRepositoryEnt) DeleteByNameEnt(ctx context.Context, name string) error {
	deleteVuKhi, err := r.client.VuKhi.Query().Where(vukhi.TenVuKhiEQ(name)).Only(ctx)
	if err != nil {
		return err
	}
	return r.client.VuKhi.DeleteOneID(deleteVuKhi.ID).Exec(ctx)
}

func (r *VuKhiRepositoryEnt) GetAllEnt(ctx context.Context, kw string) ([]*models.VuKhis, error) {
	query := r.client.VuKhi.Query()

	if kw != "" {
		query = query.Where(vukhi.TenVuKhiContainsFold(kw))
	}

	entList, err := query.
		WithLoai().
		WithDoHiem().
		WithHe().
		All(ctx)

	if err != nil {
		return nil, err
	}
	return mappers.ListVuKhiToModels(entList), nil
}

func (r *VuKhiRepositoryEnt) GetAllTenVuKhiEnt(ctx context.Context) ([]string, error) {
	return r.client.VuKhi.
		Query().
		Select(vukhi.FieldTenVuKhi).
		Strings(ctx)
}

func (r *VuKhiRepositoryEnt) GetAllSatThuongCoBanEnt(ctx context.Context) ([]int, error) {
	return r.client.VuKhi.
		Query().
		Select(vukhi.FieldSatThuongCoBan).
		Ints(ctx)
}

func (r *VuKhiRepositoryEnt) GetAllTocDoDanhEnt(ctx context.Context) ([]float64, error) {
	return r.client.VuKhi.
		Query().
		Select(vukhi.FieldTocDoDanh).
		Float64s(ctx)
}

func (r *VuKhiRepositoryEnt) GetAllTamDanhEnt(ctx context.Context) ([]int, error) {
	return r.client.VuKhi.
		Query().
		Select(vukhi.FieldTamDanh).
		Ints(ctx)
}

func (r *VuKhiRepositoryEnt) SearchEnt(
	ctx context.Context,
	req models.VuKhiSearchRequest,
	limit int,
	offset int,
	arrangeAsc string,
	arrangeDesc string,
) ([]*models.VuKhiResponse, int64, error) {

	q := r.client.VuKhi.Query()

	//  ID
	if req.MaVuKhi != nil {
		q = q.Where(vukhi.IDEQ(*req.MaVuKhi))
	}

	//  Tên
	if t := strings.TrimSpace(req.TenVuKhi); t != "" {
		q = q.Where(vukhi.TenVuKhiContainsFold(t))
	}

	//  FK
	if req.MaLoaiVuKhi != nil {
		q = q.Where(vukhi.MaLoaiVuKhiEQ(*req.MaLoaiVuKhi))
	}
	if req.MaHe != nil {
		q = q.Where(vukhi.MaHeEQ(*req.MaHe))
	}
	if req.MaDoHiem != nil {
		q = q.Where(vukhi.MaDoHiemEQ(*req.MaDoHiem))
	}

	//  Damage
	if req.SatThuongCoBan != nil {
		q = q.Where(vukhi.SatThuongCoBanEQ(*req.SatThuongCoBan))
	}
	if req.MinDamage != nil {
		q = q.Where(vukhi.SatThuongCoBanGTE(*req.MinDamage))
	}
	if req.MaxDamage != nil {
		q = q.Where(vukhi.SatThuongCoBanLTE(*req.MaxDamage))
	}

	//  TocDoDanh
	if req.TocDoDanh != nil {
		q = q.Where(vukhi.TocDoDanhEQ(*req.TocDoDanh))
	}

	//  TamDanh
	if req.TamDanh != nil {
		q = q.Where(vukhi.TamDanhEQ(*req.TamDanh))
	}

	//  JOIN DoHiem

	if strings.TrimSpace(req.TenDoHiem) != "" {
		q = q.Where(vukhi.HasDoHiemWith(
			dohiem.TenDoHiemContainsFold(strings.TrimSpace(req.TenDoHiem)),
		))
	}

	if strings.TrimSpace(req.MauSac) != "" {
		q = q.Where(vukhi.HasDoHiemWith(
			dohiem.MauSacContainsFold(strings.TrimSpace(req.MauSac)),
		))
	}

	if req.SoLuong != nil {
		q = q.Where(vukhi.HasDoHiemWith(
			dohiem.SoLuongEQ(*req.SoLuong),
		))
	}

	if req.SatThuongBonus != nil {
		q = q.Where(vukhi.HasDoHiemWith(
			dohiem.SatThuongBonusEQ(*req.SatThuongBonus),
		))
	}

	if req.TocDoDanhBonus != nil {
		q = q.Where(vukhi.HasDoHiemWith(
			dohiem.TocDoDanhBonusEQ(*req.TocDoDanhBonus),
		))
	}

	//  SORT

	if arrangeAsc != "" {
		switch arrangeAsc {

		case "ma_vu_khi":
			q = q.Order(ent.Asc(vukhi.FieldID))

		case "ten_vu_khi":
			q = q.Order(ent.Asc(vukhi.FieldTenVuKhi))

		case "sat_thuong_co_ban":
			q = q.Order(ent.Asc(vukhi.FieldSatThuongCoBan))

		case "toc_do_danh":
			q = q.Order(ent.Asc(vukhi.FieldTocDoDanh))

		case "tam_danh":
			q = q.Order(ent.Asc(vukhi.FieldTamDanh))

		case "version":
			q = q.Order(ent.Asc(vukhi.FieldVersion))

		case "ten_do_hiem":
			q = q.Order(ent.Asc(dohiem.FieldTenDoHiem))

		case "so_luong":
			q = q.Order(ent.Asc(dohiem.FieldSoLuong))
		}

	} else if arrangeDesc != "" {

		switch arrangeDesc {

		case "ma_vu_khi":
			q = q.Order(ent.Desc(vukhi.FieldID))

		case "ten_vu_khi":
			q = q.Order(ent.Desc(vukhi.FieldTenVuKhi))

		case "sat_thuong_co_ban":
			q = q.Order(ent.Desc(vukhi.FieldSatThuongCoBan))

		case "toc_do_danh":
			q = q.Order(ent.Desc(vukhi.FieldTocDoDanh))

		case "tam_danh":
			q = q.Order(ent.Desc(vukhi.FieldTamDanh))

		case "version":
			q = q.Order(ent.Desc(vukhi.FieldVersion))

		case "ten_do_hiem":
			q = q.Order(ent.Desc(dohiem.FieldTenDoHiem))

		case "so_luong":
			q = q.Order(ent.Desc(dohiem.FieldSoLuong))
		}

	} else {
		q = q.Order(vukhi.ByID())
	}

	// COUNT
	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// QUERY
	entities, err := q.
		WithLoai().
		WithHe().
		WithDoHiem().
		Limit(limit).
		Offset(offset).
		All(ctx)

	if err != nil {
		return nil, 0, err
	}

	// Mapping
	modelsList := mappers.ListVuKhiToModels(entities)
	responses := mappers.ListVuKhiToResponsesSearchENT(modelsList)

	return responses, int64(total), nil
}
