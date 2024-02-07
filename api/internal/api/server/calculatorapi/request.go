package calculatorapi

type CreateExpressionRequest struct {
	ExpressionValue string `json:"expression" binding:"required" example:"2+2+2"`
}
