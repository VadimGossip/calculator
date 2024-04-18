package writer

import (
	"context"
	"fmt"
	"github.com/VadimGossip/calculator/api/internal/api/grpcservice/writergrpc"
	"github.com/VadimGossip/calculator/api/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type Client interface {
	Connect() error
	Disconnect() error
	Heartbeat(ctx context.Context, agentName string) error
	StartEval(ctx context.Context, seId int64, agentName string) (bool, error)
	StopEval(ctx context.Context, seId int64, result *float64, errMsg string) error
	GetReadySubExpressions(ctx context.Context, eId *int64, skipTimeoutSec uint32) ([]domain.ReadySubExpression, error)
	GetExpressionByReqUid(ctx context.Context, reqUid string) (*domain.Expression, error)
	CreateExpression(ctx context.Context, e *domain.Expression) (int64, error)
	CreateSubExpression(ctx context.Context, se *domain.SubExpression) (int64, error)
	GetExpressions(ctx context.Context) ([]domain.Expression, error)
	GetAgents(ctx context.Context) ([]domain.Agent, error)
	SaveOperationDuration(ctx context.Context, name string, duration uint32) error
	GetOperationDurations(ctx context.Context) ([]domain.OperationDuration, error)
	SkipAgentSubExpressions(ctx context.Context, agentName string) error
}

type client struct {
	writerClient writergrpc.WriterServiceClient
	conn         *grpc.ClientConn
	host         string
	port         int
}

var _ Client = (*client)(nil)

func NewClient(host string, port int) *client {
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

func (c *client) Heartbeat(ctx context.Context, agentName string) error {
	_, err := c.writerClient.Heartbeat(ctx, &writergrpc.HeartbeatRequest{AgentName: agentName})
	return err
}

func (c *client) StartEval(ctx context.Context, seId int64, agentName string) (bool, error) {
	response, err := c.writerClient.StartEval(ctx, &writergrpc.StartEvalRequest{
		SeId:  seId,
		Agent: agentName,
	})
	if err != nil {
		return false, err
	}

	return response.Success, nil
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

func (c *client) unwrapSubExpression(gse *writergrpc.SubExpression) *domain.SubExpression {
	var subExpression1, subExpression2 *int64
	var val1, val2 *float64
	if gse.SubExpressionId1 != -1 {
		subExpression1 = &gse.SubExpressionId1
	} else {
		val1 = &gse.Val1
	}

	if gse.SubExpressionId2 != -1 {
		subExpression2 = &gse.SubExpressionId2
	} else {
		val2 = &gse.Val2
	}

	return &domain.SubExpression{
		ExpressionId:      gse.ExpressionId,
		Val1:              val1,
		Val2:              val2,
		SubExpressionId1:  subExpression1,
		SubExpressionId2:  subExpression2,
		Operation:         gse.Operation,
		OperationDuration: gse.OperationDuration,
		IsLast:            gse.IsLast,
	}
}

func (c *client) GetReadySubExpressions(ctx context.Context, eId *int64, skipTimeoutSec uint32) ([]domain.ReadySubExpression, error) {
	var id int64 = -1
	if eId != nil {
		id = *eId
	}

	req := &writergrpc.ReadySubExpressionsRequest{
		EId:            id,
		SkipTimeoutSec: skipTimeoutSec,
	}

	response, err := c.writerClient.GetReadySubExpressions(ctx, req)
	if err != nil {
		return nil, err
	}

	readySes := make([]domain.ReadySubExpression, len(response.SubExpressions))
	for i, gse := range response.SubExpressions {
		readySes[i] = domain.ReadySubExpression{
			Id:        gse.Id,
			Val1:      gse.Val1,
			Val2:      gse.Val2,
			Operation: gse.Operation,
			Duration:  gse.OperationDuration,
			IsLast:    gse.IsLast,
		}
	}

	return readySes, err
}

func (c *client) wrapExpression(e *domain.Expression) *writergrpc.Expression {
	if e == nil {
		return &writergrpc.Expression{}
	}
	return &writergrpc.Expression{
		ReqUid: e.ReqUid,
		Value:  e.Value,
		State:  e.State,
	}
}

func (c *client) unwrapExpression(ge *writergrpc.Expression) *domain.Expression {
	if ge.Id == -1 {
		return nil
	}

	var evalStartedAt, evalFinishedAt *time.Time
	if ge.EvalStartedAt != 0 {
		evalStartedAtUnix := time.Unix(ge.EvalStartedAt, 0)
		evalStartedAt = &evalStartedAtUnix
	}

	if ge.EvalFinishedAt != 0 {
		evalFinishedAtUnix := time.Unix(ge.EvalFinishedAt, 0)
		evalFinishedAt = &evalFinishedAtUnix
	}

	var result *float64
	if ge.State == domain.ExpressionStateOK {
		result = &ge.Result
	}

	return &domain.Expression{
		Id:             ge.Id,
		ReqUid:         ge.ReqUid,
		Value:          ge.Value,
		Result:         result,
		State:          ge.State,
		ErrorMsg:       ge.Error,
		CreatedAt:      time.Unix(ge.CreatedAt, 0),
		EvalStartedAt:  evalStartedAt,
		EvalFinishedAt: evalFinishedAt,
	}
}

func (c *client) wrapSubExpression(se *domain.SubExpression) *writergrpc.SubExpression {
	if se == nil {
		return &writergrpc.SubExpression{}
	}

	var subExpression1, subExpression2 int64 = -1, -1
	var val1, val2 float64
	if se.SubExpressionId1 != nil {
		subExpression1 = *se.SubExpressionId1
	} else {
		val1 = *se.Val1
	}
	if se.SubExpressionId2 != nil {
		subExpression2 = *se.SubExpressionId2
	} else {
		val2 = *se.Val2
	}

	return &writergrpc.SubExpression{
		ExpressionId:      se.ExpressionId,
		Val1:              val1,
		Val2:              val2,
		SubExpressionId1:  subExpression1,
		SubExpressionId2:  subExpression2,
		Operation:         se.Operation,
		OperationDuration: se.OperationDuration,
		IsLast:            se.IsLast,
	}
}

func (c *client) unwrapAgent(ga *writergrpc.Agent) *domain.Agent {
	return &domain.Agent{
		Name:            ga.Name,
		CreatedAt:       time.Unix(ga.CreatedAt, 0),
		LastHeartbeatAt: time.Unix(ga.LastHbAt, 0),
	}
}

func (c *client) unwrapOperationDuration(gd *writergrpc.OperationDuration) *domain.OperationDuration {
	return &domain.OperationDuration{
		Name:      gd.Name,
		Duration:  gd.Duration,
		CreatedAt: time.Unix(gd.CreatedAt, 0),
		UpdatedAt: time.Unix(gd.UpdatedAt, 0),
	}
}

func (c *client) GetExpressionByReqUid(ctx context.Context, reqUid string) (*domain.Expression, error) {
	req := &writergrpc.ExpressionByReqUidRequest{ReqUid: reqUid}

	response, err := c.writerClient.GetExpressionByReqUid(ctx, req)
	if err != nil {
		return nil, err
	}
	return c.unwrapExpression(response), nil

}

func (c *client) CreateExpression(ctx context.Context, e *domain.Expression) (int64, error) {
	req := &writergrpc.CreateExpressionRequest{E: c.wrapExpression(e)}

	response, err := c.writerClient.CreateExpression(ctx, req)
	if err != nil {
		return 0, err
	}
	return response.Id, err
}

func (c *client) CreateSubExpression(ctx context.Context, se *domain.SubExpression) (int64, error) {
	req := &writergrpc.CreateSubExpressionRequest{Se: c.wrapSubExpression(se)}

	response, err := c.writerClient.CreateSubExpression(ctx, req)
	if err != nil {
		return 0, err
	}
	return response.Id, err
}

func (c *client) GetExpressions(ctx context.Context) ([]domain.Expression, error) {
	response, err := c.writerClient.GetExpressions(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	expressions := make([]domain.Expression, len(response.Expressions))
	for i, ge := range response.Expressions {
		expressions[i] = *c.unwrapExpression(ge)
	}

	return expressions, err
}

func (c *client) GetAgents(ctx context.Context) ([]domain.Agent, error) {
	response, err := c.writerClient.GetAgents(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	agents := make([]domain.Agent, len(response.Agents))
	for i, ga := range response.Agents {
		fmt.Println(c.unwrapAgent(ga))
		agents[i] = *c.unwrapAgent(ga)
	}
	return agents, err
}

func (c *client) SaveOperationDuration(ctx context.Context, name string, duration uint32) error {
	req := &writergrpc.CreateOperDurRequest{Name: name, Duration: duration}
	_, err := c.writerClient.SaveOperationDuration(ctx, req)
	return err
}

func (c *client) GetOperationDurations(ctx context.Context) ([]domain.OperationDuration, error) {
	response, err := c.writerClient.GetOperationDurations(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	durations := make([]domain.OperationDuration, len(response.OperationDurations))
	for i, gd := range response.OperationDurations {
		durations[i] = *c.unwrapOperationDuration(gd)
	}
	return durations, err
}

func (c *client) SkipAgentSubExpressions(ctx context.Context, agentName string) error {
	req := &writergrpc.SkipAgentSubExpressionsRequest{AgentName: agentName}
	_, err := c.writerClient.SkipAgentSubExpressions(ctx, req)
	return err
}
