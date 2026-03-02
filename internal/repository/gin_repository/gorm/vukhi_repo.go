package gin_repository

import (
	"context"
	"errors"
	"game/internal/db"
	mappers "game/internal/mapping/gorm"
	"game/internal/models"
	"strings"

	"gorm.io/gorm"
)

type VuKhiRepositoryGorm struct {
	db *gorm.DB
}

func NewVuKhiRepositoryGorm(db *gorm.DB) VuKhiRepoGorm {
	return &VuKhiRepositoryGorm{
		db: db,
	}
}

// Repository: Tạo vũ khí
func (r *VuKhiRepositoryGorm) CreateGorm(ctx context.Context, vk *models.VuKhis) (*models.VuKhis, error) {

	vukhis := mappers.ModelToGormVuKhi(vk)
	vukhis.Version = 1

	if err := r.db.WithContext(ctx).Create(vukhis).Error; err != nil {
		return nil, err
	}

	return mappers.GormVuKhiToModel(vukhis), nil
}

func (r *VuKhiRepositoryGorm) UpdateByNameGorm(ctx context.Context, name string, in models.VuKhiUpdateInput) (*models.VuKhis, error) {

	var vk db.GVuKhi

	if err := r.db.WithContext(ctx).
		Where("ten_vu_khi = ?", name).
		First(&vk).Error; err != nil {
		return nil, err
	}

	if vk.Version != in.Version {
		return nil, errors.New("data has been modified")
	}

	vk.TenVuKhi = in.TenVuKhi
	vk.SatThuongCoBan = in.SatThuongCoBan
	vk.TocDoDanh = in.TocDoDanh
	vk.TamDanh = in.TamDanh
	vk.MoTa = &in.MoTa
	vk.MaLoaiVuKhi = in.MaLoaiVuKhi
	vk.MaDoHiem = in.MaDoHiem
	vk.MaHe = in.MaHe
	vk.Version++

	if err := r.db.WithContext(ctx).Save(&vk).Error; err != nil {
		return nil, err
	}

	return mappers.GormVuKhiToModel(&vk), nil
}

func (r *VuKhiRepositoryGorm) DeleteByNameGorm(ctx context.Context, name string) error {
	return r.db.WithContext(ctx).
		Where("ten_vu_khi = ?", name).
		Delete(&db.GVuKhi{}).Error
}

func (r *VuKhiRepositoryGorm) GetAllGorm(ctx context.Context, kw string) ([]*models.VuKhis, error) {

	var entities []db.GVuKhi
	query := r.db.WithContext(ctx)

	if kw != "" {
		query = query.Where("ten_vu_khi LIKE ?", "%"+kw+"%")
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}

	result := make([]*models.VuKhis, 0, len(entities))

	for i := range entities {
		result = append(result, mappers.GormVuKhiToModel(&entities[i]))
	}

	return result, nil
}
func (r *VuKhiRepositoryGorm) GetAllTenVuKhiGorm(ctx context.Context) ([]string, error) {
	var result []string

	err := r.db.WithContext(ctx).
		Model(&db.GVuKhi{}).
		Pluck("ten_vu_khi", &result).Error

	return result, err
}

func (r *VuKhiRepositoryGorm) GetAllSatThuongCoBanGorm(ctx context.Context) ([]int, error) {
	var result []int

	err := r.db.WithContext(ctx).
		Model(&db.GVuKhi{}).
		Pluck("sat_thuong_co_ban", &result).Error

	return result, err
}

func (r *VuKhiRepositoryGorm) GetAllTocDoDanhGorm(ctx context.Context) ([]float64, error) {
	var result []float64

	err := r.db.WithContext(ctx).
		Model(&db.GVuKhi{}).
		Pluck("toc_do_danh", &result).Error

	return result, err
}

func (r *VuKhiRepositoryGorm) GetAllTamDanhGorm(ctx context.Context) ([]int, error) {
	var result []int

	err := r.db.WithContext(ctx).
		Model(&db.GVuKhi{}).
		Pluck("tam_danh", &result).Error

	return result, err
}

func (r *VuKhiRepositoryGorm) SearchGorm(
	ctx context.Context,
	req models.VuKhiSearchRequest,
	limit int,
	offset int,
	arrangeAsc string,
	arrangeDesc string,
) ([]*models.VuKhiResponse, int64, error) {

	var list []db.GVuKhi
	var total int64

	dbQuery := r.db.WithContext(ctx).
		Model(&db.GVuKhi{}).
		Preload("Loai").
		Preload("He").
		Preload("DoHiem").
		Joins("LEFT JOIN g_loai_vu_khis ON g_loai_vu_khis.ma_loai_vu_khi = g_vu_khis.ma_loai_vu_khi").
		Joins("LEFT JOIN g_hes ON g_hes.ma_he = g_vu_khis.ma_he").
		Joins("LEFT JOIN g_do_hiems ON g_do_hiems.ma_do_hiem = g_vu_khis.ma_do_hiem")

	//  FILTER

	if req.MaVuKhi != nil {
		dbQuery = dbQuery.Where("g_vu_khis.ma_vu_khi = ?", *req.MaVuKhi)
	}

	if t := strings.TrimSpace(req.TenVuKhi); t != "" {
		dbQuery = dbQuery.Where("LOWER(g_vu_khis.ten_vu_khi) LIKE LOWER(?)", "%"+t+"%")
	}

	if req.MaLoaiVuKhi != nil {
		dbQuery = dbQuery.Where("g_vu_khis.ma_loai_vu_khi = ?", *req.MaLoaiVuKhi)
	}

	if strings.TrimSpace(req.TenLoai) != "" {
		dbQuery = dbQuery.Where("LOWER(g_loai_vu_khis.ten_loai_vu_khi) LIKE LOWER(?)",
			"%"+strings.TrimSpace(req.TenLoai)+"%")
	}

	if req.MaHe != nil {
		dbQuery = dbQuery.Where("g_vu_khis.ma_he = ?", *req.MaHe)
	}

	if strings.TrimSpace(req.TenHe) != "" {
		dbQuery = dbQuery.Where("LOWER(g_hes.ten_he) LIKE LOWER(?)",
			"%"+strings.TrimSpace(req.TenHe)+"%")
	}

	if req.MaDoHiem != nil {
		dbQuery = dbQuery.Where("g_vu_khis.ma_do_hiem = ?", *req.MaDoHiem)
	}

	if strings.TrimSpace(req.TenDoHiem) != "" {
		dbQuery = dbQuery.Where("LOWER(g_do_hiems.ten_do_hiem) LIKE LOWER(?)",
			"%"+strings.TrimSpace(req.TenDoHiem)+"%")
	}

	if strings.TrimSpace(req.MauSac) != "" {
		dbQuery = dbQuery.Where("LOWER(g_do_hiems.mau_sac) LIKE LOWER(?)",
			"%"+strings.TrimSpace(req.MauSac)+"%")
	}

	if req.SatThuongCoBan != nil {
		dbQuery = dbQuery.Where("g_vu_khis.sat_thuong_co_ban = ?", *req.SatThuongCoBan)
	}

	if req.MinDamage != nil {
		dbQuery = dbQuery.Where("g_vu_khis.sat_thuong_co_ban >= ?", *req.MinDamage)
	}

	if req.MaxDamage != nil {
		dbQuery = dbQuery.Where("g_vu_khis.sat_thuong_co_ban <= ?", *req.MaxDamage)
	}

	if req.TocDoDanh != nil {
		dbQuery = dbQuery.Where("g_vu_khis.toc_do_danh = ?", *req.TocDoDanh)
	}

	if req.TamDanh != nil {
		dbQuery = dbQuery.Where("g_vu_khis.tam_danh = ?", *req.TamDanh)
	}

	if req.SoLuong != nil {
		dbQuery = dbQuery.Where("g_do_hiems.so_luong = ?", *req.SoLuong)
	}

	if req.SatThuongBonus != nil {
		dbQuery = dbQuery.Where("g_do_hiems.sat_thuong_bonus = ?", *req.SatThuongBonus)
	}

	if req.TocDoDanhBonus != nil {
		dbQuery = dbQuery.Where("g_do_hiems.toc_do_danh_bonus = ?", *req.TocDoDanhBonus)
	}

	//  COUNT
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	//  SORT

	if arrangeAsc != "" {

		dbQuery = dbQuery.Order(arrangeAsc + " ASC")

	} else if arrangeDesc != "" {

		dbQuery = dbQuery.Order(arrangeDesc + " DESC")

	} else {
		dbQuery = dbQuery.Order("g_vu_khis.ma_vu_khi ASC")
	}

	//  PAGINATION
	if limit > 0 {
		dbQuery = dbQuery.Limit(limit).Offset(offset)
	}

	//  QUERY
	if err := dbQuery.Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return mappers.ListVuKhiToResponsesSearchGORM(list), total, nil
}

// // // 3. RAW
// // type VuKhiRaw struct {
// // 	MaVuKhi             int
// // 	TenVuKhi       string
// // 	SatThuongCoBan int
// // 	TocDoDanh      float64
// // 	TamDanh        int
// // 	MoTa           *string
// // 	MaLoai         int
// // 	MaDoHiem       int
// // 	MaHe           int
// // }

// // type VuKhiRawInput struct {
// // 	TenVuKhi       string
// // 	SatThuongCoBan int
// // 	TocDoDanh      float64
// // 	TamDanh        int
// // 	MoTa           *string
// // 	MaLoai         int
// // 	MaDoHiem       int
// // 	MaHe           int
// // }

// // type VuKhiRawJoin struct {
// // 	MaVuKhi             int
// // 	TenVuKhi       string
// // 	SatThuongCoBan int
// // 	TocDoDanh      float64
// // 	TamDanh        int
// // 	MoTa           *string
// // 	MaLoai         int
// // 	MaDoHiem       int
// // 	MaHe           int

// // 	// Trường join thêm:
// // 	TenLoai        string
// // 	TenHe          string
// // 	TenDoHiem      string
// // 	SatThuongBonus float64
// // 	TocDoDanhBonus float64
// // 	MauSac         *string
// // 	SoLuong        *int
// // }

// // type VuKhiRepoRaw interface {
// // 	CreateRaw(ctx context.Context, in VuKhiRawInput) (*VuKhiRaw, error)
// // 	GetIdRaw(ctx context.Context, id int) (*VuKhiRaw, error)

// // 	GetByNameRaw(ctx context.Context, name string) (*VuKhiRaw, error)
// // 	UpdateByNameRaw(ctx context.Context, name string, in VuKhiRawInput) (*VuKhiRaw, error)
// // 	DeleteByNameRaw(ctx context.Context, name string) error

// // 	GetAllRaw(ctx context.Context, q string) ([]VuKhiRaw, error)
// // 	GetAllTenVuKhiRaw(ctx context.Context) ([]string, error)
// // 	GetAllSatThuongCoBanRaw(ctx context.Context) ([]int, error)
// // 	GetAllTocDoDanhRaw(ctx context.Context) ([]float64, error)
// // 	GetAllTamDanhRaw(ctx context.Context) ([]int, error)

// // 	SearchRaw(ctx context.Context, in VuKhiSearchRequest) ([]VuKhiRawJoin, error)
// // }

// // type vuKhiRepoRaw struct{ sql *sql.DB }

// // func CreateKhoVuKhiRaw(sqlDB *sql.DB) VuKhiRepoRaw { return &vuKhiRepoRaw{sql: sqlDB} }

// // func (r *vuKhiRepoRaw) CreateRaw(ctx context.Context, in VuKhiRawInput) (*VuKhiRaw, error) {
// // 	tx, err := r.sql.BeginTx(ctx, nil)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer func() {
// // 		if err != nil {
// // 			_ = tx.Rollback()
// // 		}
// // 	}()

// // 	var nextID int
// // 	if err = tx.QueryRowContext(ctx, `SELECT IFNULL(MAX(ma_vu_khi)+1,1) FROM vu_khi`).Scan(&nextID); err != nil {
// // 		return nil, err
// // 	}

// // 	_, err = tx.ExecContext(ctx, `
// // 		INSERT INTO vu_khi
// // 			(ma_vu_khi, ten_vu_khi, sat_thuong_co_ban, toc_do_danh, tam_danh, mo_ta, ma_loai_vu_khi, ma_do_hiem, ma_he)
// // 		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
// // 		nextID, in.TenVuKhi, in.SatThuongCoBan, in.TocDoDanh, in.TamDanh, in.MoTa, in.MaLoai, in.MaDoHiem, in.MaHe,
// // 	)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	if err = tx.Commit(); err != nil {
// // 		return nil, err
// // 	}
// // 	return r.GetIdRaw(ctx, nextID)
// // }

// // func (r *vuKhiRepoRaw) GetIdRaw(ctx context.Context, id int) (*VuKhiRaw, error) {
// // 	row := r.sql.QueryRowContext(ctx, `
// // 		SELECT ma_vu_khi, ten_vu_khi, sat_thuong_co_ban, toc_do_danh, tam_danh, mo_ta, ma_loai_vu_khi, ma_do_hiem, ma_he
// // 		FROM vu_khi WHERE ma_vu_khi = ?`, id)
// // 	var v VuKhiRaw
// // 	if err := row.Scan(&v.MaVuKhi, &v.TenVuKhi, &v.SatThuongCoBan, &v.TocDoDanh, &v.TamDanh, &v.MoTa, &v.MaLoai, &v.MaDoHiem, &v.MaHe); err != nil {
// // 		return nil, err
// // 	}
// // 	return &v, nil
// // }

// // func (r *vuKhiRepoRaw) GetAllRaw(ctx context.Context, q string) ([]VuKhiRaw, error) {
// // 	q = strings.TrimSpace(q)
// // 	var (
// // 		rows *sql.Rows
// // 		err  error
// // 	)
// // 	if q == "" {
// // 		rows, err = r.sql.QueryContext(ctx, `
// // 			SELECT ma_vu_khi, ten_vu_khi, sat_thuong_co_ban, toc_do_danh, tam_danh, mo_ta, ma_loai_vu_khi, ma_do_hiem, ma_he
// // 			FROM vu_khi ORDER BY ma_vu_khi`)
// // 	} else {
// // 		like := "%" + strings.ToLower(q) + "%"
// // 		rows, err = r.sql.QueryContext(ctx, `
// // 			SELECT ma_vu_khi, ten_vu_khi, sat_thuong_co_ban, toc_do_danh, tam_danh, mo_ta, ma_loai_vu_khi, ma_do_hiem, ma_he
// // 			FROM vu_khi
// // 			WHERE LOWER(ten_vu_khi) LIKE ?
// // 			ORDER BY ma_vu_khi`, like)
// // 	}
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer rows.Close()

// // 	var out []VuKhiRaw
// // 	for rows.Next() {
// // 		var v VuKhiRaw
// // 		if err := rows.Scan(&v.MaVuKhi, &v.TenVuKhi, &v.SatThuongCoBan, &v.TocDoDanh, &v.TamDanh, &v.MoTa, &v.MaLoai, &v.MaDoHiem, &v.MaHe); err != nil {
// // 			return nil, err
// // 		}
// // 		out = append(out, v)
// // 	}
// // 	return out, rows.Err()
// // }

// // func (r *vuKhiRepoRaw) GetByNameRaw(ctx context.Context, name string) (*VuKhiRaw, error) {
// // 	row := r.sql.QueryRowContext(ctx, `
// //         SELECT ma_vu_khi, ten_vu_khi, sat_thuong_co_ban, toc_do_danh, tam_danh,
// //                mo_ta, ma_loai_vu_khi, ma_do_hiem, ma_he
// //         FROM vu_khi WHERE ten_vu_khi = ?`, name)

// // 	var v VuKhiRaw
// // 	if err := row.Scan(&v.MaVuKhi, &v.TenVuKhi, &v.SatThuongCoBan, &v.TocDoDanh, &v.TamDanh,
// // 		&v.MoTa, &v.MaLoai, &v.MaDoHiem, &v.MaHe); err != nil {
// // 		return nil, err
// // 	}
// // 	return &v, nil
// // }

// // func (r *vuKhiRepoRaw) UpdateByNameRaw(ctx context.Context, name string, in VuKhiRawInput) (*VuKhiRaw, error) {
// // 	_, err := r.sql.ExecContext(ctx, `
// //         UPDATE vu_khi
// //         SET ten_vu_khi = ?, sat_thuong_co_ban = ?, toc_do_danh = ?, tam_danh = ?,
// //             mo_ta = ?, ma_loai_vu_khi = ?, ma_do_hiem = ?, ma_he = ?
// //         WHERE ten_vu_khi = ?`,
// // 		in.TenVuKhi, in.SatThuongCoBan, in.TocDoDanh, in.TamDanh,
// // 		in.MoTa, in.MaLoai, in.MaDoHiem, in.MaHe, name,
// // 	)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	return r.GetByNameRaw(ctx, in.TenVuKhi)
// // }

// // func (r *vuKhiRepoRaw) DeleteByNameRaw(ctx context.Context, name string) error {
// // 	_, err := r.sql.ExecContext(ctx, `DELETE FROM vu_khi WHERE ten_vu_khi = ?`, name)
// // 	return err
// // }

// // func (r *vuKhiRepoRaw) GetAllTenVuKhiRaw(ctx context.Context) ([]string, error) {
// // 	rows, err := r.sql.QueryContext(ctx, `SELECT ten_vu_khi FROM vu_khi ORDER BY ten_vu_khi`)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer rows.Close()
// // 	var names []string
// // 	for rows.Next() {
// // 		var s string
// // 		if err := rows.Scan(&s); err != nil {
// // 			return nil, err
// // 		}
// // 		names = append(names, s)
// // 	}
// // 	return names, rows.Err()
// // }

// // func (r *vuKhiRepoRaw) GetAllSatThuongCoBanRaw(ctx context.Context) ([]int, error) {
// // 	rows, err := r.sql.QueryContext(ctx, `SELECT sat_thuong_co_ban FROM vu_khi`)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer rows.Close()
// // 	var values []int
// // 	for rows.Next() {
// // 		var v int
// // 		if err := rows.Scan(&v); err != nil {
// // 			return nil, err
// // 		}
// // 		values = append(values, v)
// // 	}
// // 	return values, rows.Err()
// // }

// // func (r *vuKhiRepoRaw) GetAllTocDoDanhRaw(ctx context.Context) ([]float64, error) {
// // 	rows, err := r.sql.QueryContext(ctx, `SELECT toc_do_danh FROM vu_khi`)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer rows.Close()
// // 	var values []float64
// // 	for rows.Next() {
// // 		var v float64
// // 		if err := rows.Scan(&v); err != nil {
// // 			return nil, err
// // 		}
// // 		values = append(values, v)
// // 	}
// // 	return values, rows.Err()
// // }

// // func (r *vuKhiRepoRaw) GetAllTamDanhRaw(ctx context.Context) ([]int, error) {
// // 	rows, err := r.sql.QueryContext(ctx, `SELECT tam_danh FROM vu_khi`)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer rows.Close()
// // 	var values []int
// // 	for rows.Next() {
// // 		var v int
// // 		if err := rows.Scan(&v); err != nil {
// // 			return nil, err
// // 		}
// // 		values = append(values, v)
// // 	}
// // 	return values, rows.Err()
// // }

// // func (r *vuKhiRepoRaw) SearchRaw(ctx context.Context, in VuKhiSearchRequest) ([]VuKhiRawJoin, error) {
// // 	var (
// // 		sb   strings.Builder
// // 		args []any
// // 	)

// // 	// Base SELECT + JOIN
// // 	sb.WriteString(`
// // SELECT
// // 	v.ma_vu_khi,
// // 	v.ten_vu_khi,
// // 	v.sat_thuong_co_ban,
// // 	v.toc_do_danh,
// // 	v.tam_danh,
// // 	v.mo_ta,
// // 	v.ma_loai_vu_khi,
// // 	v.ma_do_hiem,
// // 	v.ma_he,
// // 	l.ten_loai,
// // 	h.ten_he,
// // 	d.ten_do_hiem,
// // 	d.sat_thuong_bonus,
// // 	d.toc_do_danh_bonus,
// // 	d.mau_sac,
// // 	d.so_luong
// // 	FROM vu_khi v
// // 	JOIN loai_vu_khi l ON l.ma_loai_vu_khi = v.ma_loai_vu_khi
// // 	JOIN do_hiem     d ON d.ma_do_hiem = v.ma_do_hiem
// // 	JOIN he          h ON h.ma_he = v.ma_he
// // 	WHERE 1=1
// // 	`)

// // 	// Bộ lọc (giống ENT/GORM)
// // 	if in.MaVuKhi != nil {
// // 		sb.WriteString(" AND v.ma_vu_khi = ?")
// // 		args = append(args, *in.MaVuKhi)
// // 	}
// // 	if s := strings.TrimSpace(in.TenVuKhi); s != "" {
// // 		sb.WriteString(" AND LOWER(v.ten_vu_khi) LIKE ?")
// // 		args = append(args, "%"+strings.ToLower(s)+"%")
// // 	}
// // 	if in.MaLoai != nil {
// // 		sb.WriteString(" AND v.ma_loai_vu_khi = ?")
// // 		args = append(args, *in.MaLoai)
// // 	}
// // 	if s := strings.TrimSpace(in.TenLoai); s != "" {
// // 		sb.WriteString(" AND LOWER(l.ten_loai) LIKE ?")
// // 		args = append(args, "%"+strings.ToLower(s)+"%")
// // 	}
// // 	if in.MaHe != nil {
// // 		sb.WriteString(" AND v.ma_he = ?")
// // 		args = append(args, *in.MaHe)
// // 	}
// // 	if s := strings.TrimSpace(in.TenHe); s != "" {
// // 		sb.WriteString(" AND LOWER(h.ten_he) LIKE ?")
// // 		args = append(args, "%"+strings.ToLower(s)+"%")
// // 	}
// // 	if in.MaDoHiem != nil {
// // 		sb.WriteString(" AND v.ma_do_hiem = ?")
// // 		args = append(args, *in.MaDoHiem)
// // 	}
// // 	if s := strings.TrimSpace(in.TenDoHiem); s != "" {
// // 		sb.WriteString(" AND LOWER(d.ten_do_hiem) LIKE ?")
// // 		args = append(args, "%"+strings.ToLower(s)+"%")
// // 	}
// // 	if s := strings.TrimSpace(in.MauSac); s != "" {
// // 		sb.WriteString(" AND LOWER(d.mau_sac) = LOWER(?)")
// // 		args = append(args, s)
// // 	}
// // 	if in.TocDoDanh != nil {
// // 		sb.WriteString(" AND v.toc_do_danh = ?")
// // 		args = append(args, *in.TocDoDanh)
// // 	}
// // 	if in.TamDanh != nil {
// // 		sb.WriteString(" AND v.tam_danh = ?")
// // 		args = append(args, *in.TamDanh)
// // 	}
// // 	if in.SatThuongBonus != nil {
// // 		sb.WriteString(" AND d.sat_thuong_bonus = ?")
// // 		args = append(args, *in.SatThuongBonus)
// // 	}
// // 	if in.TocDoDanhBonus != nil {
// // 		sb.WriteString(" AND d.toc_do_danh_bonus = ?")
// // 		args = append(args, *in.TocDoDanhBonus)
// // 	}
// // 	if in.MinDamage != nil {
// // 		sb.WriteString(" AND v.sat_thuong_co_ban >= ?")
// // 		args = append(args, *in.MinDamage)
// // 	}
// // 	if in.MaxDamage != nil {
// // 		sb.WriteString(" AND v.sat_thuong_co_ban <= ?")
// // 		args = append(args, *in.MaxDamage)
// // 	}
// // 	if in.SoLuong != nil {
// // 		sb.WriteString(" AND d.so_luong = ?")
// // 		args = append(args, *in.SoLuong)
// // 	}

// // 	// Sắp xếp mặc định
// // 	sb.WriteString(" ORDER BY v.ma_vu_khi")

// // 	rows, err := r.sql.QueryContext(ctx, sb.String(), args...)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer rows.Close()

// // 	var out []VuKhiRawJoin
// // 	for rows.Next() {
// // 		var v VuKhiRawJoin
// // 		if err := rows.Scan(
// // 			&v.MaVuKhi, &v.TenVuKhi, &v.SatThuongCoBan, &v.TocDoDanh, &v.TamDanh, &v.MoTa,
// // 			&v.MaLoai, &v.MaDoHiem, &v.MaHe,
// // 			&v.TenLoai, &v.TenHe, &v.TenDoHiem,
// // 			&v.SatThuongBonus, &v.TocDoDanhBonus, &v.MauSac, &v.SoLuong,
// // 		); err != nil {
// // 			return nil, err
// // 		}
// // 		out = append(out, v)
// // 	}
// // 	return out, rows.Err()
// // }
