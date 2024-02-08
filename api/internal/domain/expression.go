package domain

import "time"

const (
	ExpressionStateNew        = "New"
	ExpressionStateError      = "Error"
	ExpressionStateInProgress = "In Progress"
	ExpressionStateOK         = "Ok"
)

type Expression struct {
	Id             int64      `json:"expression_id"`
	Value          string     `json:"expression_value"`
	Result         *int       `json:"result"`
	State          string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	EvalStartedAt  *time.Time `json:"eval_started_at"`
	EvalFinishedAt *time.Time `json:"eval_finished_at"`
}

type QueueExpression struct {
	Id    int64  `json:"expression_id"`
	Value string `json:"expression_value"`
}
