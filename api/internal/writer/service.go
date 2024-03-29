package writer

import (
	"context"
	"github.com/VadimGossip/calculator/api/internal/domain"
	"time"
)

type service struct {
	repo Repository
}

type Service interface {
	CreateExpression(ctx context.Context, e *domain.Expression) error
	UpdateExpression(ctx context.Context, e domain.Expression) error
	GetExpressionSummaryBySeId(ctx context.Context, seId int64) (domain.Expression, error)
	GetExpressionBySeId(ctx context.Context, seId int64) (*domain.Expression, error)
	GetExpressionByReqUid(ctx context.Context, reqUid string) (*domain.Expression, error)
	GetExpressions(ctx context.Context) ([]domain.Expression, error)
	SaveAgentHeartbeat(ctx context.Context, name string) error
	GetAgents(ctx context.Context) ([]domain.Agent, error)
	SaveOperationDuration(ctx context.Context, name string, duration uint16) error
	GetOperationDurations(ctx context.Context) ([]domain.OperationDuration, error)
	CreateSubExpression(ctx context.Context, s *domain.SubExpression) error
	StartSubExpressionEval(ctx context.Context, seId int64, agent string) (bool, error)
	StopSubExpressionEval(ctx context.Context, seId int64, result *float64) error
	GetSubExpressionIsLast(ctx context.Context, seId int64) (bool, error)
	GetReadySubExpressions(ctx context.Context, expressionId *int64, skipTimeout time.Duration) ([]domain.SubExpression, error)
	SkipAgentSubExpressions(ctx context.Context, agent string) error
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

func (s *service) GetExpressionBySeId(ctx context.Context, seId int64) (*domain.Expression, error) {
	return s.repo.GetExpressionBySeId(ctx, seId)
}

func (s *service) GetExpressionByReqUid(ctx context.Context, reqUid string) (*domain.Expression, error) {
	return s.repo.GetExpressionByReqUid(ctx, reqUid)
}

func (s *service) GetExpressions(ctx context.Context) ([]domain.Expression, error) {
	return s.repo.GetExpressions(ctx)
}

func (s *service) SaveAgentHeartbeat(ctx context.Context, name string) error {
	success, err := s.repo.SetAgentHeartbeatAt(ctx, name)
	if err != nil {
		return err
	}
	if !success {
		return s.repo.CreateAgent(ctx, name)
	}
	return nil
}

func (s *service) GetAgents(ctx context.Context) ([]domain.Agent, error) {
	return s.repo.GetAgents(ctx)
}

func (s *service) SaveOperationDuration(ctx context.Context, name string, duration uint16) error {
	success, err := s.repo.UpdateOperationDuration(ctx, name, duration)
	if err != nil {
		return err
	}
	if !success {
		return s.repo.CreateOperationDuration(ctx, name, duration)
	}
	return nil
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

func (s *service) StopSubExpressionEval(ctx context.Context, seId int64, result *float64) error {
	return s.repo.StopSubExpressionEval(ctx, seId, result)
}

func (s *service) GetSubExpressionIsLast(ctx context.Context, seId int64) (bool, error) {
	return s.repo.GetSubExpressionIsLast(ctx, seId)
}

func (s *service) GetReadySubExpressions(ctx context.Context, expressionId *int64, skipTimeout time.Duration) ([]domain.SubExpression, error) {
	return s.repo.GetReadySubExpressions(ctx, expressionId, skipTimeout)
}

func (s *service) SkipAgentSubExpressions(ctx context.Context, agent string) error {
	return s.repo.SkipAgentSubExpressions(ctx, agent)
}
