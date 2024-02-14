package worker

import (
	"fmt"
	"github.com/VadimGossip/calculator/agent/internal/api/client/calculatorapi"
	"github.com/VadimGossip/calculator/agent/internal/domain"
	"time"
)

type Service interface {
}

type service struct {
	calculatorClient calculatorapi.ClientService
}

var _ Service = (*service)(nil)

func NewService(calculatorClient calculatorapi.ClientService) *service {
	return &service{calculatorClient: calculatorClient}
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
	startResp, err := s.calculatorClient.SendStartEvalRequest(&calculatorapi.StartSubExpressionEvalRequest{
		Id:    item.Id,
		Agent: "x",
	})

	if startResp.Skip == true {
		return nil
	}

	result, err := s.eval(item)
	if err != nil {
		return err
	}

	stopResp, err := s.calculatorClient.SendStopEvalRequest(&calculatorapi.StopSubExpressionEvalRequest{
		Id:     item.Id,
		Result: result,
		Error:  err.Error(),
	})

	if stopResp.Error != "" {
		return fmt.Errorf(stopResp.Error)
	}
	return nil
}
