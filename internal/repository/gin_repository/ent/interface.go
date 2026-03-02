package gin_repository

import (
	"context"
	"game/internal/models"
)

type VuKhiRepoEnt interface {
	CreateEnt(ctx context.Context, vk *models.VuKhis) (*models.VuKhis, error)
	GetAllEnt(ctx context.Context, kw string) ([]*models.VuKhis, error)
	DeleteByNameEnt(ctx context.Context, name string) error
	UpdateByNameEnt(ctx context.Context, name string, in models.VuKhiUpdateInput) (*models.VuKhis, error)
	GetAllTenVuKhiEnt(ctx context.Context) ([]string, error)
	GetAllSatThuongCoBanEnt(ctx context.Context) ([]int, error)
	GetAllTocDoDanhEnt(ctx context.Context) ([]float64, error)
	GetAllTamDanhEnt(ctx context.Context) ([]int, error)
	SearchEnt(ctx context.Context, req models.VuKhiSearchRequest, limit int, offset int, arrangeAsc string, arrangeDesc string) ([]*models.VuKhiResponse, int64, error)
}

// type DoHiemRepoEnt interface {
// 	CreateEnt(ctx context.Context, in TaoDoHiemInput) (*ent.DoHiem, error)
// 	GetByNameEnt(ctx context.Context, name string) (*ent.DoHiem, error)
// 	UpdateByNameEnt(ctx context.Context, name string, mutate func(*ent.DoHiemUpdateOne) *ent.DoHiemUpdateOne) (*ent.DoHiem, error)
// 	DeleteByNameEnt(ctx context.Context, name string) error
// 	GetAllEnt(ctx context.Context, q string) ([]*ent.DoHiem, error)
// 	GetAllTen(ctx context.Context) ([]string, error)
// 	GetAllSoLuong(ctx context.Context) ([]int, error)
// 	GetAllMauSac(ctx context.Context) ([]string, error)
// 	GetAllSatThuongBonus(ctx context.Context) ([]float64, error)
// 	GetAllTocDoDanhBonus(ctx context.Context) ([]float64, error)
// }

// type HeRepoEnt interface {
// 	CreateEnt(ctx context.Context, in CreateHeInput) (*ent.He, error)
// 	GetIdEnt(ctx context.Context, id int) (*ent.He, error)
// 	GetAllEnt(ctx context.Context, q string) ([]*ent.He, error)
// 	UpdateEnt(ctx context.Context, id int, mutate func(*ent.HeUpdateOne) *ent.HeUpdateOne) (*ent.He, error)
// 	DeleteEnt(ctx context.Context, id int) error
// 	GetByNameEnt(ctx context.Context, name string) (*ent.He, error)
// 	UpdateByNameEnt(ctx context.Context, name string, mutate func(*ent.HeUpdateOne) *ent.HeUpdateOne) (*ent.He, error)
// 	DeleteByNameEnt(ctx context.Context, name string) error
// 	GetAllNamesEnt(ctx context.Context, q string) ([]string, error)
// }

// type LoaiVuKhiRepoEnt interface {
// 	CreateEnt(ctx context.Context, in CreateLoaiVKInput) (*ent.LoaiVuKhi, error)
// 	GetIdEnt(ctx context.Context, id int) (*ent.LoaiVuKhi, error)
// 	GetAllEnt(ctx context.Context, q string) ([]*ent.LoaiVuKhi, error)
// 	UpdateEnt(ctx context.Context, id int, mutate func(*ent.LoaiVuKhiUpdateOne) *ent.LoaiVuKhiUpdateOne) (*ent.LoaiVuKhi, error)
// 	DeleteEnt(ctx context.Context, id int) error
// 	GetTenLoaiEnt(ctx context.Context, q string) ([]string, error)
// }
