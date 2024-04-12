package domain

import "time"

type SubExpression struct {
	Id                int64
	ExpressionId      int64
	Val1              *float64
	Val2              *float64
	SubExpressionId1  *int64
	SubExpressionId2  *int64
	Operation         string
	OperationDuration uint16
	Result            *float64
	Agent             string
	IsLast            bool
	EvalStartedAt     time.Time
}

type SubExpressionQueryItem struct {
	Id        int64   `json:"id"`
	Val1      float64 `json:"val1"`
	Val2      float64 `json:"val2"`
	Operation string  `json:"operation"`
	Duration  uint16  `json:"operation_duration"`
}
