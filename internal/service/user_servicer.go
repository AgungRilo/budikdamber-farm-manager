package service

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/model"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/repository"
)

var ErrSelfAction = errors.New("tidak bisa menonaktifkan atau menurunkan role akun sendiri")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func hashPassword(pw string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(h), err
}

func (s *UserService) List(ctx context.Context, f model.UserFilter) ([]model.User, int64, error) {
	f.Search = strings.TrimSpace(f.Search)
	return s.repo.List(ctx, f)
}

func (s *UserService) Get(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) Create(ctx context.Context, in model.UserCreateInput) (*model.User, error) {
	role := in.Role
	if role == "" {
		role = "worker"
	}
	hash, err := hashPassword(in.Password)
	if err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, strings.TrimSpace(in.Name), strings.ToLower(strings.TrimSpace(in.Email)), hash, role)
}

func (s *UserService) Update(ctx context.Context, actorID, id int64, in model.UserUpdateInput) (*model.User, error) {
	if actorID == id && in.Role != "owner" {
		return nil, ErrSelfAction
	}
	var newHash *string
	if in.Password != nil {
		h, err := hashPassword(*in.Password)
		if err != nil {
			return nil, err
		}
		newHash = &h
	}
	return s.repo.Update(ctx, id, strings.TrimSpace(in.Name), strings.ToLower(strings.TrimSpace(in.Email)), in.Role, newHash)
}

func (s *UserService) SetActive(ctx context.Context, actorID, id int64, active bool) (*model.User, error) {
	if actorID == id && !active {
		return nil, ErrSelfAction
	}
	return s.repo.SetActive(ctx, id, active)
}
