package calculatorapi

type StartSubExpressionEvalRequest struct {
	Id    int64  `json:"sub_expression_id" binding:"required" example:"1"`
	Agent string `json:"agent_name" binding:"required" example:"agent1"`
}

type StopSubExpressionEvalRequest struct {
	Id     int64    `json:"sub_expression_id" binding:"required" example:"1"`
	Result *float64 `json:"result"`
	Error  string   `json:"calculation_error"`
}
