package writer

import "context"

type service struct {
	repo Repository
}

type Service interface {
	CreateExpression(ctx context.Context, e *Expression) error
}

var _ Service = (*service)(nil)

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) CreateExpression(ctx context.Context, e *Expression) error {
	return s.repo.CreateExpression(ctx, e)
}