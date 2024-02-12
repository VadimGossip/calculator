package manager

import (
	"context"
	"github.com/VadimGossip/calculator/api/internal/domain"
	"github.com/VadimGossip/calculator/api/internal/rabbitmq"
	"github.com/VadimGossip/calculator/api/internal/writer"
)

type service struct {
	writerService writer.Service
	producer      rabbitmq.Producer
}

type Service interface {
	RegisterExpression(ctx context.Context, value string) (int64, error)
	GetExpressions(ctx context.Context) ([]domain.Expression, error)
	SaveAgentHeartbeat(ctx context.Context, name string) error
	GetAgents(ctx context.Context) ([]domain.Agent, error)
	SaveOperationDurations(ctx context.Context, data map[string]uint16) error
	GetOperationDurations(ctx context.Context) ([]domain.OperationDuration, error)
}

var _ Service = (*service)(nil)

func NewService(writerService writer.Service, producer rabbitmq.Producer) *service {
	return &service{writerService: writerService, producer: producer}
}

func (s *service) RegisterExpression(ctx context.Context, value string) (int64, error) {
	expr := domain.Expression{Value: value}

	if err := s.writerService.CreateExpression(ctx, &expr); err != nil {
		return 0, err
	}
	a := 2
	b := 3
	se := domain.SubExpression{
		ExpressionsId: expr.Id,
		Val1:          &a,
		Val2:          &b,
		OperationName: "-",
		IsLast:        true,
	}

	if err := s.writerService.CreateSubExpression(ctx, &se); err != nil {
		return 0, err
	}

	//if err := s.writerService.SaveExpressionResult(ctx, expr.Id, 6); err != nil {
	//	return 0, err
	//}

	//if err := s.producer.SendMessage("", "text/plain", expr.Id); err != nil {
	//	return 0, err
	//}
	return expr.Id, nil
}

func (s *service) GetExpressions(ctx context.Context) ([]domain.Expression, error) {
	return s.writerService.GetExpressions(ctx)
}

func (s *service) SaveAgentHeartbeat(ctx context.Context, name string) error {
	return s.writerService.SaveAgentHeartbeat(ctx, name)
}

func (s *service) GetAgents(ctx context.Context) ([]domain.Agent, error) {
	return s.writerService.GetAgents(ctx)
}

func (s *service) SaveOperationDurations(ctx context.Context, data map[string]uint16) error {
	for key, value := range data {
		if err := s.writerService.SaveOperationDuration(ctx, key, value); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) GetOperationDurations(ctx context.Context) ([]domain.OperationDuration, error) {
	return s.writerService.GetOperationDurations(ctx)
}
