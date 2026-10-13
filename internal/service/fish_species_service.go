// ===== internal/service/fish_species_service.go =====
package service

import (
	"context"
	"strings"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/repository"
)

type FishSpeciesService struct {
	repo *repository.FishSpeciesRepository
}

func NewFishSpeciesService(repo *repository.FishSpeciesRepository) *FishSpeciesService {
	return &FishSpeciesService{repo: repo}
}

func (s *FishSpeciesService) List(ctx context.Context, f model.FishSpeciesFilter) ([]model.FishSpecies, int64, error) {
	f.Search = strings.TrimSpace(f.Search)
	return s.repo.List(ctx, f)
}

func (s *FishSpeciesService) Get(ctx context.Context, id int64) (*model.FishSpecies, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *FishSpeciesService) Create(ctx context.Context, in model.FishSpeciesInput) (*model.FishSpecies, error) {
	in.Name = strings.TrimSpace(in.Name)
	return s.repo.Create(ctx, in)
}

func (s *FishSpeciesService) Update(ctx context.Context, id int64, in model.FishSpeciesInput) (*model.FishSpecies, error) {
	in.Name = strings.TrimSpace(in.Name)
	return s.repo.Update(ctx, id, in)
}

func (s *FishSpeciesService) Delete(ctx context.Context, id int64) error {
	return s.repo.SoftDelete(ctx, id)
}
