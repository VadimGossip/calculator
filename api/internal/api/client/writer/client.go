package writer

import (
	"context"
	"fmt"
	"github.com/VadimGossip/calculator/api/internal/api/grpcservice/writergrpc"
	"github.com/VadimGossip/calculator/api/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"time"
)

type Client interface {
	Connect() error
	Disconnect() error
	GetReadySubExpressions(ctx context.Context, seId *int64, skipTimeoutSec uint32) ([]domain.ReadySubExpression, error)
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

func (c *client) mapSubExpressions(gses []*writergrpc.SubExpression) []domain.ReadySubExpression {
	readySes := make([]domain.ReadySubExpression, 0, len(gses))
	for _, gse := range gses {
		se := domain.ReadySubExpression{
			Id:        gse.SeId,
			Val1:      gse.Val1,
			Val2:      gse.Val2,
			Operation: gse.Operation,
			IsLast:    gse.IsLast,
		}
		readySes = append(readySes, se)
	}
	return readySes
}

func (c *client) GetReadySubExpressions(ctx context.Context, seId *int64, skipTimeoutSec uint32) ([]domain.ReadySubExpression, error) {
	var id int64
	if seId != nil {
		id = *seId
	}
	res, err := c.writerClient.GetReadySubExpressions(ctx, &writergrpc.ReadySubExpressionsRequest{
		SeId:           id,
		SeIsValid:      seId != nil,
		SkipTimeoutSec: skipTimeoutSec,
	})
	return c.mapSubExpressions(res.SubExpressions), err
}
