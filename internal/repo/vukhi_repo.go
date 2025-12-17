package repo

import (
	"context"
	"database/sql"
	"strings"

	"game/ent"
	"game/ent/dohiem"
	"game/ent/he"
	"game/ent/loaivukhi"
	"game/ent/vukhi"
	"game/internal/db"

	"gorm.io/gorm"
)

// 1. ENT

type CreateVuKhiInput struct {
	TenVuKhi       string
	SatThuongCoBan int
	TocDoDanh      float64
	TamDanh        int
	MoTa           *string
	MaLoai         int
	MaDoHiem       int
	MaHe           int
}

type VuKhiSearchRequest struct {
	MaVuKhi        *int     `json:"ma_vu_khi"`
	MaLoai         *int     `json:"ma_loai"`
	MaHe           *int     `json:"ma_he"`
	MaDoHiem       *int     `json:"ma_do_hiem"`
	TenVuKhi       string   `json:"ten_vu_khi"`
	TenLoai        string   `json:"ten_loai"`
	TenHe          string   `json:"ten_he"`
	TenDoHiem      string   `json:"ten_do_hiem"`
	SatThuongCoBan *int     `json:"sat_thuong_co_ban"`
	TocDoDanh      *float64 `json:"toc_do_danh"`
	TamDanh        *int     `json:"tam_danh"`
	MauSac         string   `json:"mau_sac"`
	SatThuongBonus *float64 `json:"sat_thuong_bonus"`
	TocDoDanhBonus *float64 `json:"toc_do_danh_bonus"`
	MinDamage      *int     `json:"min_damage"`
	MaxDamage      *int     `json:"max_damage"`
	SoLuong        *int     `json:"so_luong"`
}

type VuKhiRepo interface {
	CreateEnt(ctx context.Context, in CreateVuKhiInput) (*ent.VuKhi, error)
	UpdateByNameEnt(ctx context.Context, name string, mutate func(*ent.VuKhiUpdateOne) *ent.VuKhiUpdateOne) (*ent.VuKhi, error)
	DeleteByNameEnt(ctx context.Context, name string) error

	GetAllEnt(ctx context.Context, q string) ([]*ent.VuKhi, error)
	GetAllTenVuKhiEnt(ctx context.Context) ([]string, error)
	GetAllSatThuongCoBanEnt(ctx context.Context) ([]int, error)
	GetAllTocDoDanhEnt(ctx context.Context) ([]float64, error)
	GetAllTamDanhEnt(ctx context.Context) ([]int, error)

	SearchEnt(ctx context.Context, in VuKhiSearchRequest) ([]*ent.VuKhi, error)
}

type vuKhiRepo struct{ c *ent.Client }

func CreateKhoVuKhiEnt(c *ent.Client) VuKhiRepo { return &vuKhiRepo{c: c} }

func (r *vuKhiRepo) CreateEnt(ctx context.Context, in CreateVuKhiInput) (*ent.VuKhi, error) {
	return r.c.VuKhi.Create().
		SetTenVuKhi(in.TenVuKhi).
		SetSatThuongCoBan(in.SatThuongCoBan).
		SetTocDoDanh(in.TocDoDanh).
		SetTamDanh(in.TamDanh).
		SetNillableMoTa(in.MoTa).
		SetMaLoai(in.MaLoai).
		SetMaDoHiem(in.MaDoHiem).
		SetMaHe(in.MaHe).
		Save(ctx)
}

func (r *vuKhiRepo) GetAllEnt(ctx context.Context, q string) ([]*ent.VuKhi, error) {
	qb := r.c.VuKhi.Query()
	if q != "" {
		qb = qb.Where(vukhi.TenVuKhiContainsFold(q))
	}
	return qb.WithLoai().WithDoHiem().WithHe().All(ctx)
}

func (r vuKhiRepo) UpdateByNameEnt(ctx context.Context, name string, mutate func(*ent.VuKhiUpdateOne) *ent.VuKhiUpdateOne) (*ent.VuKhi, error) {
	cur, err := r.c.VuKhi.Query().Where(vukhi.TenVuKhiEQ(name)).Only(ctx)
	if err != nil {
		return nil, err
	}
	up := r.c.VuKhi.UpdateOneID(cur.ID)
	return mutate(up).Save(ctx)
}

func (r vuKhiRepo) DeleteByNameEnt(ctx context.Context, name string) error {
	cur, err := r.c.VuKhi.Query().Where(vukhi.TenVuKhiEQ(name)).Only(ctx)
	if err != nil {
		return err
	}
	return r.c.VuKhi.DeleteOneID(cur.ID).Exec(ctx)
}

func (r *vuKhiRepo) GetAllTenVuKhiEnt(ctx context.Context) ([]string, error) {
	return r.c.VuKhi.
		Query().
		Select(vukhi.FieldTenVuKhi).
		Strings(ctx)
}

func (r *vuKhiRepo) GetAllSatThuongCoBanEnt(ctx context.Context) ([]int, error) {
	return r.c.VuKhi.
		Query().
		Select(vukhi.FieldSatThuongCoBan).
		Ints(ctx)
}

func (r *vuKhiRepo) GetAllTocDoDanhEnt(ctx context.Context) ([]float64, error) {
	return r.c.VuKhi.
		Query().
		Select(vukhi.FieldTocDoDanh).
		Float64s(ctx)
}

func (r *vuKhiRepo) GetAllTamDanhEnt(ctx context.Context) ([]int, error) {
	return r.c.VuKhi.
		Query().
		Select(vukhi.FieldTamDanh).
		Ints(ctx)
}

func (r *vuKhiRepo) SearchEnt(ctx context.Context, in VuKhiSearchRequest) ([]*ent.VuKhi, error) {
	qb := r.c.VuKhi.Query().WithLoai().WithDoHiem().WithHe()

	if in.MaVuKhi != nil {
		qb = qb.Where(vukhi.IDEQ(*in.MaVuKhi))
	}
	if strings.TrimSpace(in.TenVuKhi) != "" {
		qb = qb.Where(vukhi.TenVuKhiContainsFold(strings.TrimSpace(in.TenVuKhi)))
	}
	if in.MaLoai != nil {
		qb = qb.Where(vukhi.MaLoaiEQ(*in.MaLoai))
	}
	if strings.TrimSpace(in.TenLoai) != "" {
		qb = qb.Where(vukhi.HasLoaiWith(loaivukhi.TenLoaiContainsFold(strings.TrimSpace(in.TenLoai))))
	}
	if in.MaHe != nil {
		qb = qb.Where(vukhi.MaHeEQ(*in.MaHe))
	}
	if strings.TrimSpace(in.TenHe) != "" {
		qb = qb.Where(vukhi.HasHeWith(he.TenHeContainsFold(strings.TrimSpace(in.TenHe))))
	}
	if in.MaDoHiem != nil {
		qb = qb.Where(vukhi.MaDoHiemEQ(*in.MaDoHiem))
	}
	if strings.TrimSpace(in.TenDoHiem) != "" {
		qb = qb.Where(vukhi.HasDoHiemWith(dohiem.TenDoHiemContainsFold(strings.TrimSpace(in.TenDoHiem))))
	}
	if strings.TrimSpace(in.MauSac) != "" {
		qb = qb.Where(vukhi.HasDoHiemWith(dohiem.MauSacEQ(strings.TrimSpace(in.MauSac))))
	}
	if in.TocDoDanh != nil {
		qb = qb.Where(vukhi.TocDoDanhEQ(*in.TocDoDanh))
	}
	if in.TamDanh != nil {
		qb = qb.Where(vukhi.TamDanhEQ(*in.TamDanh))
	}
	if in.SatThuongBonus != nil {
		qb = qb.Where(vukhi.HasDoHiemWith(dohiem.SatThuongBonusEQ(float64(*in.SatThuongBonus))))
	}
	if in.TocDoDanhBonus != nil {
		qb = qb.Where(vukhi.HasDoHiemWith(dohiem.TocDoDanhBonusEQ(float64(*in.TocDoDanhBonus))))
	}
	if in.MinDamage != nil {
		qb = qb.Where(vukhi.SatThuongCoBanGTE(*in.MinDamage))
	}
	if in.MaxDamage != nil {
		qb = qb.Where(vukhi.SatThuongCoBanLTE(*in.MaxDamage))
	}
	if in.SoLuong != nil {
		qb = qb.Where(vukhi.HasDoHiemWith(dohiem.SoLuongEQ(*in.SoLuong)))
	}

	return qb.All(ctx)
}

// 2. GORM
type VuKhiGormInput struct {
	TenVuKhi       string
	SatThuongCoBan int
	TocDoDanh      float64
	TamDanh        int
	MoTa           *string
	MaLoai         int
	MaDoHiem       int
	MaHe           int
}

type VuKhiRepoGorm interface {
	CreateGorm(ctx context.Context, in VuKhiGormInput) (*db.GVuKhi, error)
	UpdateByNameGorm(ctx context.Context, name string, in VuKhiGormInput) (*db.GVuKhi, error)
	DeleteByNameGorm(ctx context.Context, name string) error

	GetAllGorm(ctx context.Context, q string) ([]db.GVuKhi, error)
	GetAllTenVuKhiGorm(ctx context.Context) ([]string, error)
	GetAllSatThuongCoBanGorm(ctx context.Context) ([]int, error)
	GetAllTocDoDanhGorm(ctx context.Context) ([]float64, error)
	GetAllTamDanhGorm(ctx context.Context) ([]int, error)

	SearchGorm(ctx context.Context, in VuKhiSearchRequest) ([]db.GVuKhi, error)
}

type vuKhiRepoGorm struct{ g *gorm.DB }

func CreateKhoVuKhiGorm(g *gorm.DB) VuKhiRepoGorm { return &vuKhiRepoGorm{g: g} }

func (r *vuKhiRepoGorm) CreateGorm(ctx context.Context, in VuKhiGormInput) (*db.GVuKhi, error) {
	row := db.GVuKhi{
		TenVuKhi:       in.TenVuKhi,
		SatThuongCoBan: in.SatThuongCoBan,
		TocDoDanh:      in.TocDoDanh,
		TamDanh:        in.TamDanh,
		MoTa:           in.MoTa,
		MaLoai:         in.MaLoai,
		MaDoHiem:       in.MaDoHiem,
		MaHe:           in.MaHe,
	}
	if err := r.g.WithContext(ctx).Create(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *vuKhiRepoGorm) GetAllGorm(ctx context.Context, q string) ([]db.GVuKhi, error) {
	var rows []db.GVuKhi
	tx := r.g.WithContext(ctx).Model(&db.GVuKhi{}).Preload("Loai").Preload("DoHiem").Preload("He")
	if q = strings.TrimSpace(q); q != "" {
		tx = tx.Where("LOWER(ten_vu_khi) LIKE ?", "%"+strings.ToLower(q)+"%")
	}
	if err := tx.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *vuKhiRepoGorm) UpdateByNameGorm(ctx context.Context, name string, in VuKhiGormInput) (*db.GVuKhi, error) {
	var row db.GVuKhi
	if err := r.g.WithContext(ctx).First(&row, "ten_vu_khi = ?", name).Error; err != nil {
		return nil, err
	}
	row.TenVuKhi = in.TenVuKhi
	row.SatThuongCoBan = in.SatThuongCoBan
	row.TocDoDanh = in.TocDoDanh
	row.TamDanh = in.TamDanh
	row.MoTa = in.MoTa
	row.MaLoai = in.MaLoai
	row.MaDoHiem = in.MaDoHiem
	row.MaHe = in.MaHe

	if err := r.g.WithContext(ctx).Save(&row).Error; err != nil {
		return nil, err
	}
	if err := r.g.WithContext(ctx).
		Preload("Loai").Preload("DoHiem").Preload("He").
		First(&row, "ten_vu_khi = ?", in.TenVuKhi).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *vuKhiRepoGorm) DeleteByNameGorm(ctx context.Context, name string) error {
	return r.g.WithContext(ctx).Delete(&db.GVuKhi{}, "ten_vu_khi = ?", name).Error
}

func (r *vuKhiRepoGorm) GetAllTenVuKhiGorm(ctx context.Context) ([]string, error) {
	var names []string
	err := r.g.WithContext(ctx).Model(&db.GVuKhi{}).Pluck("ten_vu_khi", &names).Error
	return names, err
}

func (r *vuKhiRepoGorm) GetAllSatThuongCoBanGorm(ctx context.Context) ([]int, error) {
	var values []int
	err := r.g.WithContext(ctx).Model(&db.GVuKhi{}).Pluck("sat_thuong_co_ban", &values).Error
	return values, err
}

func (r *vuKhiRepoGorm) GetAllTocDoDanhGorm(ctx context.Context) ([]float64, error) {
	var values []float64
	err := r.g.WithContext(ctx).Model(&db.GVuKhi{}).Pluck("toc_do_danh", &values).Error
	return values, err
}

func (r *vuKhiRepoGorm) GetAllTamDanhGorm(ctx context.Context) ([]int, error) {
	var values []int
	err := r.g.WithContext(ctx).Model(&db.GVuKhi{}).Pluck("tam_danh", &values).Error
	return values, err
}

func (r *vuKhiRepoGorm) SearchGorm(ctx context.Context, in VuKhiSearchRequest) ([]db.GVuKhi, error) {
	tx := r.g.WithContext(ctx).
		Model(&db.GVuKhi{}).
		// Join sớm để lọc thoải mái theo bảng liên quan
		Joins("JOIN loai_vu_khi l ON l.ma_loai = vu_khi.ma_loai").
		Joins("JOIN do_hiem d ON d.ma_do_hiem = vu_khi.ma_do_hiem").
		Joins("JOIN he h ON h.ma_he = vu_khi.ma_he").
		Preload("Loai").Preload("DoHiem").Preload("He")

	if in.MaVuKhi != nil {
		tx = tx.Where("vu_khi.ma_vu_khi = ?", *in.MaVuKhi)
	}
	if s := strings.TrimSpace(in.TenVuKhi); s != "" {
		tx = tx.Where("LOWER(vu_khi.ten_vu_khi) LIKE ?", "%"+strings.ToLower(s)+"%")
	}
	if in.MaLoai != nil {
		tx = tx.Where("vu_khi.ma_loai = ?", *in.MaLoai)
	}
	if s := strings.TrimSpace(in.TenLoai); s != "" {
		tx = tx.Where("LOWER(l.ten_loai) LIKE ?", "%"+strings.ToLower(s)+"%")
	}
	if in.MaHe != nil {
		tx = tx.Where("vu_khi.ma_he = ?", *in.MaHe)
	}
	if s := strings.TrimSpace(in.TenHe); s != "" {
		tx = tx.Where("LOWER(h.ten_he) LIKE ?", "%"+strings.ToLower(s)+"%")
	}
	if in.MaDoHiem != nil {
		tx = tx.Where("vu_khi.ma_do_hiem = ?", *in.MaDoHiem)
	}
	if s := strings.TrimSpace(in.TenDoHiem); s != "" {
		tx = tx.Where("LOWER(d.ten_do_hiem) LIKE ?", "%"+strings.ToLower(s)+"%")
	}
	if s := strings.TrimSpace(in.MauSac); s != "" {
		tx = tx.Where("LOWER(d.mau_sac) = LOWER(?)", s)
	}
	if in.TocDoDanh != nil {
		tx = tx.Where("vu_khi.toc_do_danh = ?", *in.TocDoDanh)
	}
	if in.TamDanh != nil {
		tx = tx.Where("vu_khi.tam_danh = ?", *in.TamDanh)
	}
	if in.SatThuongBonus != nil {
		tx = tx.Where("d.sat_thuong_bonus = ?", *in.SatThuongBonus)
	}
	if in.TocDoDanhBonus != nil {
		tx = tx.Where("d.toc_do_danh_bonus = ?", *in.TocDoDanhBonus)
	}
	if in.MinDamage != nil {
		tx = tx.Where("vu_khi.sat_thuong_co_ban >= ?", *in.MinDamage)
	}
	if in.MaxDamage != nil {
		tx = tx.Where("vu_khi.sat_thuong_co_ban <= ?", *in.MaxDamage)
	}
	if in.SoLuong != nil {
		tx = tx.Where("d.so_luong = ?", *in.SoLuong)
	}

	var rows []db.GVuKhi
	if err := tx.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// 3. RAW
type VuKhiRaw struct {
	ID             int
	TenVuKhi       string
	SatThuongCoBan int
	TocDoDanh      float64
	TamDanh        int
	MoTa           *string
	MaLoai         int
	MaDoHiem       int
	MaHe           int
}

type VuKhiRawInput struct {
	TenVuKhi       string
	SatThuongCoBan int
	TocDoDanh      float64
	TamDanh        int
	MoTa           *string
	MaLoai         int
	MaDoHiem       int
	MaHe           int
}

type VuKhiRawJoin struct {
	ID             int
	TenVuKhi       string
	SatThuongCoBan int
	TocDoDanh      float64
	TamDanh        int
	MoTa           *string
	MaLoai         int
	MaDoHiem       int
	MaHe           int

	// Trường join thêm:
	TenLoai        string
	TenHe          string
	TenDoHiem      string
	SatThuongBonus float64
	TocDoDanhBonus float64
	MauSac         *string
	SoLuong        *int
}

type VuKhiRepoRaw interface {
	CreateRaw(ctx context.Context, in VuKhiRawInput) (*VuKhiRaw, error)
	GetIdRaw(ctx context.Context, id int) (*VuKhiRaw, error)

	GetByNameRaw(ctx context.Context, name string) (*VuKhiRaw, error)
	UpdateByNameRaw(ctx context.Context, name string, in VuKhiRawInput) (*VuKhiRaw, error)
	DeleteByNameRaw(ctx context.Context, name string) error

	GetAllRaw(ctx context.Context, q string) ([]VuKhiRaw, error)
	GetAllTenVuKhiRaw(ctx context.Context) ([]string, error)
	GetAllSatThuongCoBanRaw(ctx context.Context) ([]int, error)
	GetAllTocDoDanhRaw(ctx context.Context) ([]float64, error)
	GetAllTamDanhRaw(ctx context.Context) ([]int, error)

	SearchRaw(ctx context.Context, in VuKhiSearchRequest) ([]VuKhiRawJoin, error)
}

type vuKhiRepoRaw struct{ sql *sql.DB }

func CreateKhoVuKhiRaw(sqlDB *sql.DB) VuKhiRepoRaw { return &vuKhiRepoRaw{sql: sqlDB} }

func (r *vuKhiRepoRaw) CreateRaw(ctx context.Context, in VuKhiRawInput) (*VuKhiRaw, error) {
	tx, err := r.sql.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var nextID int
	if err = tx.QueryRowContext(ctx, `SELECT IFNULL(MAX(ma_vu_khi)+1,1) FROM vu_khi`).Scan(&nextID); err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO vu_khi
			(ma_vu_khi, ten_vu_khi, sat_thuong_co_ban, toc_do_danh, tam_danh, mo_ta, ma_loai, ma_do_hiem, ma_he)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		nextID, in.TenVuKhi, in.SatThuongCoBan, in.TocDoDanh, in.TamDanh, in.MoTa, in.MaLoai, in.MaDoHiem, in.MaHe,
	)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return r.GetIdRaw(ctx, nextID)
}

func (r *vuKhiRepoRaw) GetIdRaw(ctx context.Context, id int) (*VuKhiRaw, error) {
	row := r.sql.QueryRowContext(ctx, `
		SELECT ma_vu_khi, ten_vu_khi, sat_thuong_co_ban, toc_do_danh, tam_danh, mo_ta, ma_loai, ma_do_hiem, ma_he
		FROM vu_khi WHERE ma_vu_khi = ?`, id)
	var v VuKhiRaw
	if err := row.Scan(&v.ID, &v.TenVuKhi, &v.SatThuongCoBan, &v.TocDoDanh, &v.TamDanh, &v.MoTa, &v.MaLoai, &v.MaDoHiem, &v.MaHe); err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *vuKhiRepoRaw) GetAllRaw(ctx context.Context, q string) ([]VuKhiRaw, error) {
	q = strings.TrimSpace(q)
	var (
		rows *sql.Rows
		err  error
	)
	if q == "" {
		rows, err = r.sql.QueryContext(ctx, `
			SELECT ma_vu_khi, ten_vu_khi, sat_thuong_co_ban, toc_do_danh, tam_danh, mo_ta, ma_loai, ma_do_hiem, ma_he
			FROM vu_khi ORDER BY ma_vu_khi`)
	} else {
		like := "%" + strings.ToLower(q) + "%"
		rows, err = r.sql.QueryContext(ctx, `
			SELECT ma_vu_khi, ten_vu_khi, sat_thuong_co_ban, toc_do_danh, tam_danh, mo_ta, ma_loai, ma_do_hiem, ma_he
			FROM vu_khi
			WHERE LOWER(ten_vu_khi) LIKE ?
			ORDER BY ma_vu_khi`, like)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []VuKhiRaw
	for rows.Next() {
		var v VuKhiRaw
		if err := rows.Scan(&v.ID, &v.TenVuKhi, &v.SatThuongCoBan, &v.TocDoDanh, &v.TamDanh, &v.MoTa, &v.MaLoai, &v.MaDoHiem, &v.MaHe); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func (r *vuKhiRepoRaw) GetByNameRaw(ctx context.Context, name string) (*VuKhiRaw, error) {
	row := r.sql.QueryRowContext(ctx, `
        SELECT ma_vu_khi, ten_vu_khi, sat_thuong_co_ban, toc_do_danh, tam_danh,
               mo_ta, ma_loai, ma_do_hiem, ma_he
        FROM vu_khi WHERE ten_vu_khi = ?`, name)

	var v VuKhiRaw
	if err := row.Scan(&v.ID, &v.TenVuKhi, &v.SatThuongCoBan, &v.TocDoDanh, &v.TamDanh,
		&v.MoTa, &v.MaLoai, &v.MaDoHiem, &v.MaHe); err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *vuKhiRepoRaw) UpdateByNameRaw(ctx context.Context, name string, in VuKhiRawInput) (*VuKhiRaw, error) {
	_, err := r.sql.ExecContext(ctx, `
        UPDATE vu_khi
        SET ten_vu_khi = ?, sat_thuong_co_ban = ?, toc_do_danh = ?, tam_danh = ?,
            mo_ta = ?, ma_loai = ?, ma_do_hiem = ?, ma_he = ?
        WHERE ten_vu_khi = ?`,
		in.TenVuKhi, in.SatThuongCoBan, in.TocDoDanh, in.TamDanh,
		in.MoTa, in.MaLoai, in.MaDoHiem, in.MaHe, name,
	)
	if err != nil {
		return nil, err
	}
	return r.GetByNameRaw(ctx, in.TenVuKhi)
}

func (r *vuKhiRepoRaw) DeleteByNameRaw(ctx context.Context, name string) error {
	_, err := r.sql.ExecContext(ctx, `DELETE FROM vu_khi WHERE ten_vu_khi = ?`, name)
	return err
}

func (r *vuKhiRepoRaw) GetAllTenVuKhiRaw(ctx context.Context) ([]string, error) {
	rows, err := r.sql.QueryContext(ctx, `SELECT ten_vu_khi FROM vu_khi ORDER BY ten_vu_khi`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		names = append(names, s)
	}
	return names, rows.Err()
}

func (r *vuKhiRepoRaw) GetAllSatThuongCoBanRaw(ctx context.Context) ([]int, error) {
	rows, err := r.sql.QueryContext(ctx, `SELECT sat_thuong_co_ban FROM vu_khi`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var values []int
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	return values, rows.Err()
}

func (r *vuKhiRepoRaw) GetAllTocDoDanhRaw(ctx context.Context) ([]float64, error) {
	rows, err := r.sql.QueryContext(ctx, `SELECT toc_do_danh FROM vu_khi`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var values []float64
	for rows.Next() {
		var v float64
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	return values, rows.Err()
}

func (r *vuKhiRepoRaw) GetAllTamDanhRaw(ctx context.Context) ([]int, error) {
	rows, err := r.sql.QueryContext(ctx, `SELECT tam_danh FROM vu_khi`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var values []int
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	return values, rows.Err()
}

func (r *vuKhiRepoRaw) SearchRaw(ctx context.Context, in VuKhiSearchRequest) ([]VuKhiRawJoin, error) {
	var (
		sb   strings.Builder
		args []any
	)

	// Base SELECT + JOIN
	sb.WriteString(`
SELECT
	v.ma_vu_khi,
	v.ten_vu_khi,
	v.sat_thuong_co_ban,
	v.toc_do_danh,
	v.tam_danh,
	v.mo_ta,
	v.ma_loai,
	v.ma_do_hiem,
	v.ma_he,
	l.ten_loai,
	h.ten_he,
	d.ten_do_hiem,
	d.sat_thuong_bonus,
	d.toc_do_danh_bonus,
	d.mau_sac,
	d.so_luong
	FROM vu_khi v
	JOIN loai_vu_khi l ON l.ma_loai = v.ma_loai
	JOIN do_hiem     d ON d.ma_do_hiem = v.ma_do_hiem
	JOIN he          h ON h.ma_he = v.ma_he
	WHERE 1=1
	`)

	// Bộ lọc (giống ENT/GORM)
	if in.MaVuKhi != nil {
		sb.WriteString(" AND v.ma_vu_khi = ?")
		args = append(args, *in.MaVuKhi)
	}
	if s := strings.TrimSpace(in.TenVuKhi); s != "" {
		sb.WriteString(" AND LOWER(v.ten_vu_khi) LIKE ?")
		args = append(args, "%"+strings.ToLower(s)+"%")
	}
	if in.MaLoai != nil {
		sb.WriteString(" AND v.ma_loai = ?")
		args = append(args, *in.MaLoai)
	}
	if s := strings.TrimSpace(in.TenLoai); s != "" {
		sb.WriteString(" AND LOWER(l.ten_loai) LIKE ?")
		args = append(args, "%"+strings.ToLower(s)+"%")
	}
	if in.MaHe != nil {
		sb.WriteString(" AND v.ma_he = ?")
		args = append(args, *in.MaHe)
	}
	if s := strings.TrimSpace(in.TenHe); s != "" {
		sb.WriteString(" AND LOWER(h.ten_he) LIKE ?")
		args = append(args, "%"+strings.ToLower(s)+"%")
	}
	if in.MaDoHiem != nil {
		sb.WriteString(" AND v.ma_do_hiem = ?")
		args = append(args, *in.MaDoHiem)
	}
	if s := strings.TrimSpace(in.TenDoHiem); s != "" {
		sb.WriteString(" AND LOWER(d.ten_do_hiem) LIKE ?")
		args = append(args, "%"+strings.ToLower(s)+"%")
	}
	if s := strings.TrimSpace(in.MauSac); s != "" {
		sb.WriteString(" AND LOWER(d.mau_sac) = LOWER(?)")
		args = append(args, s)
	}
	if in.TocDoDanh != nil {
		sb.WriteString(" AND v.toc_do_danh = ?")
		args = append(args, *in.TocDoDanh)
	}
	if in.TamDanh != nil {
		sb.WriteString(" AND v.tam_danh = ?")
		args = append(args, *in.TamDanh)
	}
	if in.SatThuongBonus != nil {
		sb.WriteString(" AND d.sat_thuong_bonus = ?")
		args = append(args, *in.SatThuongBonus)
	}
	if in.TocDoDanhBonus != nil {
		sb.WriteString(" AND d.toc_do_danh_bonus = ?")
		args = append(args, *in.TocDoDanhBonus)
	}
	if in.MinDamage != nil {
		sb.WriteString(" AND v.sat_thuong_co_ban >= ?")
		args = append(args, *in.MinDamage)
	}
	if in.MaxDamage != nil {
		sb.WriteString(" AND v.sat_thuong_co_ban <= ?")
		args = append(args, *in.MaxDamage)
	}
	if in.SoLuong != nil {
		sb.WriteString(" AND d.so_luong = ?")
		args = append(args, *in.SoLuong)
	}

	// Sắp xếp mặc định
	sb.WriteString(" ORDER BY v.ma_vu_khi")

	rows, err := r.sql.QueryContext(ctx, sb.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []VuKhiRawJoin
	for rows.Next() {
		var v VuKhiRawJoin
		if err := rows.Scan(
			&v.ID, &v.TenVuKhi, &v.SatThuongCoBan, &v.TocDoDanh, &v.TamDanh, &v.MoTa,
			&v.MaLoai, &v.MaDoHiem, &v.MaHe,
			&v.TenLoai, &v.TenHe, &v.TenDoHiem,
			&v.SatThuongBonus, &v.TocDoDanhBonus, &v.MauSac, &v.SoLuong,
		); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
