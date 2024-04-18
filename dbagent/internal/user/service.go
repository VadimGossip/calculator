package user

import (
	"context"
	"github.com/VadimGossip/calculator/dbagent/internal/domain"
)

type service struct {
	repo Repository
}

type Service interface {
	Create(ctx context.Context, user *domain.User) error
	GetByCredentials(ctx context.Context, login string, password string) (*domain.User, error)
}

var _ Service = (*service)(nil)

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, user *domain.User) error {
	return s.repo.Create(ctx, user)
}

func (s *service) GetByCredentials(ctx context.Context, login, password string) (*domain.User, error) {
	return s.repo.GetByCredentials(ctx, login, password)
}
