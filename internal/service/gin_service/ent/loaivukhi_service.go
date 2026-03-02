package service

// import (
// 	"context"

// 	"game/ent"
// 	"game/internal/db"
// 	"game/internal/repo"
// )

// // 1. ENT
// type LoaiVuKhiInputEnt struct {
// 	TenLoai string
// 	MoTa    *string
// }

// type LoaiVuKhiServiceEnt interface {
// 	CreateEnt(ctx context.Context, in LoaiVuKhiInputEnt) (*ent.LoaiVuKhi, error)
// 	GetIdEnt(ctx context.Context, id int) (*ent.LoaiVuKhi, error)
// 	GetAllEnt(ctx context.Context, q string) ([]*ent.LoaiVuKhi, error)
// 	UpdateEnt(ctx context.Context, id int, in LoaiVuKhiInputEnt) (*ent.LoaiVuKhi, error)
// 	DeleteEnt(ctx context.Context, id int) error

// 	GetTenLoaiEnt(ctx context.Context, q string) ([]string, error)
// }

// type loaiVuKhiServiceEnt struct{ repo repo.LoaiVuKhiRepo }

// func CreateLoaiVuKhiServiceEnt(r repo.LoaiVuKhiRepo) LoaiVuKhiServiceEnt {
// 	return &loaiVuKhiServiceEnt{repo: r}
// }

// func (s *loaiVuKhiServiceEnt) CreateEnt(ctx context.Context, in LoaiVuKhiInputEnt) (*ent.LoaiVuKhi, error) {
// 	return s.repo.CreateEnt(ctx, repo.CreateLoaiVKInput{
// 		TenLoai: in.TenLoai,
// 		MoTa:    in.MoTa,
// 	})
// }

// func (s *loaiVuKhiServiceEnt) GetIdEnt(ctx context.Context, id int) (*ent.LoaiVuKhi, error) {
// 	return s.repo.GetIdEnt(ctx, id)
// }

// func (s *loaiVuKhiServiceEnt) GetAllEnt(ctx context.Context, q string) ([]*ent.LoaiVuKhi, error) {
// 	return s.repo.GetAllEnt(ctx, q)
// }

// func (s *loaiVuKhiServiceEnt) UpdateEnt(ctx context.Context, id int, in LoaiVuKhiInputEnt) (*ent.LoaiVuKhi, error) {
// 	return s.repo.UpdateEnt(ctx, id, func(up *ent.LoaiVuKhiUpdateOne) *ent.LoaiVuKhiUpdateOne {
// 		return up.
// 			SetTenLoai(in.TenLoai).
// 			SetNillableMoTa(in.MoTa)
// 	})
// }

// func (s *loaiVuKhiServiceEnt) DeleteEnt(ctx context.Context, id int) error {
// 	return s.repo.DeleteEnt(ctx, id)
// }

// func (s *loaiVuKhiServiceEnt) GetTenLoaiEnt(ctx context.Context, q string) ([]string, error) {
// 	return s.repo.GetTenLoaiEnt(ctx, q)
// }

// // 2. GORM
// type LoaiVuKhiInputGorm struct {
// 	TenLoai string
// 	MoTa    *string
// }

// type LoaiVuKhiServiceGorm interface {
// 	CreateGorm(ctx context.Context, in LoaiVuKhiInputGorm) (*db.GLoaiVuKhi, error)
// 	GetIdGorm(ctx context.Context, id int) (*db.GLoaiVuKhi, error)
// 	GetAllGorm(ctx context.Context, q string) ([]db.GLoaiVuKhi, error)
// 	UpdateGorm(ctx context.Context, id int, in LoaiVuKhiInputGorm) (*db.GLoaiVuKhi, error)
// 	DeleteGorm(ctx context.Context, id int) error

// 	GetTenLoaiGorm(ctx context.Context, q string) ([]string, error)
// }

// type loaiVuKhiServiceGorm struct{ repo repo.LoaiVuKhiRepoGorm }

// func CreateLoaiVuKhiServiceGorm(r repo.LoaiVuKhiRepoGorm) LoaiVuKhiServiceGorm {
// 	return &loaiVuKhiServiceGorm{repo: r}
// }

// func (s *loaiVuKhiServiceGorm) CreateGorm(ctx context.Context, in LoaiVuKhiInputGorm) (*db.GLoaiVuKhi, error) {
// 	return s.repo.CreateGorm(ctx, repo.LoaiVuKhiGormInput{
// 		TenLoai: in.TenLoai,
// 		MoTa:    in.MoTa,
// 	})
// }

// func (s *loaiVuKhiServiceGorm) GetIdGorm(ctx context.Context, id int) (*db.GLoaiVuKhi, error) {
// 	return s.repo.GetIdGorm(ctx, id)
// }

// func (s *loaiVuKhiServiceGorm) GetAllGorm(ctx context.Context, q string) ([]db.GLoaiVuKhi, error) {
// 	return s.repo.GetAllGorm(ctx, q)
// }

// func (s *loaiVuKhiServiceGorm) UpdateGorm(ctx context.Context, id int, in LoaiVuKhiInputGorm) (*db.GLoaiVuKhi, error) {
// 	return s.repo.UpdateGorm(ctx, id, repo.LoaiVuKhiGormInput{
// 		TenLoai: in.TenLoai,
// 		MoTa:    in.MoTa,
// 	})
// }

// func (s *loaiVuKhiServiceGorm) DeleteGorm(ctx context.Context, id int) error {
// 	return s.repo.DeleteGorm(ctx, id)
// }

// func (s *loaiVuKhiServiceGorm) GetTenLoaiGorm(ctx context.Context, q string) ([]string, error) {
// 	return s.repo.GetTenLoaiGorm(ctx, q)
// }

// // 3. RAW
// type LoaiVuKhiInputRaw struct {
// 	TenLoai string
// 	MoTa    *string
// }

// type LoaiVuKhiServiceRaw interface {
// 	CreateRaw(ctx context.Context, in LoaiVuKhiInputRaw) (*repo.LoaiVuKhiRaw, error)
// 	GetIdRaw(ctx context.Context, id int) (*repo.LoaiVuKhiRaw, error)
// 	GetAllRaw(ctx context.Context, q string) ([]repo.LoaiVuKhiRaw, error)
// 	UpdateRaw(ctx context.Context, id int, in LoaiVuKhiInputRaw) (*repo.LoaiVuKhiRaw, error)
// 	DeleteRaw(ctx context.Context, id int) error

// 	GetTenLoaiRaw(ctx context.Context, q string) ([]string, error)
// }

// type loaiVuKhiServiceRaw struct{ repo repo.LoaiVuKhiRepoRaw }

// func CreateLoaiVuKhiServiceRaw(r repo.LoaiVuKhiRepoRaw) LoaiVuKhiServiceRaw {
// 	return &loaiVuKhiServiceRaw{repo: r}
// }

// func (s *loaiVuKhiServiceRaw) CreateRaw(ctx context.Context, in LoaiVuKhiInputRaw) (*repo.LoaiVuKhiRaw, error) {
// 	return s.repo.CreateRaw(ctx, repo.LoaiVuKhiRawInput{
// 		TenLoai: in.TenLoai,
// 		MoTa:    in.MoTa,
// 	})
// }

// func (s *loaiVuKhiServiceRaw) GetIdRaw(ctx context.Context, id int) (*repo.LoaiVuKhiRaw, error) {
// 	return s.repo.GetIdRaw(ctx, id)
// }

// func (s *loaiVuKhiServiceRaw) GetAllRaw(ctx context.Context, q string) ([]repo.LoaiVuKhiRaw, error) {
// 	return s.repo.GetAllRaw(ctx, q)
// }

// func (s *loaiVuKhiServiceRaw) UpdateRaw(ctx context.Context, id int, in LoaiVuKhiInputRaw) (*repo.LoaiVuKhiRaw, error) {
// 	return s.repo.UpdateRaw(ctx, id, repo.LoaiVuKhiRawInput{
// 		TenLoai: in.TenLoai,
// 		MoTa:    in.MoTa,
// 	})
// }

// func (s *loaiVuKhiServiceRaw) DeleteRaw(ctx context.Context, id int) error {
// 	return s.repo.DeleteRaw(ctx, id)
// }

// func (s *loaiVuKhiServiceRaw) GetTenLoaiRaw(ctx context.Context, q string) ([]string, error) {
// 	return s.repo.GetTenLoaiRaw(ctx, q)
// }
