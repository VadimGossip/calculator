package worker

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/VadimGossip/calculator/agent/internal/api/client/calculatorapi"
	"github.com/VadimGossip/calculator/agent/internal/api/client/writer"
	"github.com/VadimGossip/calculator/agent/internal/domain"
)

type Service interface {
	Do(item domain.SubExpressionQueryItem) error
	GetMaxProcessAllowed() int
	RunHeartbeat(ctx context.Context)
}

type service struct {
	cfg                 domain.AgentCfg
	calculatorApiClient calculatorapi.ClientService
	writerClient        writer.Client
}

var _ Service = (*service)(nil)

func NewService(cfg domain.AgentCfg, calculatorClient calculatorapi.ClientService, writerClient writer.Client) *service {
	return &service{cfg: cfg, calculatorApiClient: calculatorClient, writerClient: writerClient}
}

func (s *service) eval(item domain.SubExpressionQueryItem) (*float64, error) {
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

func (s *service) Do(item domain.SubExpressionQueryItem) error {
	startResp, err := s.calculatorApiClient.SendStartEvalRequest(&calculatorapi.StartSubExpressionEvalRequest{
		Id:    item.Id,
		Agent: s.cfg.Name,
	})
	if startResp.Error != "" {
		logrus.Infof("received error on start eval attempt %s", startResp.Error)
		return errors.New(startResp.Error)
	}

	if !startResp.Success {
		logrus.Infof("failed to start eval attempt startResp.Success = false")
		return nil
	}

	result, err := s.eval(item)
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}

	stopResp, err := s.calculatorApiClient.SendStopEvalRequest(&calculatorapi.StopSubExpressionEvalRequest{
		Id:     item.Id,
		Result: result,
		Error:  errMsg,
	})

	if stopResp.Error != "" {
		return fmt.Errorf(stopResp.Error)
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
