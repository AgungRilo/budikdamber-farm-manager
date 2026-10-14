package service

import (
	"context"
	"errors"
	"strings"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/repository"
)

var ErrChannelRequired = errors.New("kanal penjualan wajib diisi untuk penjualan (langsung/pasar/pengepul/online)")

type SaleService struct {
	repo *repository.SaleRepository
}

func NewSaleService(repo *repository.SaleRepository) *SaleService {
	return &SaleService{repo: repo}
}

func (s *SaleService) List(ctx context.Context, f model.SaleFilter) ([]model.Sale, int64, error) {
	return s.repo.List(ctx, f)
}

func (s *SaleService) Create(ctx context.Context, cycleID int64, in model.SaleInput, createdBy int64) (*model.SaleCreateResult, error) {
	if err := validateLogDate(in.SaleDate); err != nil {
		return nil, err
	}
	if in.SaleType == "" {
		in.SaleType = "sale"
	}

	if in.SaleType == "own_use" {
		// R07: konsumsi sendiri → tanpa harga, tanpa kanal, otomatis "lunas"
		in.PricePerFish, in.Channel = nil, nil
		paid := true
		in.IsPaid = &paid
	} else if in.Channel == nil {
		return nil, ErrChannelRequired
	}

	if in.IsPaid == nil {
		paid := true
		in.IsPaid = &paid
	}
	if in.BuyerName != nil {
		b := strings.TrimSpace(*in.BuyerName)
		in.BuyerName = &b
		if b == "" {
			in.BuyerName = nil
		}
	}
	in.Notes = trimNotes(in.Notes)
	return s.repo.Create(ctx, cycleID, in, createdBy)
}

func (s *SaleService) SetPaid(ctx context.Context, id int64, isPaid bool) (*model.Sale, error) {
	return s.repo.SetPaid(ctx, id, isPaid)
}

func (s *SaleService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
