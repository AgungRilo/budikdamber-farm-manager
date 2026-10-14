package service

import (
	"context"
	"strings"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/repository"
)

type MortalityLogService struct {
	repo *repository.MortalityLogRepository
}

func NewMortalityLogService(repo *repository.MortalityLogRepository) *MortalityLogService {
	return &MortalityLogService{repo: repo}
}

func (s *MortalityLogService) List(ctx context.Context, f model.LogFilter) ([]model.MortalityLog, int64, error) {
	return s.repo.List(ctx, f)
}

func (s *MortalityLogService) Create(ctx context.Context, cycleID int64, in model.MortalityLogInput, createdBy int64) (*model.MortalityCreateResult, error) {
	if err := validateLogDate(in.LogDate); err != nil {
		return nil, err
	}
	in.Notes = trimNotes(in.Notes)
	if in.Cause != nil {
		c := strings.TrimSpace(*in.Cause)
		in.Cause = &c
		if c == "" {
			in.Cause = nil
		}
	}
	return s.repo.Create(ctx, cycleID, in, createdBy)
}

func (s *MortalityLogService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
