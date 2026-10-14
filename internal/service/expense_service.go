package service

import (
	"context"
	"strings"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/repository"
)

type ExpenseService struct {
	repo *repository.ExpenseRepository
}

func NewExpenseService(repo *repository.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

func normalizeExpense(in *model.ExpenseInput) error {
	if err := validateLogDate(in.ExpenseDate); err != nil {
		return err
	}
	in.Description = strings.TrimSpace(in.Description)
	return nil
}

func (s *ExpenseService) List(ctx context.Context, f model.ExpenseFilter) ([]model.Expense, int64, error) {
	return s.repo.List(ctx, f)
}

func (s *ExpenseService) Get(ctx context.Context, id int64) (*model.Expense, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ExpenseService) Create(ctx context.Context, in model.ExpenseInput, createdBy int64) (*model.Expense, error) {
	if err := normalizeExpense(&in); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, in, createdBy)
}

func (s *ExpenseService) Update(ctx context.Context, id int64, in model.ExpenseInput) (*model.Expense, error) {
	if err := normalizeExpense(&in); err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, id, in)
}

func (s *ExpenseService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
