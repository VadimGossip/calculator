package token

import (
	"context"
	"github.com/VadimGossip/calculator/dbagent/internal/domain"
)

type service struct {
	repo Repository
}

type Service interface {
	Create(ctx context.Context, token *domain.Token) error
	Get(ctx context.Context, token string) (*domain.Token, error)
}

var _ Service = (*service)(nil)

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, token *domain.Token) error {
	return s.repo.Create(ctx, token)
}

func (s *service) Get(ctx context.Context, tokenValue string) (*domain.Token, error) {
	return s.repo.Get(ctx, tokenValue)
}
