package gin_repository

// import (
// 	"context"

// 	"game/ent"
// 	"game/ent/dohiem"
// )

// // 1. ENT
// type TaoDoHiemInput struct {
// 	TenDoHiem      string
// 	MauSac         string
// 	SoLuong        int
// 	SatThuongBonus float64
// 	TocDoDanhBonus float64
// }

// type doHiemRepo struct {
// 	c *ent.Client
// }

// func CreateKhoDoHiemEnt(c *ent.Client) DoHiemRepo {
// 	return &doHiemRepo{c: c}
// }

// func (r *doHiemRepo) CreateEnt(ctx context.Context, in TaoDoHiemInput) (*ent.DoHiem, error) {
// 	return r.c.DoHiem.Create().
// 		SetTenDoHiem(in.TenDoHiem).
// 		SetMauSac(in.MauSac).
// 		SetSoLuong(in.SoLuong).
// 		SetSatThuongBonus(in.SatThuongBonus).
// 		SetTocDoDanhBonus(in.TocDoDanhBonus).
// 		Save(ctx)
// }

// func (r *doHiemRepo) GetAllEnt(ctx context.Context, q string) ([]*ent.DoHiem, error) {
// 	qb := r.c.DoHiem.Query()
// 	if q != "" {
// 		qb = qb.Where(dohiem.TenDoHiemContainsFold(q))
// 	}
// 	return qb.All(ctx)
// }

// func (r *doHiemRepo) GetByNameEnt(ctx context.Context, name string) (*ent.DoHiem, error) {
// 	return r.c.DoHiem.Query().Where(dohiem.TenDoHiemEQ(name)).Only(ctx)
// }

// func (r *doHiemRepo) UpdateByNameEnt(ctx context.Context, name string, mutate func(*ent.DoHiemUpdateOne) *ent.DoHiemUpdateOne) (*ent.DoHiem, error) {
// 	cur, err := r.c.DoHiem.Query().Where(dohiem.TenDoHiemEQ(name)).Only(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	up := r.c.DoHiem.UpdateOneID(cur.ID)
// 	return mutate(up).Save(ctx)
// }

// func (r *doHiemRepo) DeleteByNameEnt(ctx context.Context, name string) error {
// 	cur, err := r.c.DoHiem.Query().Where(dohiem.TenDoHiemEQ(name)).Only(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	return r.c.DoHiem.DeleteOneID(cur.ID).Exec(ctx)
// }

// func (r *doHiemRepo) GetAllTen(ctx context.Context) ([]string, error) {
// 	return r.c.DoHiem.
// 		Query().
// 		Select(dohiem.FieldTenDoHiem).
// 		Strings(ctx)
// }

// func (r *doHiemRepo) GetAllSoLuong(ctx context.Context) ([]int, error) {
// 	return r.c.DoHiem.
// 		Query().
// 		Select(dohiem.FieldSoLuong).
// 		Ints(ctx)
// }

// func (r *doHiemRepo) GetAllMauSac(ctx context.Context) ([]string, error) {
// 	return r.c.DoHiem.
// 		Query().
// 		Select(dohiem.FieldMauSac).
// 		Strings(ctx)
// }

// func (r *doHiemRepo) GetAllSatThuongBonus(ctx context.Context) ([]float64, error) {
// 	return r.c.DoHiem.
// 		Query().
// 		Select(dohiem.FieldSatThuongBonus).
// 		Float64s(ctx)
// }

// func (r *doHiemRepo) GetAllTocDoDanhBonus(ctx context.Context) ([]float64, error) {
// 	return r.c.DoHiem.
// 		Query().
// 		Select(dohiem.FieldTocDoDanhBonus).
// 		Float64s(ctx)
// }

// // // 2. GORM
// // type DoHiemGormInput struct {
// // 	TenDoHiem      string
// // 	MauSac         string
// // 	SoLuong        int
// // 	SatThuongBonus float64
// // 	TocDoDanhBonus float64
// // }

// // type DoHiemRepoGorm interface {
// // 	CreateGorm(ctx context.Context, in DoHiemGormInput) (*db.GDoHiem, error)
// // 	GetIdGorm(ctx context.Context, id int) (*db.GDoHiem, error)
// // 	GetAllGorm(ctx context.Context, q string) ([]db.GDoHiem, error)
// // 	UpdateGorm(ctx context.Context, id int, in DoHiemGormInput) (*db.GDoHiem, error)
// // 	DeleteGorm(ctx context.Context, id int) error

// // 	GetByNameGorm(ctx context.Context, name string) (*db.GDoHiem, error)
// // 	UpdateByNameGorm(ctx context.Context, name string, in DoHiemGormInput) (*db.GDoHiem, error)
// // 	DeleteByNameGorm(ctx context.Context, name string) error

// // 	GetAllTenGorm(ctx context.Context) ([]string, error)
// // 	GetAllSoLuongGorm(ctx context.Context) ([]int, error)
// // 	GetAllMauSacGorm(ctx context.Context) ([]string, error)
// // 	GetAllSatThuongBonusGorm(ctx context.Context) ([]float64, error)
// // 	GetAllTocDoDanhBonusGorm(ctx context.Context) ([]float64, error)
// // }

// // type doHiemRepoGorm struct{ g *gorm.DB }

// // func CreateKhoDoHiemGorm(g *gorm.DB) DoHiemRepoGorm { return &doHiemRepoGorm{g: g} }

// // func (r *doHiemRepoGorm) CreateGorm(ctx context.Context, in DoHiemGormInput) (*db.GDoHiem, error) {
// // 	row := db.GDoHiem{
// // 		TenDoHiem:      in.TenDoHiem,
// // 		MauSac:         in.MauSac,
// // 		SoLuong:        in.SoLuong,
// // 		SatThuongBonus: in.SatThuongBonus,
// // 		TocDoDanhBonus: in.TocDoDanhBonus,
// // 	}
// // 	if err := r.g.WithContext(ctx).Create(&row).Error; err != nil {
// // 		return nil, err
// // 	}
// // 	return &row, nil
// // }

// // func (r *doHiemRepoGorm) GetIdGorm(ctx context.Context, id int) (*db.GDoHiem, error) {
// // 	var row db.GDoHiem
// // 	if err := r.g.WithContext(ctx).First(&row, "ma_do_hiem = ?", id).Error; err != nil {
// // 		return nil, err
// // 	}
// // 	return &row, nil
// // }

// // func (r *doHiemRepoGorm) GetAllGorm(ctx context.Context, q string) ([]db.GDoHiem, error) {
// // 	var rows []db.GDoHiem
// // 	tx := r.g.WithContext(ctx).Model(&db.GDoHiem{})
// // 	if q = strings.TrimSpace(q); q != "" {
// // 		tx = tx.Where("LOWER(ten_do_hiem) LIKE ?", "%"+strings.ToLower(q)+"%")
// // 	}
// // 	if err := tx.Find(&rows).Error; err != nil {
// // 		return nil, err
// // 	}
// // 	return rows, nil
// // }

// // func (r *doHiemRepoGorm) UpdateGorm(ctx context.Context, id int, in DoHiemGormInput) (*db.GDoHiem, error) {
// // 	var row db.GDoHiem
// // 	if err := r.g.WithContext(ctx).First(&row, "ma_do_hiem = ?", id).Error; err != nil {
// // 		return nil, err
// // 	}
// // 	row.TenDoHiem = in.TenDoHiem
// // 	row.MauSac = in.MauSac
// // 	row.SoLuong = in.SoLuong
// // 	row.SatThuongBonus = in.SatThuongBonus
// // 	row.TocDoDanhBonus = in.TocDoDanhBonus

// // 	if err := r.g.WithContext(ctx).Save(&row).Error; err != nil {
// // 		return nil, err
// // 	}
// // 	return &row, nil
// // }

// // func (r *doHiemRepoGorm) DeleteGorm(ctx context.Context, id int) error {
// // 	return r.g.WithContext(ctx).Delete(&db.GDoHiem{}, "ma_do_hiem = ?", id).Error
// // }

// // type DoHiemRaw struct {
// // 	ID             int
// // 	TenDoHiem      string
// // 	MauSac         *string
// // 	SoLuong        int
// // 	SatThuongBonus float64
// // 	TocDoDanhBonus float64
// // }

// // func (r *doHiemRepoGorm) GetByNameGorm(ctx context.Context, name string) (*db.GDoHiem, error) {
// // 	var row db.GDoHiem
// // 	if err := r.g.WithContext(ctx).Where("ten_do_hiem = ?", name).First(&row).Error; err != nil {
// // 		return nil, err
// // 	}
// // 	return &row, nil
// // }
// // func (r *doHiemRepoGorm) UpdateByNameGorm(ctx context.Context, name string, in DoHiemGormInput) (*db.GDoHiem, error) {
// // 	var row db.GDoHiem
// // 	if err := r.g.WithContext(ctx).Where("ten_do_hiem = ?", name).First(&row).Error; err != nil {
// // 		return nil, err
// // 	}
// // 	updates := map[string]interface{}{
// // 		"ten_do_hiem":       in.TenDoHiem,
// // 		"mau_sac":           in.MauSac,
// // 		"so_luong":          in.SoLuong,
// // 		"sat_thuong_bonus":  in.SatThuongBonus,
// // 		"toc_do_danh_bonus": in.TocDoDanhBonus,
// // 	}
// // 	if err := r.g.WithContext(ctx).Model(&row).Updates(updates).Error; err != nil {
// // 		return nil, err
// // 	}
// // 	// Reload updated record
// // 	if err := r.g.WithContext(ctx).First(&row, row.ID).Error; err != nil {
// // 		return nil, err
// // 	}
// // 	return &row, nil
// // }
// // func (r *doHiemRepoGorm) DeleteByNameGorm(ctx context.Context, name string) error {
// // 	res := r.g.WithContext(ctx).Where("ten_do_hiem = ?", name).Delete(&db.GDoHiem{})
// // 	if res.Error != nil {
// // 		return res.Error
// // 	}
// // 	if res.RowsAffected == 0 {
// // 		return gorm.ErrRecordNotFound
// // 	}
// // 	return nil
// // }

// // func (r *doHiemRepoGorm) GetAllTenGorm(ctx context.Context) ([]string, error) {
// // 	var items []string
// // 	err := r.g.WithContext(ctx).
// // 		Model(&db.GDoHiem{}).
// // 		Pluck("ten_do_hiem", &items).Error
// // 	return items, err
// // }

// // func (r *doHiemRepoGorm) GetAllSoLuongGorm(ctx context.Context) ([]int, error) {
// // 	var items []int
// // 	err := r.g.WithContext(ctx).
// // 		Model(&db.GDoHiem{}).
// // 		Pluck("so_luong", &items).Error
// // 	return items, err
// // }

// // func (r *doHiemRepoGorm) GetAllMauSacGorm(ctx context.Context) ([]string, error) {
// // 	var items []string
// // 	err := r.g.WithContext(ctx).
// // 		Model(&db.GDoHiem{}).
// // 		Pluck("mau_sac", &items).Error
// // 	return items, err
// // }

// // func (r *doHiemRepoGorm) GetAllSatThuongBonusGorm(ctx context.Context) ([]float64, error) {
// // 	var items []float64
// // 	err := r.g.WithContext(ctx).
// // 		Model(&db.GDoHiem{}).
// // 		Pluck("sat_thuong_bonus", &items).Error
// // 	return items, err
// // }

// // func (r *doHiemRepoGorm) GetAllTocDoDanhBonusGorm(ctx context.Context) ([]float64, error) {
// // 	var items []float64
// // 	err := r.g.WithContext(ctx).
// // 		Model(&db.GDoHiem{}).
// // 		Pluck("toc_do_danh_bonus", &items).Error
// // 	return items, err
// // }

// // // 3.RAW
// // type DoHiemRawInput struct {
// // 	TenDoHiem      string
// // 	MauSac         string
// // 	SoLuong        int
// // 	SatThuongBonus float64
// // 	TocDoDanhBonus float64
// // }

// // type DoHiemRepoRaw interface {
// // 	CreateRaw(ctx context.Context, in DoHiemRawInput) (*DoHiemRaw, error)
// // 	GetIdRaw(ctx context.Context, id int) (*DoHiemRaw, error)
// // 	GetAllRaw(ctx context.Context, q string) ([]DoHiemRaw, error)
// // 	UpdateRaw(ctx context.Context, id int, in DoHiemRawInput) (*DoHiemRaw, error)
// // 	DeleteRaw(ctx context.Context, id int) error

// // 	GetByNameRaw(ctx context.Context, name string) (*DoHiemRaw, error)
// // 	UpdateByNameRaw(ctx context.Context, name string, in DoHiemRawInput) (*DoHiemRaw, error)
// // 	DeleteByNameRaw(ctx context.Context, name string) error

// // 	GetAllTenRaw(ctx context.Context) ([]string, error)
// // 	GetAllSoLuongRaw(ctx context.Context) ([]int, error)
// // 	GetAllMauSacRaw(ctx context.Context) ([]string, error)
// // 	GetAllSatThuongBonusRaw(ctx context.Context) ([]float64, error)
// // 	GetAllTocDoDanhBonusRaw(ctx context.Context) ([]float64, error)
// // }

// // type doHiemRepoRaw struct{ sql *sql.DB }

// // func CreateKhoDoHiemRaw(sqlDB *sql.DB) DoHiemRepoRaw { return &doHiemRepoRaw{sql: sqlDB} }

// // func (r *doHiemRepoRaw) CreateRaw(ctx context.Context, in DoHiemRawInput) (*DoHiemRaw, error) {
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
// // 	if err = tx.QueryRowContext(ctx, `SELECT IFNULL(MAX(ma_do_hiem)+1,1) FROM do_hiem`).Scan(&nextID); err != nil {
// // 		return nil, err
// // 	}

// // 	if _, err = tx.ExecContext(ctx, `
// // 		INSERT INTO do_hiem (ma_do_hiem, ten_do_hiem, so_luong, mau_sac, sat_thuong_bonus, toc_do_danh_bonus)
// // 		VALUES (?, ?, ?, ?, ?, ?)`,
// // 		nextID, in.TenDoHiem, in.SoLuong, in.MauSac, in.SatThuongBonus, in.TocDoDanhBonus,
// // 	); err != nil {
// // 		return nil, err
// // 	}

// // 	if err = tx.Commit(); err != nil {
// // 		return nil, err
// // 	}

// // 	return r.GetIdRaw(ctx, nextID)
// // }

// // func (r *doHiemRepoRaw) GetIdRaw(ctx context.Context, id int) (*DoHiemRaw, error) {
// // 	row := r.sql.QueryRowContext(ctx, `
// // 		SELECT ma_do_hiem, ten_do_hiem, so_luong, mau_sac, sat_thuong_bonus, toc_do_danh_bonus
// // 		FROM do_hiem WHERE ma_do_hiem = ?`, id)

// // 	var v DoHiemRaw
// // 	if err := row.Scan(&v.ID, &v.TenDoHiem, &v.SoLuong, &v.MauSac, &v.SatThuongBonus, &v.TocDoDanhBonus); err != nil {
// // 		return nil, err
// // 	}
// // 	return &v, nil
// // }

// // func (r *doHiemRepoRaw) GetAllRaw(ctx context.Context, q string) ([]DoHiemRaw, error) {
// // 	q = strings.TrimSpace(q)
// // 	var (
// // 		rows *sql.Rows
// // 		err  error
// // 	)
// // 	if q == "" {
// // 		rows, err = r.sql.QueryContext(ctx, `
// // 			SELECT ma_do_hiem, ten_do_hiem, so_luong, mau_sac, sat_thuong_bonus, toc_do_danh_bonus
// // 			FROM do_hiem ORDER BY ma_do_hiem`)
// // 	} else {
// // 		like := "%" + strings.ToLower(q) + "%"
// // 		rows, err = r.sql.QueryContext(ctx, `
// // 			SELECT ma_do_hiem, ten_do_hiem, so_luong, mau_sac, sat_thuong_bonus, toc_do_danh_bonus
// // 			FROM do_hiem
// // 			WHERE LOWER(ten_do_hiem) LIKE ?
// // 			ORDER BY ma_do_hiem`, like)
// // 	}
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer rows.Close()

// // 	var out []DoHiemRaw
// // 	for rows.Next() {
// // 		var v DoHiemRaw
// // 		if err := rows.Scan(&v.ID, &v.TenDoHiem, &v.SoLuong, &v.MauSac, &v.SatThuongBonus, &v.TocDoDanhBonus); err != nil {
// // 			return nil, err
// // 		}
// // 		out = append(out, v)
// // 	}
// // 	return out, rows.Err()
// // }

// // func (r *doHiemRepoRaw) UpdateRaw(ctx context.Context, id int, in DoHiemRawInput) (*DoHiemRaw, error) {
// // 	if _, err := r.sql.ExecContext(ctx, `
// // 		UPDATE do_hiem
// // 		SET ten_do_hiem = ?, so_luong = ?, mau_sac = ?, sat_thuong_bonus = ?, toc_do_danh_bonus = ?
// // 		WHERE ma_do_hiem = ?`,
// // 		in.TenDoHiem, in.SoLuong, in.MauSac, in.SatThuongBonus, in.TocDoDanhBonus, id,
// // 	); err != nil {
// // 		return nil, err
// // 	}
// // 	return r.GetIdRaw(ctx, id)
// // }

// // func (r *doHiemRepoRaw) DeleteRaw(ctx context.Context, id int) error {
// // 	_, err := r.sql.ExecContext(ctx, `DELETE FROM do_hiem WHERE ma_do_hiem = ?`, id)
// // 	return err
// // }

// // func (r *doHiemRepoRaw) GetByNameRaw(ctx context.Context, name string) (*DoHiemRaw, error) {
// // 	row := r.sql.QueryRowContext(ctx, `
// // 		SELECT ma_do_hiem, ten_do_hiem, so_luong, mau_sac, sat_thuong_bonus, toc_do_danh_bonus
// // 		FROM do_hiem
// // 		WHERE ten_do_hiem = ?`, name)
// // 	var v DoHiemRaw
// // 	if err := row.Scan(&v.ID, &v.TenDoHiem, &v.SoLuong, &v.MauSac, &v.SatThuongBonus, &v.TocDoDanhBonus); err != nil {
// // 		return nil, err
// // 	}
// // 	return &v, nil
// // }
// // func (r *doHiemRepoRaw) UpdateByNameRaw(ctx context.Context, name string, in DoHiemRawInput) (*DoHiemRaw, error) {
// // 	var id int
// // 	if err := r.sql.QueryRowContext(ctx, `SELECT ma_do_hiem FROM do_hiem WHERE ten_do_hiem = ?`, name).Scan(&id); err != nil {
// // 		return nil, err
// // 	}

// // 	if _, err := r.sql.ExecContext(ctx, `
// // 		UPDATE do_hiem
// // 		SET ten_do_hiem = ?, so_luong = ?, mau_sac = ?, sat_thuong_bonus = ?, toc_do_danh_bonus = ?
// // 		WHERE ma_do_hiem = ?`,
// // 		in.TenDoHiem, in.SoLuong, in.MauSac, in.SatThuongBonus, in.TocDoDanhBonus, id,
// // 	); err != nil {
// // 		return nil, err
// // 	}

// // 	return r.GetIdRaw(ctx, id)
// // }

// // func (r *doHiemRepoRaw) DeleteByNameRaw(ctx context.Context, name string) error {
// // 	res, err := r.sql.ExecContext(ctx, `DELETE FROM do_hiem WHERE ten_do_hiem = ?`, name)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	if n, _ := res.RowsAffected(); n == 0 {
// // 		return sql.ErrNoRows
// // 	}
// // 	return nil
// // }

// // func (r *doHiemRepoRaw) GetAllTenRaw(ctx context.Context) ([]string, error) {
// // 	rows, err := r.sql.QueryContext(ctx, `SELECT ten_do_hiem FROM do_hiem ORDER BY ten_do_hiem`)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer rows.Close()

// // 	var items []string
// // 	for rows.Next() {
// // 		var s string
// // 		if err := rows.Scan(&s); err != nil {
// // 			return nil, err
// // 		}
// // 		items = append(items, s)
// // 	}
// // 	return items, rows.Err()
// // }

// // func (r *doHiemRepoRaw) GetAllSoLuongRaw(ctx context.Context) ([]int, error) {
// // 	rows, err := r.sql.QueryContext(ctx, `SELECT so_luong FROM do_hiem ORDER BY so_luong`)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer rows.Close()

// // 	var items []int
// // 	for rows.Next() {
// // 		var v int
// // 		if err := rows.Scan(&v); err != nil {
// // 			return nil, err
// // 		}
// // 		items = append(items, v)
// // 	}
// // 	return items, rows.Err()
// // }

// // func (r *doHiemRepoRaw) GetAllMauSacRaw(ctx context.Context) ([]string, error) {
// // 	rows, err := r.sql.QueryContext(ctx, `SELECT mau_sac FROM do_hiem ORDER BY mau_sac`)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer rows.Close()

// // 	var items []string
// // 	for rows.Next() {
// // 		var s string
// // 		if err := rows.Scan(&s); err != nil {
// // 			return nil, err
// // 		}
// // 		items = append(items, s)
// // 	}
// // 	return items, rows.Err()
// // }

// // func (r *doHiemRepoRaw) GetAllSatThuongBonusRaw(ctx context.Context) ([]float64, error) {
// // 	rows, err := r.sql.QueryContext(ctx, `SELECT sat_thuong_bonus FROM do_hiem ORDER BY sat_thuong_bonus`)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer rows.Close()

// // 	var items []float64
// // 	for rows.Next() {
// // 		var v float64
// // 		if err := rows.Scan(&v); err != nil {
// // 			return nil, err
// // 		}
// // 		items = append(items, v)
// // 	}
// // 	return items, rows.Err()
// // }

// // func (r *doHiemRepoRaw) GetAllTocDoDanhBonusRaw(ctx context.Context) ([]float64, error) {
// // 	rows, err := r.sql.QueryContext(ctx, `SELECT toc_do_danh_bonus FROM do_hiem ORDER BY toc_do_danh_bonus`)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer rows.Close()

// // 	var items []float64
// // 	for rows.Next() {
// // 		var v float64
// // 		if err := rows.Scan(&v); err != nil {
// // 			return nil, err
// // 		}
// // 		items = append(items, v)
// // 	}
// // 	return items, rows.Err()
// // }
