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
	GetReadyToEvalSubExpressions(ctx context.Context, req *writergrpc.ReadySubExpressionsRequest) (*writergrpc.ReadySubExpressionsResponse, error)
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

func (c *controller) mapSubExpressions(subExpressions []domain.SubExpression) []*writergrpc.SubExpression {
	gses := make([]*writergrpc.SubExpression, 0)
	for _, se := range subExpressions {
		gse := &writergrpc.SubExpression{
			SeId:              se.Id,
			Val1:              se.Val1,
			Val2:              se.Val2,
			Operation:         se.Operation,
			OperationDuration: se.OperationDuration,
			IsLast:            se.IsLast,
		}
		gses = append(gses, gse)
	}
	return gses
}

func (c *controller) GetReadyToEvalSubExpressions(ctx context.Context, req *writergrpc.ReadySubExpressionsRequest) (*writergrpc.ReadySubExpressionsResponse, error) {
	var eId *int64
	if req.SeIsValid {
		eId = &req.SeId
	}

	res, err := c.writer.GetReadySubExpressions(ctx, eId, req.SkipTimeoutSec)
	return &writergrpc.ReadySubExpressionsResponse{SubExpressions: c.mapSubExpressions(res)}, err
}
