package domain

type SubExpression struct {
	Id               int64
	ExpressionsId    int64
	Val1             *float64
	Val2             *float64
	SubExpressionId1 *int64
	SubExpressionId2 *int64
	Operation        string
	Result           *float64
	AgentId          int64
	IsLast           bool
}
