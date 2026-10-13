package service

import (
	"context"
	"strings"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/repository"
)

type PondService struct {
	repo *repository.PondRepository
}

func NewPondService(repo *repository.PondRepository) *PondService {
	return &PondService{repo: repo}
}

func normalizePond(in *model.PondInput) {
	in.Name = strings.TrimSpace(in.Name)
	if in.IsActive == nil {
		active := true
		in.IsActive = &active
	}
}

func (s *PondService) List(ctx context.Context, f model.PondFilter) ([]model.Pond, int64, error) {
	f.Search = strings.TrimSpace(f.Search)
	return s.repo.List(ctx, f)
}

func (s *PondService) Get(ctx context.Context, id int64) (*model.Pond, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *PondService) Create(ctx context.Context, in model.PondInput) (*model.Pond, error) {
	normalizePond(&in)
	return s.repo.Create(ctx, in)
}

func (s *PondService) Update(ctx context.Context, id int64, in model.PondInput) (*model.Pond, error) {
	normalizePond(&in)
	return s.repo.Update(ctx, id, in)
}

func (s *PondService) Delete(ctx context.Context, id int64) error {
	return s.repo.SoftDelete(ctx, id)
}
