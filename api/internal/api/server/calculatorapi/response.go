package calculatorapi

import "github.com/VadimGossip/calculator/api/internal/domain"

type CreateExpressionResponse struct {
	Id     int64  `json:"expression_id,omitempty" example:"1"`
	Error  string `json:"error,omitempty" example:"parse error"`
	Status int    `json:"status" example:"200"`
}

type GetExpressionsResponse struct {
	Expressions []domain.Expression `json:"expressions"`
	Error       string              `json:"error,omitempty"`
	Status      int                 `json:"status" example:"200"`
}

type GetAgentsResponse struct {
	Agents []domain.Agent `json:"agents"`
	Error  string         `json:"error,omitempty"`
	Status int            `json:"status" example:"200"`
}

type GetOperationDurationsResponse struct {
	OperationDuration []domain.OperationDuration `json:"operation_durations"`
	Error             string                     `json:"error,omitempty"`
	Status            int                        `json:"status" example:"200"`
}
