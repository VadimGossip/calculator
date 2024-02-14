package calculatorapi

type StartSubExpressionEvalResponse struct {
	Error  string `json:"error,omitempty"`
	Status int    `json:"status" example:"200"`
	Skip   bool   `json:"skip" example:"false"`
}
type CommonResponse struct {
	Error  string `json:"error,omitempty"`
	Status int    `json:"status" example:"200"`
}
