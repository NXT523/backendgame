package repo

import (
	"context"
	"database/sql"
	"strings"

	"game/ent"
	"game/ent/he"
	"game/internal/db"

	"gorm.io/gorm"
)

// 1. ENT
type CreateHeInput struct {
	TenHe string
	MoTa  *string
}

type HeSearchRequest struct {
	ID   *int   `json:"id"`
	Name string `json:"name"`
	Q    string `json:"q"`
}

type HeRepo interface {
	CreateEnt(ctx context.Context, in CreateHeInput) (*ent.He, error)
	GetIdEnt(ctx context.Context, id int) (*ent.He, error)
	GetAllEnt(ctx context.Context, q string) ([]*ent.He, error)
	UpdateEnt(ctx context.Context, id int, mutate func(*ent.HeUpdateOne) *ent.HeUpdateOne) (*ent.He, error)
	DeleteEnt(ctx context.Context, id int) error

	GetByNameEnt(ctx context.Context, name string) (*ent.He, error)
	UpdateByNameEnt(ctx context.Context, name string, mutate func(*ent.HeUpdateOne) *ent.HeUpdateOne) (*ent.He, error)
	DeleteByNameEnt(ctx context.Context, name string) error

	GetAllNamesEnt(ctx context.Context, q string) ([]string, error)
}

type heRepo struct {
	c *ent.Client
}

func CreateKhoHeEnt(c *ent.Client) HeRepo {
	return &heRepo{c: c}
}

func (r *heRepo) CreateEnt(ctx context.Context, in CreateHeInput) (*ent.He, error) {
	return r.c.He.Create().
		SetTenHe(in.TenHe).
		SetNillableMoTa(in.MoTa).
		Save(ctx)
}

func (r *heRepo) GetIdEnt(ctx context.Context, id int) (*ent.He, error) {
	return r.c.He.Query().
		Where(he.IDEQ(id)).
		Only(ctx)
}

func (r *heRepo) GetAllEnt(ctx context.Context, q string) ([]*ent.He, error) {
	qb := r.c.He.Query()
	if q != "" {
		qb = qb.Where(he.TenHeContainsFold(q))
	}
	return qb.All(ctx)
}

func (r *heRepo) UpdateEnt(ctx context.Context, id int, mutate func(*ent.HeUpdateOne) *ent.HeUpdateOne) (*ent.He, error) {
	up := r.c.He.UpdateOneID(id)
	return mutate(up).Save(ctx)
}

func (r *heRepo) DeleteEnt(ctx context.Context, id int) error {
	return r.c.He.DeleteOneID(id).Exec(ctx)
}

func (r *heRepo) GetByNameEnt(ctx context.Context, name string) (*ent.He, error) {
	return r.c.He.Query().Where(he.TenHeEQ(name)).Only(ctx)
}

func (r *heRepo) UpdateByNameEnt(ctx context.Context, name string, mutate func(*ent.HeUpdateOne) *ent.HeUpdateOne) (*ent.He, error) {
	cur, err := r.c.He.Query().Where(he.TenHeEQ(name)).Only(ctx)
	if err != nil {
		return nil, err
	}
	up := r.c.He.UpdateOneID(cur.ID)
	return mutate(up).Save(ctx)
}

func (r *heRepo) DeleteByNameEnt(ctx context.Context, name string) error {
	cur, err := r.c.He.Query().Where(he.TenHeEQ(name)).Only(ctx)
	if err != nil {
		return err
	}
	return r.c.He.DeleteOneID(cur.ID).Exec(ctx)
}
func (r *heRepo) GetAllNamesEnt(ctx context.Context, q string) ([]string, error) {
	qb := r.c.He.Query()
	if q != "" {
		qb = qb.Where(he.TenHeContainsFold(q))
	}
	return qb.Select(he.FieldTenHe).Strings(ctx)
}



// 2. GORM
type HeGormInput struct {
	TenHe string
	MoTa  *string
}

type HeRepoGorm interface {
	CreateGorm(ctx context.Context, in HeGormInput) (*db.GHe, error)
	GetIdGorm(ctx context.Context, id int) (*db.GHe, error)
	GetAllGorm(ctx context.Context, q string) ([]db.GHe, error)
	UpdateGorm(ctx context.Context, id int, in HeGormInput) (*db.GHe, error)
	DeleteGorm(ctx context.Context, id int) error

	GetByNameGorm(ctx context.Context, name string) (*db.GHe, error)
	UpdateByNameGorm(ctx context.Context, name string, in HeGormInput) (*db.GHe, error)
	DeleteByNameGorm(ctx context.Context, name string) error

	GetAllNamesGorm(ctx context.Context, q string) ([]string, error)
}

type heRepoGorm struct{ g *gorm.DB }

func CreateKhoHeGorm(g *gorm.DB) HeRepoGorm { return &heRepoGorm{g: g} }

func (r *heRepoGorm) CreateGorm(ctx context.Context, in HeGormInput) (*db.GHe, error) {
	row := db.GHe{
		TenHe: in.TenHe,
		MoTa:  in.MoTa,
	}
	if err := r.g.WithContext(ctx).Create(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *heRepoGorm) GetIdGorm(ctx context.Context, id int) (*db.GHe, error) {
	var row db.GHe
	if err := r.g.WithContext(ctx).First(&row, "ma_he = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *heRepoGorm) GetAllGorm(ctx context.Context, q string) ([]db.GHe, error) {
	var rows []db.GHe
	tx := r.g.WithContext(ctx).Model(&db.GHe{})
	if q = strings.TrimSpace(q); q != "" {
		tx = tx.Where("LOWER(ten_he) LIKE ?", "%"+strings.ToLower(q)+"%")
	}
	if err := tx.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *heRepoGorm) UpdateGorm(ctx context.Context, id int, in HeGormInput) (*db.GHe, error) {
	var row db.GHe
	if err := r.g.WithContext(ctx).First(&row, "ma_he = ?", id).Error; err != nil {
		return nil, err
	}
	row.TenHe = in.TenHe
	row.MoTa = in.MoTa
	if err := r.g.WithContext(ctx).Save(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *heRepoGorm) DeleteGorm(ctx context.Context, id int) error {
	return r.g.WithContext(ctx).Delete(&db.GHe{}, "ma_he = ?", id).Error
}

func (r *heRepoGorm) GetByNameGorm(ctx context.Context, name string) (*db.GHe, error) {
	var row db.GHe
	if err := r.g.WithContext(ctx).Where("ten_he = ?", name).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *heRepoGorm) UpdateByNameGorm(ctx context.Context, name string, in HeGormInput) (*db.GHe, error) {
	var row db.GHe
	if err := r.g.WithContext(ctx).Where("ten_he = ?", name).First(&row).Error; err != nil {
		return nil, err
	}
	row.TenHe = in.TenHe
	row.MoTa = in.MoTa
	if err := r.g.WithContext(ctx).Save(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *heRepoGorm) DeleteByNameGorm(ctx context.Context, name string) error {
	res := r.g.WithContext(ctx).Where("ten_he = ?", name).Delete(&db.GHe{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *heRepoGorm) GetAllNamesGorm(ctx context.Context, q string) ([]string, error) {
	var names []string
	tx := r.g.WithContext(ctx).Model(&db.GHe{})
	if q = strings.TrimSpace(q); q != "" {
		tx = tx.Where("LOWER(ten_he) LIKE ?", "%"+strings.ToLower(q)+"%")
	}
	if err := tx.Pluck("ten_he", &names).Error; err != nil {
		return nil, err
	}
	return names, nil
}

// 3. RAW
type HeRaw struct {
	ID    int
	TenHe string
	MoTa  *string
}

type HeRawInput struct {
	TenHe string
	MoTa  *string
}

type HeRepoRaw interface {
	CreateRaw(ctx context.Context, in HeRawInput) (*HeRaw, error)
	GetIdRaw(ctx context.Context, id int) (*HeRaw, error)
	GetAllRaw(ctx context.Context, q string) ([]HeRaw, error)
	UpdateRaw(ctx context.Context, id int, in HeRawInput) (*HeRaw, error)
	DeleteRaw(ctx context.Context, id int) error

	GetByNameRaw(ctx context.Context, name string) (*HeRaw, error)
	UpdateByNameRaw(ctx context.Context, name string, in HeRawInput) (*HeRaw, error)
	DeleteByNameRaw(ctx context.Context, name string) error

	GetAllNamesRaw(ctx context.Context, q string) ([]string, error)
}

type heRepoRaw struct{ sql *sql.DB }

func CreateKhoHeRaw(sqlDB *sql.DB) HeRepoRaw { return &heRepoRaw{sql: sqlDB} }

func (r *heRepoRaw) CreateRaw(ctx context.Context, in HeRawInput) (*HeRaw, error) {
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
	if err = tx.QueryRowContext(ctx, `SELECT IFNULL(MAX(ma_he)+1,1) FROM he`).Scan(&nextID); err != nil {
		return nil, err
	}
	if _, err = tx.ExecContext(ctx, `
		INSERT INTO he (ma_he, ten_he, mo_ta)
		VALUES (?, ?, ?)`,
		nextID, in.TenHe, in.MoTa,
	); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return r.GetIdRaw(ctx, nextID)
}

func (r *heRepoRaw) GetIdRaw(ctx context.Context, id int) (*HeRaw, error) {
	row := r.sql.QueryRowContext(ctx, `
		SELECT ma_he, ten_he, mo_ta
		FROM he
		WHERE ma_he = ?`, id)

	var v HeRaw
	if err := row.Scan(&v.ID, &v.TenHe, &v.MoTa); err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *heRepoRaw) GetAllRaw(ctx context.Context, q string) ([]HeRaw, error) {
	q = strings.TrimSpace(q)
	var (
		rows *sql.Rows
		err  error
	)
	if q == "" {
		rows, err = r.sql.QueryContext(ctx, `
			SELECT ma_he, ten_he, mo_ta
			FROM he
			ORDER BY ma_he`)
	} else {
		like := "%" + strings.ToLower(q) + "%"
		rows, err = r.sql.QueryContext(ctx, `
			SELECT ma_he, ten_he, mo_ta
			FROM he
			WHERE LOWER(ten_he) LIKE ?
			ORDER BY ma_he`, like)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []HeRaw
	for rows.Next() {
		var v HeRaw
		if err := rows.Scan(&v.ID, &v.TenHe, &v.MoTa); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func (r *heRepoRaw) UpdateRaw(ctx context.Context, id int, in HeRawInput) (*HeRaw, error) {
	if _, err := r.sql.ExecContext(ctx, `
		UPDATE he
		SET ten_he = ?, mo_ta = ?
		WHERE ma_he = ?`,
		in.TenHe, in.MoTa, id,
	); err != nil {
		return nil, err
	}
	return r.GetIdRaw(ctx, id)
}

func (r *heRepoRaw) DeleteRaw(ctx context.Context, id int) error {
	_, err := r.sql.ExecContext(ctx, `DELETE FROM he WHERE ma_he = ?`, id)
	return err
}

func (r *heRepoRaw) GetByNameRaw(ctx context.Context, name string) (*HeRaw, error) {
	row := r.sql.QueryRowContext(ctx, `
        SELECT ma_he, ten_he, mo_ta
        FROM he
        WHERE ten_he = ?`, name)

	var v HeRaw
	if err := row.Scan(&v.ID, &v.TenHe, &v.MoTa); err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *heRepoRaw) UpdateByNameRaw(ctx context.Context, name string, in HeRawInput) (*HeRaw, error) {
	var id int
	if err := r.sql.QueryRowContext(ctx, `SELECT ma_he FROM he WHERE ten_he = ?`, name).Scan(&id); err != nil {
		return nil, err
	}
	if _, err := r.sql.ExecContext(ctx, `
        UPDATE he
        SET ten_he = ?, mo_ta = ?
        WHERE ma_he = ?`,
		in.TenHe, in.MoTa, id,
	); err != nil {
		return nil, err
	}
	return r.GetIdRaw(ctx, id)
}

func (r *heRepoRaw) DeleteByNameRaw(ctx context.Context, name string) error {
	res, err := r.sql.ExecContext(ctx, `DELETE FROM he WHERE ten_he = ?`, name)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *heRepoRaw) GetAllNamesRaw(ctx context.Context, q string) ([]string, error) {
	q = strings.TrimSpace(q)
	var (
		rows *sql.Rows
		err  error
	)
	if q == "" {
		rows, err = r.sql.QueryContext(ctx, `SELECT ten_he FROM he ORDER BY ten_he`)
	} else {
		like := "%" + strings.ToLower(q) + "%"
		rows, err = r.sql.QueryContext(ctx, `SELECT ten_he FROM he WHERE LOWER(ten_he) LIKE ? ORDER BY ten_he`, like)
	}
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
