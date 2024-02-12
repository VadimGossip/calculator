package domain

type SubExpression struct {
	Id               int64
	ExpressionsId    int64
	Val1             *int
	Val2             *int
	SubExpressionId1 *int64
	SubExpressionId2 *int64
	OperationName    string
	Result           *int
	AgentId          int64
	IsLast           bool
}
