package service

import (
	"context"

	"game/ent"
	"game/internal/db"
	"game/internal/repo"
)

type VuKhiInputEnt struct {
	TenVuKhi       string
	SatThuongCoBan int
	TocDoDanh      float64
	TamDanh        int
	MoTa           *string
	MaLoai         int
	MaDoHiem       int
	MaHe           int
}

type VuKhiView struct {
	TenVuKhi       string  `json:"ten_vu_khi"`
	TenLoai        string  `json:"ten_loai,omitempty"`
	TenHe          string  `json:"ten_he,omitempty"`
	TenDoHiem      string  `json:"ten_do_hiem,omitempty"`
	SatThuongCoBan int     `json:"sat_thuong_co_ban"`
	SatThuongBonus float64 `json:"sat_thuong_bonus"`
	TongSatThuong  float64 `json:"tong_sat_thuong"`
	TocDoDanh      float64 `json:"toc_do_danh"`
	TocDoDanhBonus float64 `json:"toc_do_danh_bonus"`
	TongTocDoDanh  float64 `json:"tong_toc_do_danh"`
	TamDanh        int     `json:"tam_danh,omitempty"`
}

type VuKhiServiceEnt interface {
	CreateEnt(ctx context.Context, in VuKhiInputEnt) (*ent.VuKhi, error)
	UpdateByNameEnt(ctx context.Context, name string, in VuKhiInputEnt) (*ent.VuKhi, error)
	DeleteByNameEnt(ctx context.Context, name string) error

	GetAllEnt(ctx context.Context, q string) ([]*ent.VuKhi, error)
	GetAllTenVuKhiEnt(ctx context.Context) ([]string, error)
	GetAllSatThuongCoBanEnt(ctx context.Context) ([]int, error)
	GetAllTocDoDanhEnt(ctx context.Context) ([]float64, error)
	GetAllTamDanhEnt(ctx context.Context) ([]int, error)

	SearchEnt(ctx context.Context, in repo.VuKhiSearchRequest) ([]*ent.VuKhi, error)
}

type vuKhiServiceEnt struct{ repo repo.VuKhiRepo }

func CreateVuKhiServiceEnt(r repo.VuKhiRepo) VuKhiServiceEnt { return &vuKhiServiceEnt{repo: r} }

func (s *vuKhiServiceEnt) CreateEnt(ctx context.Context, in VuKhiInputEnt) (*ent.VuKhi, error) {
	return s.repo.CreateEnt(ctx, repo.CreateVuKhiInput{
		TenVuKhi:       in.TenVuKhi,
		SatThuongCoBan: in.SatThuongCoBan,
		TocDoDanh:      in.TocDoDanh,
		TamDanh:        in.TamDanh,
		MoTa:           in.MoTa,
		MaLoai:         in.MaLoai,
		MaDoHiem:       in.MaDoHiem,
		MaHe:           in.MaHe,
	})
}

func (s *vuKhiServiceEnt) GetAllEnt(ctx context.Context, q string) ([]*ent.VuKhi, error) {
	return s.repo.GetAllEnt(ctx, q)
}

func (s vuKhiServiceEnt) UpdateByNameEnt(ctx context.Context, name string, in VuKhiInputEnt) (*ent.VuKhi, error) {
	return s.repo.UpdateByNameEnt(ctx, name, func(up *ent.VuKhiUpdateOne) *ent.VuKhiUpdateOne {
		return up.
			SetTenVuKhi(in.TenVuKhi).
			SetSatThuongCoBan(in.SatThuongCoBan).
			SetTocDoDanh(in.TocDoDanh).
			SetTamDanh(in.TamDanh).
			SetNillableMoTa(in.MoTa).
			SetMaLoai(in.MaLoai).
			SetMaDoHiem(in.MaDoHiem).
			SetMaHe(in.MaHe)
	})
}

func (s vuKhiServiceEnt) DeleteByNameEnt(ctx context.Context, name string) error {
	return s.repo.DeleteByNameEnt(ctx, name)
}

func (s *vuKhiServiceEnt) GetAllTenVuKhiEnt(ctx context.Context) ([]string, error) {
	return s.repo.GetAllTenVuKhiEnt(ctx)
}

func (s *vuKhiServiceEnt) GetAllSatThuongCoBanEnt(ctx context.Context) ([]int, error) {
	return s.repo.GetAllSatThuongCoBanEnt(ctx)
}

func (s *vuKhiServiceEnt) GetAllTocDoDanhEnt(ctx context.Context) ([]float64, error) {
	return s.repo.GetAllTocDoDanhEnt(ctx)
}

func (s *vuKhiServiceEnt) GetAllTamDanhEnt(ctx context.Context) ([]int, error) {
	return s.repo.GetAllTamDanhEnt(ctx)
}

func (s *vuKhiServiceEnt) SearchEnt(ctx context.Context, in repo.VuKhiSearchRequest) ([]*ent.VuKhi, error) {
	return s.repo.SearchEnt(ctx, in)
}

// 2. GORM
type VuKhiInputGorm struct {
	TenVuKhi       string
	SatThuongCoBan int
	TocDoDanh      float64
	TamDanh        int
	MoTa           *string
	MaLoai         int
	MaDoHiem       int
	MaHe           int
}

type VuKhiServiceGorm interface {
	CreateGorm(ctx context.Context, in VuKhiInputGorm) (*db.GVuKhi, error)
	GetAllGorm(ctx context.Context, q string) ([]db.GVuKhi, error)

	// GetByNameGorm(ctx context.Context, name string) (*db.GVuKhi, error)
	UpdateByNameGorm(ctx context.Context, name string, in VuKhiInputGorm) (*db.GVuKhi, error)
	DeleteByNameGorm(ctx context.Context, name string) error

	GetAllTenVuKhiGorm(ctx context.Context) ([]string, error)
	GetAllSatThuongCoBanGorm(ctx context.Context) ([]int, error)
	GetAllTocDoDanhGorm(ctx context.Context) ([]float64, error)
	GetAllTamDanhGorm(ctx context.Context) ([]int, error)

	SearchGorm(ctx context.Context, in repo.VuKhiSearchRequest) ([]db.GVuKhi, error)
}

type vuKhiServiceGorm struct{ repo repo.VuKhiRepoGorm }

func CreateVuKhiServiceGorm(r repo.VuKhiRepoGorm) VuKhiServiceGorm { return &vuKhiServiceGorm{repo: r} }

func (s *vuKhiServiceGorm) CreateGorm(ctx context.Context, in VuKhiInputGorm) (*db.GVuKhi, error) {
	return s.repo.CreateGorm(ctx, repo.VuKhiGormInput{
		TenVuKhi:       in.TenVuKhi,
		SatThuongCoBan: in.SatThuongCoBan,
		TocDoDanh:      in.TocDoDanh,
		TamDanh:        in.TamDanh,
		MoTa:           in.MoTa,
		MaLoai:         in.MaLoai,
		MaDoHiem:       in.MaDoHiem,
		MaHe:           in.MaHe,
	})
}

func (s *vuKhiServiceGorm) GetAllGorm(ctx context.Context, q string) ([]db.GVuKhi, error) {
	return s.repo.GetAllGorm(ctx, q)
}

func (s *vuKhiServiceGorm) UpdateByNameGorm(ctx context.Context, name string, in VuKhiInputGorm) (*db.GVuKhi, error) {
	return s.repo.UpdateByNameGorm(ctx, name, repo.VuKhiGormInput{
		TenVuKhi:       in.TenVuKhi,
		SatThuongCoBan: in.SatThuongCoBan,
		TocDoDanh:      in.TocDoDanh,
		TamDanh:        in.TamDanh,
		MoTa:           in.MoTa,
		MaLoai:         in.MaLoai,
		MaDoHiem:       in.MaDoHiem,
		MaHe:           in.MaHe,
	})
}

func (s *vuKhiServiceGorm) DeleteByNameGorm(ctx context.Context, name string) error {
	return s.repo.DeleteByNameGorm(ctx, name)
}

func (s *vuKhiServiceGorm) GetAllTenVuKhiGorm(ctx context.Context) ([]string, error) {
	return s.repo.GetAllTenVuKhiGorm(ctx)
}
func (s *vuKhiServiceGorm) GetAllSatThuongCoBanGorm(ctx context.Context) ([]int, error) {
	return s.repo.GetAllSatThuongCoBanGorm(ctx)
}
func (s *vuKhiServiceGorm) GetAllTocDoDanhGorm(ctx context.Context) ([]float64, error) {
	return s.repo.GetAllTocDoDanhGorm(ctx)
}
func (s *vuKhiServiceGorm) GetAllTamDanhGorm(ctx context.Context) ([]int, error) {
	return s.repo.GetAllTamDanhGorm(ctx)
}

func (s *vuKhiServiceGorm) SearchGorm(ctx context.Context, in repo.VuKhiSearchRequest) ([]db.GVuKhi, error) {
	return s.repo.SearchGorm(ctx, in)
}

// 3. RAW
type VuKhiInputRaw struct {
	TenVuKhi       string
	SatThuongCoBan int
	TocDoDanh      float64
	TamDanh        int
	MoTa           *string
	MaLoai         int
	MaDoHiem       int
	MaHe           int
}

// Chỉ giữ các hàm khớp router rút gọn
type VuKhiServiceRaw interface {
	CreateRaw(ctx context.Context, in VuKhiInputRaw) (*repo.VuKhiRaw, error)
	GetAllRaw(ctx context.Context, q string) ([]repo.VuKhiRaw, error)

	UpdateByNameRaw(ctx context.Context, name string, in VuKhiInputRaw) (*repo.VuKhiRaw, error)
	DeleteByNameRaw(ctx context.Context, name string) error

	GetAllTenVuKhiRaw(ctx context.Context) ([]string, error)
	GetAllSatThuongCoBanRaw(ctx context.Context) ([]int, error)
	GetAllTocDoDanhRaw(ctx context.Context) ([]float64, error)
	GetAllTamDanhRaw(ctx context.Context) ([]int, error)

	// search tổng hợp
	SearchRaw(ctx context.Context, req repo.VuKhiSearchRequest) ([]repo.VuKhiRawJoin, error)
}

type vuKhiServiceRaw struct{ repo repo.VuKhiRepoRaw }

func CreateVuKhiServiceRaw(r repo.VuKhiRepoRaw) VuKhiServiceRaw { return &vuKhiServiceRaw{repo: r} }

func (s *vuKhiServiceRaw) CreateRaw(ctx context.Context, in VuKhiInputRaw) (*repo.VuKhiRaw, error) {
	return s.repo.CreateRaw(ctx, repo.VuKhiRawInput{
		TenVuKhi:       in.TenVuKhi,
		SatThuongCoBan: in.SatThuongCoBan,
		TocDoDanh:      in.TocDoDanh,
		TamDanh:        in.TamDanh,
		MoTa:           in.MoTa,
		MaLoai:         in.MaLoai,
		MaDoHiem:       in.MaDoHiem,
		MaHe:           in.MaHe,
	})
}

func (s *vuKhiServiceRaw) GetAllRaw(ctx context.Context, q string) ([]repo.VuKhiRaw, error) {
	return s.repo.GetAllRaw(ctx, q)
}

func (s *vuKhiServiceRaw) UpdateByNameRaw(ctx context.Context, name string, in VuKhiInputRaw) (*repo.VuKhiRaw, error) {
	return s.repo.UpdateByNameRaw(ctx, name, repo.VuKhiRawInput{
		TenVuKhi:       in.TenVuKhi,
		SatThuongCoBan: in.SatThuongCoBan,
		TocDoDanh:      in.TocDoDanh,
		TamDanh:        in.TamDanh,
		MoTa:           in.MoTa,
		MaLoai:         in.MaLoai,
		MaDoHiem:       in.MaDoHiem,
		MaHe:           in.MaHe,
	})
}

func (s *vuKhiServiceRaw) DeleteByNameRaw(ctx context.Context, name string) error {
	return s.repo.DeleteByNameRaw(ctx, name)
}

func (s *vuKhiServiceRaw) GetAllTenVuKhiRaw(ctx context.Context) ([]string, error) {
	return s.repo.GetAllTenVuKhiRaw(ctx)
}
func (s *vuKhiServiceRaw) GetAllSatThuongCoBanRaw(ctx context.Context) ([]int, error) {
	return s.repo.GetAllSatThuongCoBanRaw(ctx)
}
func (s *vuKhiServiceRaw) GetAllTocDoDanhRaw(ctx context.Context) ([]float64, error) {
	return s.repo.GetAllTocDoDanhRaw(ctx)
}
func (s *vuKhiServiceRaw) GetAllTamDanhRaw(ctx context.Context) ([]int, error) {
	return s.repo.GetAllTamDanhRaw(ctx)
}

func (s *vuKhiServiceRaw) SearchRaw(ctx context.Context, req repo.VuKhiSearchRequest) ([]repo.VuKhiRawJoin, error) {
	return s.repo.SearchRaw(ctx, req)
}
