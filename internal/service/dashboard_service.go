package service

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/repository"
)

type DashboardService struct {
	repo *repository.DashboardRepository
}

func NewDashboardService(repo *repository.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

// Get menjalankan 5 query yang tidak saling bergantung secara PARALEL.
//
//   - Tiap g.Go(...) berjalan di goroutine sendiri.
//   - Tiap goroutine menulis ke field yang BERBEDA di struct d, jadi tidak ada
//     data race dan tidak perlu mutex.
//   - Kalau satu query gagal, ctx dibatalkan → query lain ikut berhenti,
//     dan g.Wait() mengembalikan error pertama.
//   - pgxpool aman dipakai bersamaan; tiap goroutine meminjam koneksinya sendiri
//     (MaxConns 10, cukup untuk 5 query sekaligus).
func (s *DashboardService) Get(ctx context.Context, f model.DashboardFilter) (*model.Dashboard, error) {
	f.Today = todayWIB()
	d := &model.Dashboard{}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { return s.repo.Stock(ctx, &d.Stock) })
	g.Go(func() error { return s.repo.Finance(ctx, &d.Finance) })
	g.Go(func() error { return s.repo.SalesByChannel(ctx, f, &d.SalesByChannel) })
	g.Go(func() error { return s.repo.Receivables(ctx, f, &d.Receivables) })
	g.Go(func() error { return s.repo.ActiveCycles(ctx, f, &d.ActiveCycles) })

	if err := g.Wait(); err != nil {
		return nil, err
	}
	d.GeneratedAt = time.Now().In(wib)
	return d, nil
}
