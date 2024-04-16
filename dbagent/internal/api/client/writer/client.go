package writer

import (
	"context"
	"fmt"
	"github.com/VadimGossip/calculator/dbagent/internal/api/grpcservice/writergrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"time"
)

type Client interface {
	Connect() error
	Disconnect() error
	StartEval(ctx context.Context, seId int64, agentName string) (*writergrpc.StartEvalResponse, error)
	StopEval(ctx context.Context, seId int64, result *float64, errMsg string) error
	Heartbeat(ctx context.Context, agentName string) error
}

type client struct {
	writerClient writergrpc.WriterServiceClient
	conn         *grpc.ClientConn
	host         string
	port         uint32
}

var _ Client = (*client)(nil)

func NewClient(host string, port uint32) *client {
	return &client{host: host, port: port}
}

func (c *client) Connect() error {
	kap := keepalive.ClientParameters{
		Time:                5 * time.Minute,
		Timeout:             10 * time.Second,
		PermitWithoutStream: true,
	}

	dialOptions := []grpc.DialOption{
		grpc.WithKeepaliveParams(kap),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	uri := fmt.Sprintf("%s:%d", c.host, c.port)
	conn, err := grpc.Dial(uri, dialOptions...)
	if err != nil {
		return err
	}
	c.conn = conn
	c.writerClient = writergrpc.NewWriterServiceClient(c.conn)
	return nil
}

func (c *client) Disconnect() error {
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("error while disconection grps client %s", err)
	}
	return nil
}

func (c *client) StartEval(ctx context.Context, seId int64, agentName string) (*writergrpc.StartEvalResponse, error) {
	return c.writerClient.StartEval(ctx, &writergrpc.StartEvalRequest{
		SeId:  seId,
		Agent: agentName,
	})
}

func (c *client) StopEval(ctx context.Context, seId int64, result *float64, errMsg string) error {
	req := &writergrpc.StopEvalRequest{
		SeId:  seId,
		Error: errMsg,
	}

	if result != nil {
		req.Result = *result
	}

	_, err := c.writerClient.StopEval(ctx, req)
	return err
}

func (c *client) Heartbeat(ctx context.Context, agentName string) error {
	_, err := c.writerClient.Heartbeat(ctx, &writergrpc.HeartbeatRequest{AgentName: agentName})
	return err
}
