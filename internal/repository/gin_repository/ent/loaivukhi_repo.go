package gin_repository

// import (
// 	"context"
// 	"database/sql"
// 	"strings"

// 	"game/ent"
// 	"game/ent/loaivukhi"
// 	"game/internal/db"

// 	"gorm.io/gorm"
// )

// // 1. ENT
// type CreateLoaiVKInput struct {
// 	TenLoai string
// 	MoTa    *string
// }


// type loaiVuKhiRepo struct{ c *ent.Client }

// func CreateKhoLoaiVuKhiEnt(c *ent.Client) LoaiVuKhiRepo {
// 	return &loaiVuKhiRepo{c: c}
// }

// func (r *loaiVuKhiRepo) CreateEnt(ctx context.Context, in CreateLoaiVKInput) (*ent.LoaiVuKhi, error) {
// 	return r.c.LoaiVuKhi.Create().
// 		SetTenLoaiVuKhi(in.TenLoai).
// 		SetNillableMoTa(in.MoTa).
// 		Save(ctx)
// }

// func (r *loaiVuKhiRepo) GetIdEnt(ctx context.Context, id int) (*ent.LoaiVuKhi, error) {
// 	return r.c.LoaiVuKhi.Query().
// 		Where(loaivukhi.IDEQ(id)).
// 		Only(ctx)
// }

// func (r *loaiVuKhiRepo) GetAllEnt(ctx context.Context, q string) ([]*ent.LoaiVuKhi, error) {
// 	qb := r.c.LoaiVuKhi.Query()
// 	if q != "" {
// 		qb = qb.Where(loaivukhi.TenLoaiVuKhiContainsFold(q))
// 	}
// 	return qb.All(ctx)
// }

// func (r *loaiVuKhiRepo) UpdateEnt(ctx context.Context, id int, mutate func(*ent.LoaiVuKhiUpdateOne) *ent.LoaiVuKhiUpdateOne) (*ent.LoaiVuKhi, error) {
// 	up := r.c.LoaiVuKhi.UpdateOneID(id)
// 	return mutate(up).Save(ctx)
// }

// func (r *loaiVuKhiRepo) DeleteEnt(ctx context.Context, id int) error {
// 	return r.c.LoaiVuKhi.DeleteOneID(id).Exec(ctx)
// }

// func (r *loaiVuKhiRepo) GetTenLoaiEnt(ctx context.Context, q string) ([]string, error) {
// 	qb := r.c.LoaiVuKhi.Query()
// 	if q != "" {
// 		qb = qb.Where(loaivukhi.TenLoaiVuKhiContainsFold(q))
// 	}
// 	return qb.
// 		Select(loaivukhi.FieldTenLoaiVuKhi).
// 		Strings(ctx)
// }

// // 2. GORM
// type LoaiVuKhiGormInput struct {
// 	TenLoai string
// 	MoTa    *string
// }

// type LoaiVuKhiRepoGorm interface {
// 	CreateGorm(ctx context.Context, in LoaiVuKhiGormInput) (*db.GLoaiVuKhi, error)
// 	GetIdGorm(ctx context.Context, id int) (*db.GLoaiVuKhi, error)
// 	GetAllGorm(ctx context.Context, q string) ([]db.GLoaiVuKhi, error)
// 	UpdateGorm(ctx context.Context, id int, in LoaiVuKhiGormInput) (*db.GLoaiVuKhi, error)
// 	DeleteGorm(ctx context.Context, id int) error

// 	GetTenLoaiGorm(ctx context.Context, q string) ([]string, error)
// }

// type loaiVuKhiRepoGorm struct{ g *gorm.DB }

// func CreateKhoLoaiVuKhiGorm(g *gorm.DB) LoaiVuKhiRepoGorm { return &loaiVuKhiRepoGorm{g: g} }

// func (r *loaiVuKhiRepoGorm) CreateGorm(ctx context.Context, in LoaiVuKhiGormInput) (*db.GLoaiVuKhi, error) {
// 	row := db.GLoaiVuKhi{
// 		TenLoaiVuKhi: in.TenLoai,
// 		MoTa:    in.MoTa,
// 	}
// 	if err := r.g.WithContext(ctx).Create(&row).Error; err != nil {
// 		return nil, err
// 	}
// 	return &row, nil
// }

// func (r *loaiVuKhiRepoGorm) GetIdGorm(ctx context.Context, id int) (*db.GLoaiVuKhi, error) {
// 	var row db.GLoaiVuKhi
// 	if err := r.g.WithContext(ctx).First(&row, "ma_loai_vu_khi = ?", id).Error; err != nil {
// 		return nil, err
// 	}
// 	return &row, nil
// }

// func (r *loaiVuKhiRepoGorm) GetAllGorm(ctx context.Context, q string) ([]db.GLoaiVuKhi, error) {
// 	var rows []db.GLoaiVuKhi
// 	tx := r.g.WithContext(ctx).Model(&db.GLoaiVuKhi{})
// 	if q = strings.TrimSpace(q); q != "" {
// 		tx = tx.Where("LOWER(ten_loai) LIKE ?", "%"+strings.ToLower(q)+"%")
// 	}
// 	if err := tx.Find(&rows).Error; err != nil {
// 		return nil, err
// 	}
// 	return rows, nil
// }

// func (r *loaiVuKhiRepoGorm) UpdateGorm(ctx context.Context, id int, in LoaiVuKhiGormInput) (*db.GLoaiVuKhi, error) {
// 	var row db.GLoaiVuKhi
// 	if err := r.g.WithContext(ctx).First(&row, "ma_loai_vu_khi = ?", id).Error; err != nil {
// 		return nil, err
// 	}
// 	row.TenLoaiVuKhi = in.TenLoai
// 	row.MoTa = in.MoTa
// 	if err := r.g.WithContext(ctx).Save(&row).Error; err != nil {
// 		return nil, err
// 	}
// 	return &row, nil
// }

// func (r *loaiVuKhiRepoGorm) DeleteGorm(ctx context.Context, id int) error {
// 	return r.g.WithContext(ctx).Delete(&db.GLoaiVuKhi{}, "ma_loai_vu_khi = ?", id).Error
// }

// func (r *loaiVuKhiRepoGorm) GetTenLoaiGorm(ctx context.Context, q string) ([]string, error) {
// 	var names []string
// 	tx := r.g.WithContext(ctx).Model(&db.GLoaiVuKhi{})
// 	if q = strings.TrimSpace(q); q != "" {
// 		tx = tx.Where("LOWER(ten_loai) LIKE ?", "%"+strings.ToLower(q)+"%")
// 	}
// 	if err := tx.Pluck("ten_loai", &names).Error; err != nil {
// 		return nil, err
// 	}
// 	return names, nil
// }

// // 3. RAW
// type LoaiVuKhiRaw struct {
// 	ID      int
// 	TenLoai string
// 	MoTa    *string
// }

// type LoaiVuKhiRawInput struct {
// 	TenLoai string
// 	MoTa    *string
// }

// type LoaiVuKhiRepoRaw interface {
// 	CreateRaw(ctx context.Context, in LoaiVuKhiRawInput) (*LoaiVuKhiRaw, error)
// 	GetIdRaw(ctx context.Context, id int) (*LoaiVuKhiRaw, error)
// 	GetAllRaw(ctx context.Context, q string) ([]LoaiVuKhiRaw, error)
// 	UpdateRaw(ctx context.Context, id int, in LoaiVuKhiRawInput) (*LoaiVuKhiRaw, error)
// 	DeleteRaw(ctx context.Context, id int) error

// 	GetTenLoaiRaw(ctx context.Context, q string) ([]string, error)
// }

// type loaiVuKhiRepoRaw struct{ sql *sql.DB }

// func CreateKhoLoaiVuKhiRaw(sqlDB *sql.DB) LoaiVuKhiRepoRaw { return &loaiVuKhiRepoRaw{sql: sqlDB} }

// func (r *loaiVuKhiRepoRaw) CreateRaw(ctx context.Context, in LoaiVuKhiRawInput) (*LoaiVuKhiRaw, error) {
// 	tx, err := r.sql.BeginTx(ctx, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer func() {
// 		if err != nil {
// 			_ = tx.Rollback()
// 		}
// 	}()

// 	var nextID int
// 	if err = tx.QueryRowContext(ctx, `SELECT IFNULL(MAX(ma_loai_vu_khi)+1,1) FROM loai_vu_khi`).Scan(&nextID); err != nil {
// 		return nil, err
// 	}
// 	if _, err = tx.ExecContext(ctx, `
// 		INSERT INTO loai_vu_khi (ma_loai_vu_khi, ten_loai, mo_ta)
// 		VALUES (?, ?, ?)`,
// 		nextID, in.TenLoai, in.MoTa,
// 	); err != nil {
// 		return nil, err
// 	}
// 	if err = tx.Commit(); err != nil {
// 		return nil, err
// 	}
// 	return r.GetIdRaw(ctx, nextID)
// }

// func (r *loaiVuKhiRepoRaw) GetIdRaw(ctx context.Context, id int) (*LoaiVuKhiRaw, error) {
// 	row := r.sql.QueryRowContext(ctx, `
// 		SELECT ma_loai_vu_khi, ten_loai, mo_ta
// 		FROM loai_vu_khi
// 		WHERE ma_loai_vu_khi = ?`, id)

// 	var v LoaiVuKhiRaw
// 	if err := row.Scan(&v.ID, &v.TenLoai, &v.MoTa); err != nil {
// 		return nil, err
// 	}
// 	return &v, nil
// }

// func (r *loaiVuKhiRepoRaw) GetAllRaw(ctx context.Context, q string) ([]LoaiVuKhiRaw, error) {
// 	q = strings.TrimSpace(q)
// 	var (
// 		rows *sql.Rows
// 		err  error
// 	)
// 	if q == "" {
// 		rows, err = r.sql.QueryContext(ctx, `
// 			SELECT ma_loai_vu_khi, ten_loai, mo_ta
// 			FROM loai_vu_khi
// 			ORDER BY ma_loai_vu_khi`)
// 	} else {
// 		like := "%" + strings.ToLower(q) + "%"
// 		rows, err = r.sql.QueryContext(ctx, `
// 			SELECT ma_loai_vu_khi, ten_loai, mo_ta
// 			FROM loai_vu_khi
// 			WHERE LOWER(ten_loai) LIKE ?
// 			ORDER BY ma_loai_vu_khi`, like)
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var out []LoaiVuKhiRaw
// 	for rows.Next() {
// 		var v LoaiVuKhiRaw
// 		if err := rows.Scan(&v.ID, &v.TenLoai, &v.MoTa); err != nil {
// 			return nil, err
// 		}
// 		out = append(out, v)
// 	}
// 	return out, rows.Err()
// }

// func (r *loaiVuKhiRepoRaw) UpdateRaw(ctx context.Context, id int, in LoaiVuKhiRawInput) (*LoaiVuKhiRaw, error) {
// 	if _, err := r.sql.ExecContext(ctx, `
// 		UPDATE loai_vu_khi
// 		SET ten_loai = ?, mo_ta = ?
// 		WHERE ma_loai_vu_khi = ?`,
// 		in.TenLoai, in.MoTa, id,
// 	); err != nil {
// 		return nil, err
// 	}
// 	return r.GetIdRaw(ctx, id)
// }

// func (r *loaiVuKhiRepoRaw) DeleteRaw(ctx context.Context, id int) error {
// 	_, err := r.sql.ExecContext(ctx, `DELETE FROM loai_vu_khi WHERE ma_loai_vu_khi = ?`, id)
// 	return err
// }

// func (r *loaiVuKhiRepoRaw) GetTenLoaiRaw(ctx context.Context, q string) ([]string, error) {
// 	q = strings.TrimSpace(q)
// 	var (
// 		rows *sql.Rows
// 		err  error
// 	)
// 	if q == "" {
// 		rows, err = r.sql.QueryContext(ctx, `SELECT ten_loai FROM loai_vu_khi ORDER BY ten_loai`)
// 	} else {
// 		like := "%" + strings.ToLower(q) + "%"
// 		rows, err = r.sql.QueryContext(ctx, `SELECT ten_loai FROM loai_vu_khi WHERE LOWER(ten_loai) LIKE ? ORDER BY ten_loai`, like)
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var names []string
// 	for rows.Next() {
// 		var s string
// 		if err := rows.Scan(&s); err != nil {
// 			return nil, err
// 		}
// 		names = append(names, s)
// 	}
// 	return names, rows.Err()
// }
