package service

// import (
// 	"context"

// 	"game/ent"
// 	"game/internal/db"
// 	"game/internal/repo"
// )

// // 1. ENT
// type DoHiemInputEnt struct {
// 	TenDoHiem      string
// 	MauSac         string
// 	SoLuong        int
// 	SatThuongBonus float64
// 	TocDoDanhBonus float64
// }

// type DoHiemServiceEnt interface {
// 	CreateEnt(ctx context.Context, in DoHiemInputEnt) (*ent.DoHiem, error)

// 	GetByNameEnt(ctx context.Context, name string) (*ent.DoHiem, error)
// 	UpdateByNameEnt(ctx context.Context, name string, in DoHiemInputEnt) (*ent.DoHiem, error)
// 	DeleteByNameEnt(ctx context.Context, name string) error

// 	GetAllEnt(ctx context.Context, q string) ([]*ent.DoHiem, error)
// 	GetAllTen(ctx context.Context) ([]string, error)
// 	GetAllSoLuong(ctx context.Context) ([]int, error)
// 	GetAllMauSac(ctx context.Context) ([]string, error)
// 	GetAllSatThuongBonus(ctx context.Context) ([]float64, error)
// 	GetAllTocDoDanhBonus(ctx context.Context) ([]float64, error)
// }

// type doHiemService struct {
// 	repo repo.DoHiemRepo
// }

// func CreateDoHiemServiceEnt(r repo.DoHiemRepo) DoHiemServiceEnt {
// 	return &doHiemService{repo: r}
// }

// func (s *doHiemService) CreateEnt(ctx context.Context, in DoHiemInputEnt) (*ent.DoHiem, error) {
// 	return s.repo.CreateEnt(ctx, repo.TaoDoHiemInput{
// 		TenDoHiem:      in.TenDoHiem,
// 		MauSac:         in.MauSac,
// 		SoLuong:        in.SoLuong,
// 		SatThuongBonus: in.SatThuongBonus,
// 		TocDoDanhBonus: in.TocDoDanhBonus,
// 	})
// }

// func (s *doHiemService) GetByNameEnt(ctx context.Context, name string) (*ent.DoHiem, error) {
// 	return s.repo.GetByNameEnt(ctx, name)
// }

// func (s *doHiemService) UpdateByNameEnt(ctx context.Context, name string, in DoHiemInputEnt) (*ent.DoHiem, error) {
// 	return s.repo.UpdateByNameEnt(ctx, name, func(up *ent.DoHiemUpdateOne) *ent.DoHiemUpdateOne {
// 		return up.
// 			SetTenDoHiem(in.TenDoHiem).
// 			SetMauSac(in.MauSac).
// 			SetSoLuong(in.SoLuong).
// 			SetSatThuongBonus(in.SatThuongBonus).
// 			SetTocDoDanhBonus(in.TocDoDanhBonus)
// 	})
// }

// func (s *doHiemService) DeleteByNameEnt(ctx context.Context, name string) error {
// 	return s.repo.DeleteByNameEnt(ctx, name)
// }

// func (s *doHiemService) GetAllEnt(ctx context.Context, q string) ([]*ent.DoHiem, error) {
// 	return s.repo.GetAllEnt(ctx, q)
// }
// func (s *doHiemService) GetAllTen(ctx context.Context) ([]string, error) {
// 	return s.repo.GetAllTen(ctx)
// }
// func (s *doHiemService) GetAllSoLuong(ctx context.Context) ([]int, error) {
// 	return s.repo.GetAllSoLuong(ctx)
// }
// func (s *doHiemService) GetAllMauSac(ctx context.Context) ([]string, error) {
// 	return s.repo.GetAllMauSac(ctx)
// }
// func (s *doHiemService) GetAllSatThuongBonus(ctx context.Context) ([]float64, error) {
// 	return s.repo.GetAllSatThuongBonus(ctx)
// }
// func (s *doHiemService) GetAllTocDoDanhBonus(ctx context.Context) ([]float64, error) {
// 	return s.repo.GetAllTocDoDanhBonus(ctx)
// }

// // 2. GORM
// type DoHiemInputGorm struct {
// 	TenDoHiem      string
// 	MauSac         string
// 	SoLuong        int
// 	SatThuongBonus float64
// 	TocDoDanhBonus float64
// }

// type DoHiemServiceGorm interface {
// 	CreateGorm(ctx context.Context, in DoHiemInputGorm) (*db.GDoHiem, error)
// 	GetIdGorm(ctx context.Context, id int) (*db.GDoHiem, error)
// 	GetAllGorm(ctx context.Context, q string) ([]db.GDoHiem, error)
// 	UpdateGorm(ctx context.Context, id int, in DoHiemInputGorm) (*db.GDoHiem, error)
// 	DeleteGorm(ctx context.Context, id int) error

// 	GetByNameGorm(ctx context.Context, name string) (*db.GDoHiem, error)
// 	UpdateByNameGorm(ctx context.Context, name string, in DoHiemInputGorm) (*db.GDoHiem, error)
// 	DeleteByNameGorm(ctx context.Context, name string) error

// 	GetAllTenGorm(ctx context.Context) ([]string, error)
// 	GetAllSoLuongGorm(ctx context.Context) ([]int, error)
// 	GetAllMauSacGorm(ctx context.Context) ([]string, error)
// 	GetAllSatThuongBonusGorm(ctx context.Context) ([]float64, error)
// 	GetAllTocDoDanhBonusGorm(ctx context.Context) ([]float64, error)
// }

// type doHiemServiceGorm struct {
// 	repo repo.DoHiemRepoGorm
// }

// func CreateDoHiemServiceGorm(r repo.DoHiemRepoGorm) DoHiemServiceGorm {
// 	return &doHiemServiceGorm{repo: r}
// }

// func (s *doHiemServiceGorm) CreateGorm(ctx context.Context, in DoHiemInputGorm) (*db.GDoHiem, error) {
// 	return s.repo.CreateGorm(ctx, repo.DoHiemGormInput{
// 		TenDoHiem:      in.TenDoHiem,
// 		MauSac:         in.MauSac,
// 		SoLuong:        in.SoLuong,
// 		SatThuongBonus: in.SatThuongBonus,
// 		TocDoDanhBonus: in.TocDoDanhBonus,
// 	})
// }

// func (s *doHiemServiceGorm) GetIdGorm(ctx context.Context, id int) (*db.GDoHiem, error) {
// 	return s.repo.GetIdGorm(ctx, id)
// }

// func (s *doHiemServiceGorm) GetAllGorm(ctx context.Context, q string) ([]db.GDoHiem, error) {
// 	return s.repo.GetAllGorm(ctx, q)
// }

// func (s *doHiemServiceGorm) UpdateGorm(ctx context.Context, id int, in DoHiemInputGorm) (*db.GDoHiem, error) {
// 	return s.repo.UpdateGorm(ctx, id, repo.DoHiemGormInput{
// 		TenDoHiem:      in.TenDoHiem,
// 		MauSac:         in.MauSac,
// 		SoLuong:        in.SoLuong,
// 		SatThuongBonus: in.SatThuongBonus,
// 		TocDoDanhBonus: in.TocDoDanhBonus,
// 	})
// }

// func (s *doHiemServiceGorm) DeleteGorm(ctx context.Context, id int) error {
// 	return s.repo.DeleteGorm(ctx, id)
// }

// func (s *doHiemServiceGorm) GetByNameGorm(ctx context.Context, name string) (*db.GDoHiem, error) {
// 	return s.repo.GetByNameGorm(ctx, name)
// }

// func (s *doHiemServiceGorm) UpdateByNameGorm(ctx context.Context, name string, in DoHiemInputGorm) (*db.GDoHiem, error) {
// 	return s.repo.UpdateByNameGorm(ctx, name, repo.DoHiemGormInput{
// 		TenDoHiem:      in.TenDoHiem,
// 		MauSac:         in.MauSac,
// 		SoLuong:        in.SoLuong,
// 		SatThuongBonus: in.SatThuongBonus,
// 		TocDoDanhBonus: in.TocDoDanhBonus,
// 	})
// }

// func (s *doHiemServiceGorm) DeleteByNameGorm(ctx context.Context, name string) error {
// 	return s.repo.DeleteByNameGorm(ctx, name)
// }

// func (s *doHiemServiceGorm) GetAllTenGorm(ctx context.Context) ([]string, error) {
// 	return s.repo.GetAllTenGorm(ctx)
// }

// func (s *doHiemServiceGorm) GetAllSoLuongGorm(ctx context.Context) ([]int, error) {
// 	return s.repo.GetAllSoLuongGorm(ctx)
// }

// func (s *doHiemServiceGorm) GetAllMauSacGorm(ctx context.Context) ([]string, error) {
// 	return s.repo.GetAllMauSacGorm(ctx)
// }

// func (s *doHiemServiceGorm) GetAllSatThuongBonusGorm(ctx context.Context) ([]float64, error) {
// 	return s.repo.GetAllSatThuongBonusGorm(ctx)
// }

// func (s *doHiemServiceGorm) GetAllTocDoDanhBonusGorm(ctx context.Context) ([]float64, error) {
// 	return s.repo.GetAllTocDoDanhBonusGorm(ctx)
// }

// // 3. RAW
// type DoHiemInputRaw struct {
// 	TenDoHiem      string
// 	MauSac         string
// 	SoLuong        int
// 	SatThuongBonus float64
// 	TocDoDanhBonus float64
// }

// type DoHiemServiceRaw interface {
// 	CreateRaw(ctx context.Context, in DoHiemInputRaw) (*repo.DoHiemRaw, error)
// 	GetIdRaw(ctx context.Context, id int) (*repo.DoHiemRaw, error)
// 	GetAllRaw(ctx context.Context, q string) ([]repo.DoHiemRaw, error)
// 	UpdateRaw(ctx context.Context, id int, in DoHiemInputRaw) (*repo.DoHiemRaw, error)
// 	DeleteRaw(ctx context.Context, id int) error

// 	GetByNameRaw(ctx context.Context, name string) (*repo.DoHiemRaw, error)
// 	UpdateByNameRaw(ctx context.Context, name string, in DoHiemInputRaw) (*repo.DoHiemRaw, error)
// 	DeleteByNameRaw(ctx context.Context, name string) error

// 	GetAllTenRaw(ctx context.Context) ([]string, error)
// 	GetAllSoLuongRaw(ctx context.Context) ([]int, error)
// 	GetAllMauSacRaw(ctx context.Context) ([]string, error)
// 	GetAllSatThuongBonusRaw(ctx context.Context) ([]float64, error)
// 	GetAllTocDoDanhBonusRaw(ctx context.Context) ([]float64, error)
// }

// type doHiemServiceRaw struct{ repo repo.DoHiemRepoRaw }

// func CreateDoHiemServiceRaw(r repo.DoHiemRepoRaw) DoHiemServiceRaw {
// 	return &doHiemServiceRaw{repo: r}
// }

// func (s *doHiemServiceRaw) CreateRaw(ctx context.Context, in DoHiemInputRaw) (*repo.DoHiemRaw, error) {
// 	return s.repo.CreateRaw(ctx, repo.DoHiemRawInput{
// 		TenDoHiem:      in.TenDoHiem,
// 		MauSac:         in.MauSac,
// 		SoLuong:        in.SoLuong,
// 		SatThuongBonus: in.SatThuongBonus,
// 		TocDoDanhBonus: in.TocDoDanhBonus,
// 	})
// }

// func (s *doHiemServiceRaw) GetIdRaw(ctx context.Context, id int) (*repo.DoHiemRaw, error) {
// 	return s.repo.GetIdRaw(ctx, id)
// }

// func (s *doHiemServiceRaw) GetAllRaw(ctx context.Context, q string) ([]repo.DoHiemRaw, error) {
// 	return s.repo.GetAllRaw(ctx, q)
// }

// func (s *doHiemServiceRaw) UpdateRaw(ctx context.Context, id int, in DoHiemInputRaw) (*repo.DoHiemRaw, error) {
// 	return s.repo.UpdateRaw(ctx, id, repo.DoHiemRawInput{
// 		TenDoHiem:      in.TenDoHiem,
// 		MauSac:         in.MauSac,
// 		SoLuong:        in.SoLuong,
// 		SatThuongBonus: in.SatThuongBonus,
// 		TocDoDanhBonus: in.TocDoDanhBonus,
// 	})
// }

// func (s *doHiemServiceRaw) DeleteRaw(ctx context.Context, id int) error {
// 	return s.repo.DeleteRaw(ctx, id)
// }

// func (s *doHiemServiceRaw) GetByNameRaw(ctx context.Context, name string) (*repo.DoHiemRaw, error) {
// 	return s.repo.GetByNameRaw(ctx, name)
// }

// func (s *doHiemServiceRaw) UpdateByNameRaw(ctx context.Context, name string, in DoHiemInputRaw) (*repo.DoHiemRaw, error) {
// 	return s.repo.UpdateByNameRaw(ctx, name, repo.DoHiemRawInput{
// 		TenDoHiem:      in.TenDoHiem,
// 		MauSac:         in.MauSac,
// 		SoLuong:        in.SoLuong,
// 		SatThuongBonus: in.SatThuongBonus,
// 		TocDoDanhBonus: in.TocDoDanhBonus,
// 	})
// }

// func (s *doHiemServiceRaw) DeleteByNameRaw(ctx context.Context, name string) error {
// 	return s.repo.DeleteByNameRaw(ctx, name)
// }

// func (s *doHiemServiceRaw) GetAllTenRaw(ctx context.Context) ([]string, error) {
// 	return s.repo.GetAllTenRaw(ctx)
// }
// func (s *doHiemServiceRaw) GetAllSoLuongRaw(ctx context.Context) ([]int, error) {
// 	return s.repo.GetAllSoLuongRaw(ctx)
// }
// func (s *doHiemServiceRaw) GetAllMauSacRaw(ctx context.Context) ([]string, error) {
// 	return s.repo.GetAllMauSacRaw(ctx)
// }
// func (s *doHiemServiceRaw) GetAllSatThuongBonusRaw(ctx context.Context) ([]float64, error) {
// 	return s.repo.GetAllSatThuongBonusRaw(ctx)
// }
// func (s *doHiemServiceRaw) GetAllTocDoDanhBonusRaw(ctx context.Context) ([]float64, error) {
// 	return s.repo.GetAllTocDoDanhBonusRaw(ctx)
// }
