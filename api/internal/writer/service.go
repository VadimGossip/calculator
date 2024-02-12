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
	SaveAgentHeartbeat(ctx context.Context, name string) error
	GetAgents(ctx context.Context) ([]domain.Agent, error)
	SaveOperationDuration(ctx context.Context, name string, duration uint16) error
	GetOperationDurations(ctx context.Context) ([]domain.OperationDuration, error)
	CreateSubExpression(ctx context.Context, s *domain.SubExpression) error
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

func (s *service) SaveAgentHeartbeat(ctx context.Context, name string) error {
	agent, err := s.repo.GetAgent(ctx, name)
	if err != nil {
		return err
	}
	if agent.Name == "" {
		return s.repo.CreateAgent(ctx, name)
	}
	return s.repo.SetAgentHeartbeatAt(ctx, name)
}

func (s *service) GetAgents(ctx context.Context) ([]domain.Agent, error) {
	return s.repo.GetAgents(ctx)
}

func (s *service) SaveOperationDuration(ctx context.Context, name string, duration uint16) error {
	operationDuration, err := s.repo.GetOperationDuration(ctx, name)
	if err != nil {
		return err
	}
	if operationDuration.Name == "" {
		return s.repo.CreateOperationDuration(ctx, name, duration)
	}
	return s.repo.UpdateOperationDuration(ctx, name, duration)
}

func (s *service) GetOperationDurations(ctx context.Context) ([]domain.OperationDuration, error) {
	return s.repo.GetOperationDurations(ctx)
}

func (s *service) CreateSubExpression(ctx context.Context, se *domain.SubExpression) error {
	return s.repo.CreateSubExpression(ctx, se)
}
