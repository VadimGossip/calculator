package writer

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/VadimGossip/calculator/dbagent/internal/api/grpcservice/writergrpc"
	"github.com/VadimGossip/calculator/dbagent/internal/writer"
)

type Controller interface {
	StartEval(ctx context.Context, req *writergrpc.StartEvalRequest) (*writergrpc.StartEvalResponse, error)
	StopEval(ctx context.Context, req *writergrpc.StopEvalRequest) (*emptypb.Empty, error)
	Heartbeat(ctx context.Context, req *writergrpc.HeartbeatRequest) (*emptypb.Empty, error)
}

type controller struct {
	writer writer.Service
}

var _ Controller = (*controller)(nil)

func NewController(writer writer.Service) *controller {
	return &controller{writer: writer}
}

func (c *controller) StartEval(ctx context.Context, req *writergrpc.StartEvalRequest) (*writergrpc.StartEvalResponse, error) {
	res, err := c.writer.StartSubExpressionEval(ctx, req.SeId, req.Agent)
	return &writergrpc.StartEvalResponse{Success: res}, err
}

func (c *controller) StopEval(ctx context.Context, req *writergrpc.StopEvalRequest) (*emptypb.Empty, error) {
	result := new(float64)
	if req.Error == "" {
		result = &req.Result
	}
	return &emptypb.Empty{}, c.writer.StopSubExpressionEval(ctx, req.SeId, result, req.Error)
}

func (c *controller) Heartbeat(ctx context.Context, req *writergrpc.HeartbeatRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, c.writer.SaveAgentHeartbeat(ctx, req.AgentName)
}
