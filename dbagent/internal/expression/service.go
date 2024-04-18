package expression

import (
	"context"
	"github.com/VadimGossip/calculator/dbagent/internal/domain"
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
	SaveOperationDuration(ctx context.Context, name string, duration uint32) error
	GetOperationDurations(ctx context.Context) ([]domain.OperationDuration, error)
	CreateSubExpression(ctx context.Context, se *domain.SubExpression) error
	StartSubExpressionEval(ctx context.Context, seId int64, agent string) (bool, error)
	StopSubExpressionEval(ctx context.Context, seId int64, result *float64, errMsg string) error
	GetSubExpressionIsLast(ctx context.Context, seId int64) (bool, error)
	GetReadySubExpressions(ctx context.Context, eId *int64, skipTimeoutSec uint32) ([]domain.SubExpression, error)
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

func (s *service) SaveOperationDuration(ctx context.Context, name string, duration uint32) error {
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
	e, err := s.repo.GetExpressionBySeId(ctx, seId)
	if err != nil {
		return false, err
	}
	if e.State == domain.ExpressionStateNew {
		e.State = domain.ExpressionStateInProgress
		if err = s.repo.UpdateExpression(ctx, *e); err != nil {
			return false, err
		}
	}
	return s.repo.StartSubExpressionEval(ctx, seId, agent)
}

func (s *service) StopSubExpressionEval(ctx context.Context, seId int64, result *float64, errMsg string) error {
	if err := s.repo.StopSubExpressionEval(ctx, seId, result); err != nil {
		return err
	}

	e, err := s.repo.GetExpressionSummaryBySeId(ctx, seId)
	if err != nil {
		return err
	}

	if result == nil {
		e.ErrorMsg = errMsg
		e.State = domain.ExpressionStateError
		return s.repo.UpdateExpression(ctx, e)
	}

	isLast, err := s.repo.GetSubExpressionIsLast(ctx, seId)
	if err != nil {
		return err
	}

	if isLast {
		e.State = domain.ExpressionStateOK
		return s.repo.UpdateExpression(ctx, e)
	}

	return nil
}

func (s *service) StopSubExpressionEval2(ctx context.Context, seId int64, result *float64) error {
	return s.repo.StopSubExpressionEval(ctx, seId, result)
}

func (s *service) GetSubExpressionIsLast(ctx context.Context, seId int64) (bool, error) {
	return s.repo.GetSubExpressionIsLast(ctx, seId)
}

func (s *service) GetReadySubExpressions(ctx context.Context, eId *int64, skipTimeoutSec uint32) ([]domain.SubExpression, error) {
	return s.repo.GetReadySubExpressions(ctx, eId, skipTimeoutSec)
}

func (s *service) SkipAgentSubExpressions(ctx context.Context, agent string) error {
	return s.repo.SkipAgentSubExpressions(ctx, agent)
}
