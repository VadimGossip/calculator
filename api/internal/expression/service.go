package expression

import (
	"context"
	"github.com/VadimGossip/calculator/api/internal/domain"
	"github.com/VadimGossip/calculator/api/internal/parser"
	"github.com/VadimGossip/calculator/api/internal/rabbitmq"
	"github.com/VadimGossip/calculator/api/internal/validation"
	"github.com/VadimGossip/calculator/api/internal/writer"
	"github.com/sirupsen/logrus"
	"time"
)

type service struct {
	cfg               domain.ExpressionCfg
	parseService      parser.Service
	validationService validation.Service
	writerService     writer.Service
	producer          rabbitmq.Producer
}

type Service interface {
	ValidateAndSimplify(value string) (string, error)
	RegisterExpression(ctx context.Context, value string) (int64, error)
	GetExpressions(ctx context.Context) ([]domain.Expression, error)
	SaveAgentHeartbeat(ctx context.Context, name string) error
	GetAgents(ctx context.Context) ([]domain.Agent, error)
	SaveOperationDurations(ctx context.Context, data map[string]uint16) error
	GetOperationDurations(ctx context.Context) ([]domain.OperationDuration, error)
	StartSubExpressionEval(ctx context.Context, seId int64, agent string) (bool, error)
	StopSubExpressionEval(ctx context.Context, seId int64, result *float64, errMsg string) error
	RunProcessWatchers(ctx context.Context)
}

var _ Service = (*service)(nil)

func NewService(cfg domain.ExpressionCfg, parseService parser.Service, validationService validation.Service, writerService writer.Service, producer rabbitmq.Producer) *service {
	return &service{cfg: cfg, parseService: parseService, validationService: validationService, writerService: writerService, producer: producer}
}

func (s *service) prepareSubExpressionQueryData(ctx context.Context, expressionId *int64) ([]domain.SubExpressionQueryItem, error) {
	seReady, err := s.writerService.GetReadySubExpressions(ctx, expressionId, time.Duration(s.cfg.HungTimeout)*time.Minute)
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
			Duration:  se.OperationDuration,
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

func (s *service) ValidateAndSimplify(value string) (string, error) {
	return s.validationService.ValidateAndSimplify(value)
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

func (s *service) StartSubExpressionEval(ctx context.Context, seId int64, agent string) (bool, error) {
	return s.writerService.StartSubExpressionEval(ctx, seId, agent)
}

func (s *service) StopSubExpressionEval(ctx context.Context, seId int64, result *float64, errMsg string) error {
	if err := s.writerService.StopSubExpressionEval(ctx, seId, result); err != nil {
		return err
	}
	e, err := s.writerService.GetExpressionSummaryBySeId(ctx, seId)
	if err != nil {
		return err
	}

	if result == nil {
		e.ErrorMsg = errMsg
		e.State = domain.ExpressionStateError
		return s.writerService.UpdateExpression(ctx, e)
	}

	isLast, err := s.writerService.GetSubExpressionIsLast(ctx, seId)
	if err != nil {
		return err
	}
	if isLast {
		e.State = domain.ExpressionStateOK
		return s.writerService.UpdateExpression(ctx, e)
	}
	return s.prepareAndPublish(ctx, &e.Id)
}

func (s *service) runHungProcessWatcher(ctx context.Context) {
	logrus.Info("HungProcessWatcher started")
	defer logrus.Info("HungProcessWatcher stopped")
	ticker := time.NewTicker(time.Duration(s.cfg.HungCheckPeriod) * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			logrus.Info("HungProcessWatcher. Checking hung subexpressions")
			if err := s.prepareAndPublish(ctx, nil); err != nil {
				logrus.Infof("HungProcessWatcher. prepare and publish err %s", err)
			}
		}
	}
}

func (s *service) checkAgents(ctx context.Context) error {
	agents, err := s.writerService.GetAgents(ctx)
	if err != nil {
		return err
	}
	for _, agent := range agents {
		tp := time.Since(agent.LastHeartbeatAt)
		if tp > time.Duration(s.cfg.AgentDownTimeout)*time.Minute {
			logrus.Infof("Agent %s Last HeartBeatAt %s TimePassed %s Allowed %s. All Agent Task will be skipped.", agent.Name, agent.LastHeartbeatAt, tp, 5*time.Minute)
			return s.writerService.SkipAgentSubExpressions(ctx, agent.Name)
		}
	}
	return nil
}

func (s *service) runAgentsWatcher(ctx context.Context) {
	logrus.Info("AgentsWatcher started")
	defer logrus.Info("AgentsWatcher stopped")
	ticker := time.NewTicker(time.Duration(s.cfg.AgentDownCheckPeriod) * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			logrus.Info("AgentsWatcher. Checking last agents heartbeat")
			if err := s.checkAgents(ctx); err != nil {
				logrus.Infof("AgentsWatcher. checkAgents err %s", err)
			}
		}
	}
}

func (s *service) RunProcessWatchers(ctx context.Context) {
	go s.runHungProcessWatcher(ctx)
	go s.runAgentsWatcher(ctx)
}
