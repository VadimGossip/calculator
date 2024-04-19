package writer

import (
	"context"
	"github.com/VadimGossip/calculator/dbagent/internal/domain"
	"github.com/VadimGossip/calculator/dbagent/internal/expression"
	"github.com/VadimGossip/calculator/dbagent/internal/token"
	"github.com/VadimGossip/calculator/dbagent/internal/user"
)

type service struct {
	exprService  expression.Service
	userService  user.Service
	tokenService token.Service
}

type Service interface {
	SaveAgentHeartbeat(ctx context.Context, name string) error
	StartSubExpressionEval(ctx context.Context, seId int64, agent string) (bool, error)
	StopSubExpressionEval(ctx context.Context, seId int64, result *float64, errMsg string) error
	GetReadySubExpressions(ctx context.Context, eId *int64, skipTimeoutSec uint32) ([]domain.SubExpression, error)
	GetExpressionByReqUid(ctx context.Context, reqUid string) (*domain.Expression, error)
	CreateExpression(ctx context.Context, e *domain.Expression) error
	CreateSubExpression(ctx context.Context, se *domain.SubExpression) error
	GetExpressions(ctx context.Context) ([]domain.Expression, error)
	GetAgents(ctx context.Context) ([]domain.Agent, error)
	SaveOperationDuration(ctx context.Context, name string, duration uint32) error
	GetOperationDurations(ctx context.Context) ([]domain.OperationDuration, error)
	SkipAgentSubExpressions(ctx context.Context, agent string) error
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByCredentials(ctx context.Context, login string, password string) (*domain.User, error)
	CreateToken(ctx context.Context, token *domain.Token) error
	GetToken(ctx context.Context, token string) (*domain.Token, error)
}

var _ Service = (*service)(nil)

func NewService(exprService expression.Service, userService user.Service, tokenService token.Service) *service {
	return &service{exprService: exprService, userService: userService, tokenService: tokenService}
}

func (s *service) SaveAgentHeartbeat(ctx context.Context, name string) error {
	return s.exprService.SaveAgentHeartbeat(ctx, name)
}

func (s *service) StartSubExpressionEval(ctx context.Context, seId int64, agent string) (bool, error) {
	return s.exprService.StartSubExpressionEval(ctx, seId, agent)
}

func (s *service) StopSubExpressionEval(ctx context.Context, seId int64, result *float64, errMsg string) error {
	return s.exprService.StopSubExpressionEval(ctx, seId, result, errMsg)
}

func (s *service) GetReadySubExpressions(ctx context.Context, eId *int64, skipTimeoutSec uint32) ([]domain.SubExpression, error) {
	return s.exprService.GetReadySubExpressions(ctx, eId, skipTimeoutSec)
}

func (s *service) GetExpressionByReqUid(ctx context.Context, reqUid string) (*domain.Expression, error) {
	return s.exprService.GetExpressionByReqUid(ctx, reqUid)
}

func (s *service) CreateExpression(ctx context.Context, e *domain.Expression) error {
	return s.exprService.CreateExpression(ctx, e)
}

func (s *service) CreateSubExpression(ctx context.Context, se *domain.SubExpression) error {
	return s.exprService.CreateSubExpression(ctx, se)
}

func (s *service) GetExpressions(ctx context.Context) ([]domain.Expression, error) {
	return s.exprService.GetExpressions(ctx)
}

func (s *service) GetAgents(ctx context.Context) ([]domain.Agent, error) {
	return s.exprService.GetAgents(ctx)
}

func (s *service) SaveOperationDuration(ctx context.Context, name string, duration uint32) error {
	return s.exprService.SaveOperationDuration(ctx, name, duration)
}

func (s *service) GetOperationDurations(ctx context.Context) ([]domain.OperationDuration, error) {
	return s.exprService.GetOperationDurations(ctx)
}

func (s *service) SkipAgentSubExpressions(ctx context.Context, agent string) error {
	return s.exprService.SkipAgentSubExpressions(ctx, agent)
}

func (s *service) CreateUser(ctx context.Context, user *domain.User) error {
	return s.userService.Create(ctx, user)
}

func (s *service) GetUserByCredentials(ctx context.Context, login, password string) (*domain.User, error) {
	return s.userService.GetByCredentials(ctx, login, password)
}

func (s *service) CreateToken(ctx context.Context, token *domain.Token) error {
	return s.tokenService.Create(ctx, token)
}

func (s *service) GetToken(ctx context.Context, token string) (*domain.Token, error) {
	return s.tokenService.Get(ctx, token)
}
