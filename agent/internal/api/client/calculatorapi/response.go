package calculatorapi

type StartSubExpressionEvalResponse struct {
	Error   string `json:"error,omitempty"`
	Status  int    `json:"status" example:"200"`
	Success bool   `json:"Success" example:"true"`
}
type CommonResponse struct {
	Error  string `json:"error,omitempty"`
	Status int    `json:"status" example:"200"`
}
