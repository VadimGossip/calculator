package worker

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/VadimGossip/calculator/agent/internal/api/client/writer"
	"github.com/VadimGossip/calculator/agent/internal/domain"
)

type Service interface {
	Do(ctx context.Context, item domain.ReadySubExpression) error
	GetMaxProcessAllowed() int
	RunHeartbeat(ctx context.Context)
}

type service struct {
	cfg          domain.AgentCfg
	writerClient writer.Client
}

var _ Service = (*service)(nil)

func NewService(cfg domain.AgentCfg, writerClient writer.Client) *service {
	return &service{cfg: cfg, writerClient: writerClient}
}

func (s *service) eval(item domain.ReadySubExpression) (*float64, error) {
	var result float64
	switch item.Operation {
	case "+":
		time.Sleep(time.Duration(item.Duration) * time.Millisecond)
		result = item.Val1 + item.Val2
		return &result, nil
	case "-":
		time.Sleep(time.Duration(item.Duration) * time.Millisecond)
		result = item.Val1 - item.Val2
		return &result, nil
	case "*":
		time.Sleep(time.Duration(item.Duration) * time.Millisecond)
		result = item.Val1 * item.Val2
		return &result, nil
	case "/":
		time.Sleep(time.Duration(item.Duration) * time.Millisecond)
		if item.Val2 == 0 {
			return nil, fmt.Errorf("division on zero")
		}
		result = item.Val1 / item.Val2
		return &result, nil
	}
	return nil, fmt.Errorf("unknown operation")
}

func (s *service) Do(ctx context.Context, item domain.ReadySubExpression) error {
	startResp, err := s.writerClient.StartEval(ctx, item.Id, s.cfg.Name)
	if err != nil {
		logrus.Errorf("Received error on start eval %s", err)
		return err
	}

	if !startResp.Success {
		logrus.Infof("Failed to start eval attempt startResp.Success = false")
		return nil
	}

	result, err := s.eval(item)
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}

	if err = s.writerClient.StopEval(ctx, item.Id, result, errMsg); err != nil {
		logrus.Errorf("Received error on stop eval %s", err)
		return err
	}
	return nil
}

func (s *service) GetMaxProcessAllowed() int {
	return s.cfg.MaxProcesses
}

func (s *service) RunHeartbeat(ctx context.Context) {
	logrus.Info("Heartbeat started")
	defer logrus.Info("Heartbeat stopped")
	ticker := time.NewTicker(time.Duration(s.cfg.HeartbeatTimeout) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.writerClient.Heartbeat(ctx, s.cfg.Name); err != nil {
				logrus.Infof("Hearbeat error %s", err)
			}
		}
	}
}
