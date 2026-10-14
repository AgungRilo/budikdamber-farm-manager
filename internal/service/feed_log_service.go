package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/repository"
)

var ErrLogDateInFuture = errors.New("tanggal log tidak boleh di masa depan")

type FeedLogService struct {
	repo *repository.FeedLogRepository
}

func NewFeedLogService(repo *repository.FeedLogRepository) *FeedLogService {
	return &FeedLogService{repo: repo}
}

// validateLogDate dipakai juga oleh log kematian.
func validateLogDate(date string) error {
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return err
	}
	if isFutureDate(date) {
		return ErrLogDateInFuture
	}
	return nil
}

func trimNotes(n *string) *string {
	if n == nil {
		return nil
	}
	t := strings.TrimSpace(*n)
	if t == "" {
		return nil
	}
	return &t
}

func (s *FeedLogService) List(ctx context.Context, f model.LogFilter) ([]model.FeedLog, int64, error) {
	return s.repo.List(ctx, f)
}

func (s *FeedLogService) Create(ctx context.Context, cycleID int64, in model.FeedLogInput, createdBy int64) (*model.FeedLog, error) {
	if err := validateLogDate(in.LogDate); err != nil {
		return nil, err
	}
	in.Notes = trimNotes(in.Notes)
	return s.repo.Create(ctx, cycleID, in, createdBy)
}

func (s *FeedLogService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
