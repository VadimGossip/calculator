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
	ReqUid         string     `json:"req_uid"`
	Value          string     `json:"expression_value"`
	Result         *float64   `json:"result"`
	State          string     `json:"status"`
	ErrorMsg       string     `json:"error_msg"`
	CreatedAt      time.Time  `json:"created_at"`
	EvalStartedAt  *time.Time `json:"eval_started_at"`
	EvalFinishedAt *time.Time `json:"eval_finished_at"`
}
