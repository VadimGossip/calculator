package calculatorapi

import "github.com/VadimGossip/calculator/api/internal/domain"

type CreateExpressionResponse struct {
	Expression *domain.Expression `json:"expression,omitempty"`
	Error      string             `json:"error,omitempty" example:"parse error"`
	Status     int                `json:"status" example:"200"`
}

type GetExpressionsResponse struct {
	Expressions []domain.Expression `json:"expressions"`
	Error       string              `json:"error,omitempty"`
	Status      int                 `json:"status" example:"200"`
}

type CommonResponse struct {
	Error  string `json:"error,omitempty"`
	Status int    `json:"status" example:"200"`
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

type StartSubExpressionEvalResponse struct {
	Error   string `json:"error,omitempty"`
	Status  int    `json:"status" example:"200"`
	Success bool   `json:"success" example:"false"`
}

type RegisterUserResponse struct {
	Id     int64  `json:"id,omitempty"`
	Error  string `json:"error,omitempty"`
	Status int    `json:"status" example:"200"`
}
