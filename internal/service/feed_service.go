// ===== internal/service/feed_service.go =====
package service

import (
	"context"
	"strings"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/repository"
)

type FeedService struct {
	repo *repository.FeedRepository
}

func NewFeedService(repo *repository.FeedRepository) *FeedService {
	return &FeedService{repo: repo}
}

func normalizeFeed(in *model.FeedInput) {
	in.Name = strings.TrimSpace(in.Name)
	if in.Brand != nil {
		b := strings.TrimSpace(*in.Brand)
		if b == "" {
			in.Brand = nil // brand kosong disimpan NULL
		} else {
			in.Brand = &b
		}
	}
}

func (s *FeedService) List(ctx context.Context, f model.FeedFilter) ([]model.Feed, int64, error) {
	f.Search = strings.TrimSpace(f.Search)
	return s.repo.List(ctx, f)
}

func (s *FeedService) Get(ctx context.Context, id int64) (*model.Feed, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *FeedService) Create(ctx context.Context, in model.FeedInput) (*model.Feed, error) {
	normalizeFeed(&in)
	return s.repo.Create(ctx, in)
}

func (s *FeedService) Update(ctx context.Context, id int64, in model.FeedInput) (*model.Feed, error) {
	normalizeFeed(&in)
	return s.repo.Update(ctx, id, in)
}

func (s *FeedService) Delete(ctx context.Context, id int64) error {
	return s.repo.SoftDelete(ctx, id)
}
