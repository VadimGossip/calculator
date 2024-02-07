package writer

import (
	"context"
	"github.com/VadimGossip/calculator/api/internal/domain"
)

type service struct {
	repo Repository
}

type Service interface {
	CreateExpression(ctx context.Context, e *domain.Expression) error
	SaveExpressionResult(ctx context.Context, id int64, result int) error
	GetExpressions(ctx context.Context) ([]domain.Expression, error)
}

var _ Service = (*service)(nil)

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) CreateExpression(ctx context.Context, e *domain.Expression) error {
	return s.repo.CreateExpression(ctx, e)
}

func (s *service) SaveExpressionResult(ctx context.Context, id int64, result int) error {
	return s.repo.SaveExpressionResult(ctx, id, result)
}

func (s *service) GetExpressions(ctx context.Context) ([]domain.Expression, error) {
	return s.repo.GetExpressions(ctx)
}
