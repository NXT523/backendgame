package service

import (
	"context"
	"game/internal/models"
)

type VuKhiServiceEnt interface {
	CreateEnt(ctx context.Context, vk *models.VuKhis) (*models.VuKhis, error)
	GetAllEnt(ctx context.Context, kw string) ([]*models.VuKhis, error)
	UpdateByNameEnt(ctx context.Context, name string, in models.VuKhiUpdateInput) (*models.VuKhis, error)
	DeleteByNameEnt(ctx context.Context, name string) error
	GetAllTenVuKhiEnt(ctx context.Context) ([]string, error)
	GetAllSatThuongCoBanEnt(ctx context.Context) ([]int, error)
	GetAllTocDoDanhEnt(ctx context.Context) ([]float64, error)
	GetAllTamDanhEnt(ctx context.Context) ([]int, error)
	SearchEnt(ctx context.Context, req models.VuKhiSearchRequest, page int, pageSize int, arrangeAsc string, arrangeDesc string) ([]*models.VuKhiResponse, int64, error)
}
