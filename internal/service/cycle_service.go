package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/repository"
)

var ErrStartDateInFuture = errors.New("tanggal tebar tidak boleh di masa depan")

type CycleService struct {
	repo *repository.CycleRepository
}

func NewCycleService(repo *repository.CycleRepository) *CycleService {
	return &CycleService{repo: repo}
}

func (s *CycleService) List(ctx context.Context, f model.CycleFilter) ([]model.Cycle, int64, error) {
	return s.repo.List(ctx, f)
}

func (s *CycleService) Get(ctx context.Context, id int64) (*model.Cycle, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *CycleService) Create(ctx context.Context, in model.CycleInput, createdBy int64) (*model.Cycle, error) {
	start, _ := time.Parse("2006-01-02", in.StartDate) // format sudah divalidasi di handler
	if start.After(time.Now()) {
		return nil, ErrStartDateInFuture
	}
	if in.Notes != nil {
		n := strings.TrimSpace(*in.Notes)
		in.Notes = &n
	}
	return s.repo.Create(ctx, in, createdBy)
}
