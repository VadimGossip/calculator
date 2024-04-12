package writer

import (
	"context"
	"github.com/VadimGossip/calculator/dbagent/internal/api/grpcservice/writergrpc"
	"github.com/VadimGossip/calculator/dbagent/internal/writer"
)

type Controller interface {
	Heartbeat(ctx context.Context, req *writergrpc.HeartbeatRequest) (*writergrpc.HeartbeatResponse, error)
}

type controller struct {
	writer writer.Service
}

var _ Controller = (*controller)(nil)

func NewController(writer writer.Service) *controller {
	return &controller{writer: writer}
}
func (c *controller) Heartbeat(ctx context.Context, req *writergrpc.HeartbeatRequest) (*writergrpc.HeartbeatResponse, error) {
	return &writergrpc.HeartbeatResponse{Error: c.writer.SaveAgentHeartbeat(ctx, req.AgentName).Error()}, nil
}
