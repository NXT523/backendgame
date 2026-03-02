package service

import (
	"context"
	"game/internal/models"
)

type VuKhiServiceGorm interface {
	CreateGorm(ctx context.Context, vk *models.VuKhis) (*models.VuKhis, error)
	GetAllGorm(ctx context.Context, kw string) ([]*models.VuKhis, error)
	UpdateByNameGorm(ctx context.Context, name string, in models.VuKhiUpdateInput) (*models.VuKhis, error)
	DeleteByNameGorm(ctx context.Context, name string) error
	GetAllTenVuKhiGorm(ctx context.Context) ([]string, error)
	GetAllSatThuongCoBanGorm(ctx context.Context) ([]int, error)
	GetAllTocDoDanhGorm(ctx context.Context) ([]float64, error)
	GetAllTamDanhGorm(ctx context.Context) ([]int, error)
	SearchGorm(ctx context.Context, req models.VuKhiSearchRequest, page int, pageSize int, arrangeAsc string, arrangeDesc string) ([]*models.VuKhiResponse, int64, error)
}
