package orchestrator

import (
	"context"
	"github.com/VadimGossip/calculator/api/internal/rabbitmq"
	"github.com/VadimGossip/calculator/api/internal/writer"
)

type service struct {
	writerService writer.Service
	producer      rabbitmq.Producer
}

type Service interface {
	RegisterExpression(ctx context.Context, value string) (int64, error)
}

var _ Service = (*service)(nil)

func NewService(writerService writer.Service, producer rabbitmq.Producer) *service {
	return &service{writerService: writerService, producer: producer}
}

func (s *service) RegisterExpression(ctx context.Context, value string) (int64, error) {
	expr := writer.Expression{Value: value}

	if err := s.writerService.CreateExpression(ctx, &expr); err != nil {
		return 0, err
	}
	if err := s.producer.SendMessage("", "text/plain", expr.Id); err != nil {
		return 0, err
	}
	return expr.Id, nil
}
