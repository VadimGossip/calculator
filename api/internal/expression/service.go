package expression

import (
	"context"
	"github.com/VadimGossip/calculator/api/internal/domain"
	"github.com/VadimGossip/calculator/api/internal/parser"
	"github.com/VadimGossip/calculator/api/internal/rabbitmq"
	"github.com/VadimGossip/calculator/api/internal/writer"
)

type service struct {
	parseService  parser.Service
	writerService writer.Service
	producer      rabbitmq.Producer
	events        chan *int64
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

func NewService(parseService parser.Service, writerService writer.Service, producer rabbitmq.Producer) *service {
	return &service{parseService: parseService, writerService: writerService, producer: producer}
}

func (s *service) prepareSubExpressionQueryData(ctx context.Context, expressionId *int64) ([]domain.SubExpressionQueryItem, error) {
	seReady, err := s.writerService.GetReadySubExpressions(ctx, expressionId)
	if err != nil {
		return nil, err
	}
	result := make([]domain.SubExpressionQueryItem, 0, len(seReady))
	for _, se := range seReady {
		item := domain.SubExpressionQueryItem{
			Id:        se.Id,
			Val1:      *se.Val1,
			Val2:      *se.Val2,
			Operation: se.Operation,
			Duration:  1000, // todo: read durations data
		}
		result = append(result, item)
	}
	return result, nil
}

func (s *service) publishSubExpressionQueryData(readySe []domain.SubExpressionQueryItem) error {
	for _, item := range readySe {
		if err := s.producer.SendMessage("", "application/json", item); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) prepareAndPublish(ctx context.Context, expressionId *int64) error {
	readySe, err := s.prepareSubExpressionQueryData(ctx, expressionId)
	if err != nil {
		return err
	}
	return s.publishSubExpressionQueryData(readySe)
}

func (s *service) RegisterExpression(ctx context.Context, value string) (int64, error) {
	e := domain.Expression{Value: value}

	if err := s.writerService.CreateExpression(ctx, &e); err != nil {
		return 0, err
	}

	idDict := make(map[int64]int64)
	for _, se := range s.parseService.ParseExpression(e) {
		enrichedSe := se
		if se.SubExpressionId1 != nil {
			if val, ok := idDict[*se.SubExpressionId1]; ok {
				enrichedSe.SubExpressionId1 = &val
			}
		}
		if se.SubExpressionId2 != nil {
			if val, ok := idDict[*se.SubExpressionId2]; ok {
				enrichedSe.SubExpressionId2 = &val
			}
		}
		if err := s.writerService.CreateSubExpression(ctx, &enrichedSe); err != nil {
			return 0, err
		}
		idDict[se.Id] = enrichedSe.Id
	}
	if err := s.prepareAndPublish(ctx, &e.Id); err != nil {
		return 0, err
	}

	return e.Id, nil
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
