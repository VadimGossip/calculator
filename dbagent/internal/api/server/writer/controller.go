package writer

import (
	"context"
	"github.com/VadimGossip/calculator/dbagent/internal/domain"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/VadimGossip/calculator/dbagent/internal/api/grpcservice/writergrpc"
	"github.com/VadimGossip/calculator/dbagent/internal/writer"
)

type Controller interface {
	Heartbeat(ctx context.Context, req *writergrpc.HeartbeatRequest) (*emptypb.Empty, error)
	StartEval(ctx context.Context, req *writergrpc.StartEvalRequest) (*writergrpc.StartEvalResponse, error)
	StopEval(ctx context.Context, req *writergrpc.StopEvalRequest) (*emptypb.Empty, error)
	GetReadySubExpressions(ctx context.Context, req *writergrpc.ReadySubExpressionsRequest) (*writergrpc.ReadySubExpressionsResponse, error)
	GetExpressionByReqUid(ctx context.Context, req *writergrpc.ExpressionByReqUidRequest) (*writergrpc.Expression, error)
	//CreateExpression(ctx context.Context, req *writergrpc.CreateExpressionRequest) (*writergrpc.CreateExpressionResponse, error)
	//CreateSubExpression(ctx context.Context, req *writergrpc.CreateSubExpressionRequest) (*writergrpc.CreateSubExpressionResponse, error)
	//GetExpressions(ctx context.Context, _ *emptypb.Empty) (*writergrpc.GetExpressionsResponse, error)
	//GetAgents(ctx context.Context,req *writergrpc.CreateExpressionRequest) (*writergrpc.GetAgentsResponse, error)
	//SaveOperationDuration(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error)
	//GetOperationDurations(ctx context.Context, _ *emptypb.Empty) (*writergrpc.GetOperDurResponse, error)
}

type controller struct {
	writer writer.Service
}

var _ Controller = (*controller)(nil)

func NewController(writer writer.Service) *controller {
	return &controller{writer: writer}
}

func (c *controller) Heartbeat(ctx context.Context, req *writergrpc.HeartbeatRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, c.writer.SaveAgentHeartbeat(ctx, req.AgentName)
}

func (c *controller) StartEval(ctx context.Context, req *writergrpc.StartEvalRequest) (*writergrpc.StartEvalResponse, error) {
	res, err := c.writer.StartSubExpressionEval(ctx, req.SeId, req.Agent)
	return &writergrpc.StartEvalResponse{Success: res}, err
}

func (c *controller) StopEval(ctx context.Context, req *writergrpc.StopEvalRequest) (*emptypb.Empty, error) {
	var result *float64
	if req.Error == "" {
		result = &req.Result
	}
	return &emptypb.Empty{}, c.writer.StopSubExpressionEval(ctx, req.SeId, result, req.Error)
}

func (c *controller) mapSubExpression(se *domain.SubExpression) *writergrpc.SubExpression {
	return &writergrpc.SubExpression{
		SeId:              se.Id,
		Val1:              se.Val1,
		Val2:              se.Val2,
		Operation:         se.Operation,
		OperationDuration: se.OperationDuration,
		IsLast:            se.IsLast,
	}
}

func (c *controller) mapSubExpressions(subExpressions []domain.SubExpression) []*writergrpc.SubExpression {
	gses := make([]*writergrpc.SubExpression, 0)
	for _, se := range subExpressions {
		item := se
		gses = append(gses, c.mapSubExpression(&item))
	}
	return gses
}

func (c *controller) GetReadySubExpressions(ctx context.Context, req *writergrpc.ReadySubExpressionsRequest) (*writergrpc.ReadySubExpressionsResponse, error) {
	var eId *int64
	if req.SeIsValid {
		eId = &req.SeId
	}

	res, err := c.writer.GetReadySubExpressions(ctx, eId, req.SkipTimeoutSec)
	return &writergrpc.ReadySubExpressionsResponse{SubExpressions: c.mapSubExpressions(res)}, err
}

func (c *controller) mapExpression(e *domain.Expression) *writergrpc.Expression {
	var result float64
	if e.Result != nil {
		result = *e.Result
	}
	var evalStartedAt, evalFinishedAt int64
	if e.EvalStartedAt != nil {
		evalStartedAt = (*e.EvalStartedAt).Unix()
	}

	if e.EvalFinishedAt != nil {
		evalFinishedAt = (*e.EvalFinishedAt).Unix()
	}
	return &writergrpc.Expression{
		Id:             e.Id,
		ReqUid:         e.ReqUid,
		Value:          e.Value,
		Result:         result,
		State:          e.State,
		Error:          e.ErrorMsg,
		CreatedAt:      e.CreatedAt.Unix(),
		EvalStartedAt:  evalStartedAt,
		EvalFinishedAt: evalFinishedAt,
	}
}

func (c *controller) GetExpressionByReqUid(ctx context.Context, req *writergrpc.ExpressionByReqUidRequest) (*writergrpc.Expression, error) {
	e, err := c.writer.GetExpressionByReqUid(ctx, req.ReqUid)
	return c.mapExpression(e), err
}

//CreateSubExpression(ctx context.Context, req *writergrpc.CreateSubExpressionRequest) (*writergrpc.CreateSubExpressionResponse, error)
//GetExpressions(ctx context.Context, _ *emptypb.Empty) (*writergrpc.GetExpressionsResponse, error)
//GetAgents(ctx context.Context,req *writergrpc.CreateExpressionRequest) (*writergrpc.GetAgentsResponse, error)
//SaveOperationDuration(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error)
//GetOperationDurations(ctx context.Context, _ *emptypb.Empty) (*writergrpc.GetOperDurResponse, error)
