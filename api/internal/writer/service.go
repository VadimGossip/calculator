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
	UpdateExpression(ctx context.Context, e domain.Expression) error
	GetExpressionSummaryBySeId(ctx context.Context, seId int64) (domain.Expression, error)
	GetExpressions(ctx context.Context) ([]domain.Expression, error)
	SaveAgentHeartbeat(ctx context.Context, name string) error
	GetAgents(ctx context.Context) ([]domain.Agent, error)
	SaveOperationDuration(ctx context.Context, name string, duration uint16) error
	GetOperationDurations(ctx context.Context) ([]domain.OperationDuration, error)
	CreateSubExpression(ctx context.Context, s *domain.SubExpression) error
	StartSubExpressionEval(ctx context.Context, seId int64, agent string) (bool, error)
	StopSubExpressionEval(ctx context.Context, seId int64, result float64) error
	GetReadySubExpressions(ctx context.Context, expressionId *int64) ([]domain.SubExpression, error)
}

var _ Service = (*service)(nil)

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) CreateExpression(ctx context.Context, e *domain.Expression) error {
	return s.repo.CreateExpression(ctx, e)
}

func (s *service) UpdateExpression(ctx context.Context, e domain.Expression) error {
	return s.repo.UpdateExpression(ctx, e)
}
func (s *service) GetExpressionSummaryBySeId(ctx context.Context, seId int64) (domain.Expression, error) {
	return s.repo.GetExpressionSummaryBySeId(ctx, seId)
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

func (s *service) StartSubExpressionEval(ctx context.Context, seId int64, agent string) (bool, error) {
	return s.repo.StartSubExpressionEval(ctx, seId, agent)
}

func (s *service) StopSubExpressionEval(ctx context.Context, seId int64, result float64) error {
	return s.repo.StopSubExpressionEval(ctx, seId, result)
}

func (s *service) GetReadySubExpressions(ctx context.Context, expressionId *int64) ([]domain.SubExpression, error) {
	return s.repo.GetReadySubExpressions(ctx, expressionId)
}
