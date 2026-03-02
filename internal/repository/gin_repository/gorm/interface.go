package gin_repository

import (
	"context"
	"game/internal/models"
)

type VuKhiRepoGorm interface {
	CreateGorm(ctx context.Context, vk *models.VuKhis) (*models.VuKhis, error)
	GetAllGorm(ctx context.Context, kw string) ([]*models.VuKhis, error)
	DeleteByNameGorm(ctx context.Context, name string) error
	UpdateByNameGorm(ctx context.Context, name string, in models.VuKhiUpdateInput) (*models.VuKhis, error)
	GetAllTenVuKhiGorm(ctx context.Context) ([]string, error)
	GetAllSatThuongCoBanGorm(ctx context.Context) ([]int, error)
	GetAllTocDoDanhGorm(ctx context.Context) ([]float64, error)
	GetAllTamDanhGorm(ctx context.Context) ([]int, error)
	SearchGorm(ctx context.Context, req models.VuKhiSearchRequest, limit int, offset int, arrangeAsc string, arrangeDesc string) ([]*models.VuKhiResponse, int64, error)
}

// type DoHiemRepoGorm interface {
// 	CreateGorm(ctx context.Context, in TaoDoHiemInput) (*ent.DoHiem, error)
// 	GetByNameGorm(ctx context.Context, name string) (*ent.DoHiem, error)
// 	UpdateByNameGorm(ctx context.Context, name string, mutate func(*ent.DoHiemUpdateOne) *ent.DoHiemUpdateOne) (*ent.DoHiem, error)
// 	DeleteByNameGorm(ctx context.Context, name string) error
// 	GetAllGorm(ctx context.Context, q string) ([]*ent.DoHiem, error)
// 	GetAllTen(ctx context.Context) ([]string, error)
// 	GetAllSoLuong(ctx context.Context) ([]int, error)
// 	GetAllMauSac(ctx context.Context) ([]string, error)
// 	GetAllSatThuongBonus(ctx context.Context) ([]float64, error)
// 	GetAllTocDoDanhBonus(ctx context.Context) ([]float64, error)
// }

// type HeRepoGorm interface {
// 	CreateGorm(ctx context.Context, in CreateHeInput) (*ent.He, error)
// 	GetIdGorm(ctx context.Context, id int) (*ent.He, error)
// 	GetAllGorm(ctx context.Context, q string) ([]*ent.He, error)
// 	UpdateGorm(ctx context.Context, id int, mutate func(*ent.HeUpdateOne) *ent.HeUpdateOne) (*ent.He, error)
// 	DeleteGorm(ctx context.Context, id int) error
// 	GetByNameGorm(ctx context.Context, name string) (*ent.He, error)
// 	UpdateByNameGorm(ctx context.Context, name string, mutate func(*ent.HeUpdateOne) *ent.HeUpdateOne) (*ent.He, error)
// 	DeleteByNameGorm(ctx context.Context, name string) error
// 	GetAllNamesGorm(ctx context.Context, q string) ([]string, error)
// }

// type LoaiVuKhiRepoGorm interface {
// 	CreateGorm(ctx context.Context, in CreateLoaiVKInput) (*ent.LoaiVuKhi, error)
// 	GetIdGorm(ctx context.Context, id int) (*ent.LoaiVuKhi, error)
// 	GetAllGorm(ctx context.Context, q string) ([]*ent.LoaiVuKhi, error)
// 	UpdateGorm(ctx context.Context, id int, mutate func(*ent.LoaiVuKhiUpdateOne) *ent.LoaiVuKhiUpdateOne) (*ent.LoaiVuKhi, error)
// 	DeleteGorm(ctx context.Context, id int) error
// 	GetTenLoaiGorm(ctx context.Context, q string) ([]string, error)
// }
