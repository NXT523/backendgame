package service

import (
	"context"

	"game/ent"
	"game/internal/db"
	"game/internal/repo"
)

// 1. ENT
type HeInputEnt struct {
	TenHe string
	MoTa  *string
}

type HeServiceEnt interface {
	CreateEnt(ctx context.Context, in HeInputEnt) (*ent.He, error)
	GetIdEnt(ctx context.Context, id int) (*ent.He, error)
	GetAllEnt(ctx context.Context, q string) ([]*ent.He, error)
	UpdateEnt(ctx context.Context, id int, in HeInputEnt) (*ent.He, error)
	DeleteEnt(ctx context.Context, id int) error

	GetByNameEnt(ctx context.Context, name string) (*ent.He, error)
	UpdateByNameEnt(ctx context.Context, name string, in HeInputEnt) (*ent.He, error)
	DeleteByNameEnt(ctx context.Context, name string) error

	GetAllNamesEnt(ctx context.Context, q string) ([]string, error)
}

type heServiceEnt struct{ repo repo.HeRepo }

func CreateHeServiceEnt(r repo.HeRepo) HeServiceEnt { return &heServiceEnt{repo: r} }

func (s *heServiceEnt) CreateEnt(ctx context.Context, in HeInputEnt) (*ent.He, error) {
	return s.repo.CreateEnt(ctx, repo.CreateHeInput{
		TenHe: in.TenHe,
		MoTa:  in.MoTa,
	})
}

func (s *heServiceEnt) GetIdEnt(ctx context.Context, id int) (*ent.He, error) {
	return s.repo.GetIdEnt(ctx, id)
}

func (s *heServiceEnt) GetAllEnt(ctx context.Context, q string) ([]*ent.He, error) {
	return s.repo.GetAllEnt(ctx, q)
}

func (s *heServiceEnt) UpdateEnt(ctx context.Context, id int, in HeInputEnt) (*ent.He, error) {
	return s.repo.UpdateEnt(ctx, id, func(up *ent.HeUpdateOne) *ent.HeUpdateOne {
		return up.
			SetTenHe(in.TenHe).
			SetNillableMoTa(in.MoTa)
	})
}

func (s *heServiceEnt) DeleteEnt(ctx context.Context, id int) error {
	return s.repo.DeleteEnt(ctx, id)
}

func (s *heServiceEnt) GetByNameEnt(ctx context.Context, name string) (*ent.He, error) {
	return s.repo.GetByNameEnt(ctx, name)
}

func (s *heServiceEnt) UpdateByNameEnt(ctx context.Context, name string, in HeInputEnt) (*ent.He, error) {
	return s.repo.UpdateByNameEnt(ctx, name, func(up *ent.HeUpdateOne) *ent.HeUpdateOne {
		return up.SetTenHe(in.TenHe).SetNillableMoTa(in.MoTa)
	})
}

func (s *heServiceEnt) DeleteByNameEnt(ctx context.Context, name string) error {
	return s.repo.DeleteByNameEnt(ctx, name)
}

func (s *heServiceEnt) GetAllNamesEnt(ctx context.Context, q string) ([]string, error) {
	return s.repo.GetAllNamesEnt(ctx, q)
}

// 2. GORM
type HeInputGorm struct {
	TenHe string
	MoTa  *string
}

type HeServiceGorm interface {
	CreateGorm(ctx context.Context, in HeInputGorm) (*db.GHe, error)
	GetIdGorm(ctx context.Context, id int) (*db.GHe, error)
	GetAllGorm(ctx context.Context, q string) ([]db.GHe, error)
	UpdateGorm(ctx context.Context, id int, in HeInputGorm) (*db.GHe, error)
	DeleteGorm(ctx context.Context, id int) error

	GetByNameGorm(ctx context.Context, name string) (*db.GHe, error)
	UpdateByNameGorm(ctx context.Context, name string, in HeInputGorm) (*db.GHe, error)
	DeleteByNameGorm(ctx context.Context, name string) error

	GetAllNamesGorm(ctx context.Context, q string) ([]string, error)
}

type heServiceGorm struct{ repo repo.HeRepoGorm }

func CreateHeServiceGorm(r repo.HeRepoGorm) HeServiceGorm { return &heServiceGorm{repo: r} }

func (s *heServiceGorm) CreateGorm(ctx context.Context, in HeInputGorm) (*db.GHe, error) {
	return s.repo.CreateGorm(ctx, repo.HeGormInput{TenHe: in.TenHe, MoTa: in.MoTa})
}
func (s *heServiceGorm) GetIdGorm(ctx context.Context, id int) (*db.GHe, error) {
	return s.repo.GetIdGorm(ctx, id)
}
func (s *heServiceGorm) GetAllGorm(ctx context.Context, q string) ([]db.GHe, error) {
	return s.repo.GetAllGorm(ctx, q)
}
func (s *heServiceGorm) UpdateGorm(ctx context.Context, id int, in HeInputGorm) (*db.GHe, error) {
	return s.repo.UpdateGorm(ctx, id, repo.HeGormInput{TenHe: in.TenHe, MoTa: in.MoTa})
}
func (s *heServiceGorm) DeleteGorm(ctx context.Context, id int) error {
	return s.repo.DeleteGorm(ctx, id)
}

func (s *heServiceGorm) GetByNameGorm(ctx context.Context, name string) (*db.GHe, error) {
	return s.repo.GetByNameGorm(ctx, name)
}

func (s *heServiceGorm) UpdateByNameGorm(ctx context.Context, name string, in HeInputGorm) (*db.GHe, error) {
	return s.repo.UpdateByNameGorm(ctx, name, repo.HeGormInput{
		TenHe: in.TenHe,
		MoTa:  in.MoTa,
	})
}

func (s *heServiceGorm) DeleteByNameGorm(ctx context.Context, name string) error {
	return s.repo.DeleteByNameGorm(ctx, name)
}

func (s *heServiceGorm) GetAllNamesGorm(ctx context.Context, q string) ([]string, error) {
	return s.repo.GetAllNamesGorm(ctx, q)
}

// 3. RAW
type HeInputRaw struct {
	TenHe string
	MoTa  *string
}

type HeServiceRaw interface {
	CreateRaw(ctx context.Context, in HeInputRaw) (*repo.HeRaw, error)
	GetIdRaw(ctx context.Context, id int) (*repo.HeRaw, error)
	GetAllRaw(ctx context.Context, q string) ([]repo.HeRaw, error)
	UpdateRaw(ctx context.Context, id int, in HeInputRaw) (*repo.HeRaw, error)
	DeleteRaw(ctx context.Context, id int) error

	GetByNameRaw(ctx context.Context, name string) (*repo.HeRaw, error)
	UpdateByNameRaw(ctx context.Context, name string, in HeInputRaw) (*repo.HeRaw, error)
	DeleteByNameRaw(ctx context.Context, name string) error

	GetAllNamesRaw(ctx context.Context, q string) ([]string, error)
}

type heServiceRaw struct{ repo repo.HeRepoRaw }

func CreateHeServiceRaw(r repo.HeRepoRaw) HeServiceRaw { return &heServiceRaw{repo: r} }

func (s *heServiceRaw) CreateRaw(ctx context.Context, in HeInputRaw) (*repo.HeRaw, error) {
	return s.repo.CreateRaw(ctx, repo.HeRawInput{TenHe: in.TenHe, MoTa: in.MoTa})
}
func (s *heServiceRaw) GetIdRaw(ctx context.Context, id int) (*repo.HeRaw, error) {
	return s.repo.GetIdRaw(ctx, id)
}
func (s *heServiceRaw) GetAllRaw(ctx context.Context, q string) ([]repo.HeRaw, error) {
	return s.repo.GetAllRaw(ctx, q)
}
func (s *heServiceRaw) UpdateRaw(ctx context.Context, id int, in HeInputRaw) (*repo.HeRaw, error) {
	return s.repo.UpdateRaw(ctx, id, repo.HeRawInput{TenHe: in.TenHe, MoTa: in.MoTa})
}
func (s *heServiceRaw) DeleteRaw(ctx context.Context, id int) error {
	return s.repo.DeleteRaw(ctx, id)
}

func (s *heServiceRaw) GetByNameRaw(ctx context.Context, name string) (*repo.HeRaw, error) {
	return s.repo.GetByNameRaw(ctx, name)
}

func (s *heServiceRaw) UpdateByNameRaw(ctx context.Context, name string, in HeInputRaw) (*repo.HeRaw, error) {
	return s.repo.UpdateByNameRaw(ctx, name, repo.HeRawInput{
		TenHe: in.TenHe,
		MoTa:  in.MoTa,
	})
}

func (s *heServiceRaw) DeleteByNameRaw(ctx context.Context, name string) error {
	return s.repo.DeleteByNameRaw(ctx, name)
}

func (s *heServiceRaw) GetAllNamesRaw(ctx context.Context, q string) ([]string, error) {
	return s.repo.GetAllNamesRaw(ctx, q)
}
