package expression

import (
	"context"
	"github.com/VadimGossip/calculator/api/internal/api/client/writer"
	"github.com/VadimGossip/calculator/api/internal/domain"
	"github.com/VadimGossip/calculator/api/internal/parser"
	"github.com/VadimGossip/calculator/api/internal/rabbitmq"
	"github.com/VadimGossip/calculator/api/internal/validation"
	"github.com/sirupsen/logrus"
	"time"
)

type service struct {
	cfg               domain.ExpressionCfg
	parseService      parser.Service
	validationService validation.Service
	writerClient      writer.Client
	producer          rabbitmq.Producer
}

type Service interface {
	ValidateAndSimplify(value string) (string, error)
	RegisterExpression(ctx context.Context, e *domain.Expression) error
	GetExpressions(ctx context.Context) ([]domain.Expression, error)
	GetAgents(ctx context.Context) ([]domain.Agent, error)
	SaveOperationDurations(ctx context.Context, data map[string]uint16) error
	GetOperationDurations(ctx context.Context) ([]domain.OperationDuration, error)
	RunProcessWatchers(ctx context.Context)
}

var _ Service = (*service)(nil)

func NewService(cfg domain.ExpressionCfg, parseService parser.Service, validationService validation.Service, writerClient writer.Client, producer rabbitmq.Producer) *service {
	return &service{cfg: cfg, parseService: parseService, validationService: validationService, writerClient: writerClient, producer: producer}
}

func (s *service) prepareSubExpressionQueryData(ctx context.Context, eId *int64) ([]domain.ReadySubExpression, error) {
	return s.writerClient.GetReadySubExpressions(ctx, eId, uint32(s.cfg.HungTimeout))
}

func (s *service) publishSubExpressionQueryData(readySe []domain.ReadySubExpression) error {
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

func (s *service) RegisterExpression(ctx context.Context, e *domain.Expression) error {
	existing, err := s.writerService.GetExpressionByReqUid(ctx, e.ReqUid)
	if err != nil {
		return err
	}
	if existing != nil {
		*e = *existing
		return nil
	}

	if err = s.writerService.CreateExpression(ctx, e); err != nil {
		return err
	}

	if e.ErrorMsg == "" {
		idDict := make(map[int64]int64)
		for _, se := range s.parseService.ParseExpression(*e) {
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
			if err = s.writerService.CreateSubExpression(ctx, &enrichedSe); err != nil {
				return err
			}
			idDict[se.Id] = enrichedSe.Id
		}
		if err = s.prepareAndPublish(ctx, &e.Id); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) GetExpressions(ctx context.Context) ([]domain.Expression, error) {
	return s.writerService.GetExpressions(ctx)
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
			logrus.Infof("Agent %s Last HeartBeatAt %s TimePassed %s Allowed %s. All Agent Task will be skipped.", agent.Name, agent.LastHeartbeatAt, tp, time.Duration(s.cfg.AgentDownTimeout)*time.Minute)
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
