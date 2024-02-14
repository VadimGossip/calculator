package domain

type SubExpressionQueryItem struct {
	Id        int64   `json:"id"`
	Val1      float64 `json:"val1"`
	Val2      float64 `json:"val2"`
	Operation string  `json:"operation"`
	Duration  uint16  `json:"operation_duration"`
}
