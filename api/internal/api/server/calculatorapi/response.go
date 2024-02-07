package calculatorapi

type CreateExpressionResponse struct {
	Id     int64  `json:"expression_id,omitempty" example:"1"`
	Error  string `json:"error,omitempty" example:"parse error"`
	Status int    `json:"status" example:"200"`
}
