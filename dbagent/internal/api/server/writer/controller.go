package writer

import (
	"context"
	"github.com/VadimGossip/calculator/dbagent/internal/api/grpcservice/writergrpc"
	"github.com/VadimGossip/calculator/dbagent/internal/domain"
	"github.com/VadimGossip/calculator/dbagent/internal/writer"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type Controller interface {
	Heartbeat(ctx context.Context, req *writergrpc.HeartbeatRequest) (*emptypb.Empty, error)
	StartEval(ctx context.Context, req *writergrpc.StartEvalRequest) (*writergrpc.StartEvalResponse, error)
	StopEval(ctx context.Context, req *writergrpc.StopEvalRequest) (*emptypb.Empty, error)
	GetReadySubExpressions(ctx context.Context, req *writergrpc.ReadySubExpressionsRequest) (*writergrpc.ReadySubExpressionsResponse, error)
	GetExpressionByReqUid(ctx context.Context, req *writergrpc.ExpressionByReqUidRequest) (*writergrpc.Expression, error)
	CreateExpression(ctx context.Context, req *writergrpc.CreateExpressionRequest) (*writergrpc.CreateExpressionResponse, error)
	CreateSubExpression(ctx context.Context, req *writergrpc.CreateSubExpressionRequest) (*writergrpc.CreateSubExpressionResponse, error)
	GetExpressions(ctx context.Context, _ *emptypb.Empty) (*writergrpc.GetExpressionsResponse, error)
	GetAgents(ctx context.Context, _ *emptypb.Empty) (*writergrpc.GetAgentsResponse, error)
	SaveOperationDuration(ctx context.Context, req *writergrpc.CreateOperDurRequest) (*emptypb.Empty, error)
	GetOperationDurations(ctx context.Context, _ *emptypb.Empty) (*writergrpc.GetOperDurResponse, error)
	SkipAgentSubExpressions(ctx context.Context, req *writergrpc.SkipAgentSubExpressionsRequest) (*emptypb.Empty, error)
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

func (c *controller) wrapSubExpression(se *domain.SubExpression) *writergrpc.SubExpression {
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
		Id:                se.Id,
		Val1:              val1,
		Val2:              val2,
		SubExpressionId1:  subExpression1,
		SubExpressionId2:  subExpression2,
		Operation:         se.Operation,
		OperationDuration: se.OperationDuration,
		IsLast:            se.IsLast,
	}
}

func (c *controller) unwrapSubExpression(gse *writergrpc.SubExpression) *domain.SubExpression {
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

func (c *controller) wrapExpression(e *domain.Expression) *writergrpc.Expression {
	if e == nil {
		return &writergrpc.Expression{}
	}

	var evalStartedAt, evalFinishedAt int64
	if e.EvalStartedAt != nil {
		evalStartedAt = (*e.EvalStartedAt).Unix()
	}

	if e.EvalFinishedAt != nil {
		evalFinishedAt = (*e.EvalFinishedAt).Unix()
	}

	var result float64
	if e.Result != nil {
		result = *e.Result
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

func (c *controller) unwrapExpression(ge *writergrpc.Expression) *domain.Expression {
	return &domain.Expression{
		Id:     ge.Id,
		ReqUid: ge.ReqUid,
		Value:  ge.Value,
	}
}

func (c *controller) wrapAgent(a *domain.Agent) *writergrpc.Agent {
	return &writergrpc.Agent{
		Name:      a.Name,
		CreatedAt: a.CreatedAt.Unix(),
		LastHbAt:  a.LastHeartbeatAt.Unix(),
	}
}

func (c *controller) wrapOperationDuration(d *domain.OperationDuration) *writergrpc.OperationDuration {
	return &writergrpc.OperationDuration{
		Name:      d.Name,
		Duration:  d.Duration,
		CreatedAt: d.CreatedAt.Unix(),
		UpdatedAt: d.UpdatedAt.Unix(),
	}
}

func (c *controller) unwrapOperationDuration(gd *writergrpc.OperationDuration) *domain.OperationDuration {
	return &domain.OperationDuration{
		Name:      gd.Name,
		Duration:  gd.Duration,
		CreatedAt: time.Unix(0, gd.CreatedAt),
		UpdatedAt: time.Unix(0, gd.UpdatedAt),
	}
}

func (c *controller) GetReadySubExpressions(ctx context.Context, req *writergrpc.ReadySubExpressionsRequest) (*writergrpc.ReadySubExpressionsResponse, error) {
	var eId *int64
	if req.EId != -1 {
		eId = &req.EId
	}

	subExpressions, err := c.writer.GetReadySubExpressions(ctx, eId, req.SkipTimeoutSec)
	if err != nil {
		return &writergrpc.ReadySubExpressionsResponse{}, err
	}

	response := &writergrpc.ReadySubExpressionsResponse{
		SubExpressions: make([]*writergrpc.SubExpression, len(subExpressions)),
	}
	for i, se := range subExpressions {
		response.SubExpressions[i] = c.wrapSubExpression(&se)
	}
	return response, nil
}

func (c *controller) GetExpressionByReqUid(ctx context.Context, req *writergrpc.ExpressionByReqUidRequest) (*writergrpc.Expression, error) {
	e, err := c.writer.GetExpressionByReqUid(ctx, req.ReqUid)
	return c.wrapExpression(e), err
}

func (c *controller) CreateExpression(ctx context.Context, req *writergrpc.CreateExpressionRequest) (*writergrpc.CreateExpressionResponse, error) {
	e := c.unwrapExpression(req.E)
	if err := c.writer.CreateExpression(ctx, e); err != nil {
		return &writergrpc.CreateExpressionResponse{}, err
	}
	return &writergrpc.CreateExpressionResponse{Id: e.Id}, nil
}

func (c *controller) CreateSubExpression(cxt context.Context, req *writergrpc.CreateSubExpressionRequest) (*writergrpc.CreateSubExpressionResponse, error) {
	se := c.unwrapSubExpression(req.Se)
	if err := c.writer.CreateSubExpression(cxt, se); err != nil {
		return &writergrpc.CreateSubExpressionResponse{}, err
	}

	return &writergrpc.CreateSubExpressionResponse{Id: se.Id}, nil
}

func (c *controller) GetExpressions(ctx context.Context, _ *emptypb.Empty) (*writergrpc.GetExpressionsResponse, error) {
	expressions, err := c.writer.GetExpressions(ctx)
	if err != nil {
		return &writergrpc.GetExpressionsResponse{}, err
	}
	response := &writergrpc.GetExpressionsResponse{
		Expressions: make([]*writergrpc.Expression, len(expressions)),
	}
	for i, e := range expressions {
		response.Expressions[i] = c.wrapExpression(&e)
	}
	return response, nil
}

func (c *controller) GetAgents(ctx context.Context, _ *emptypb.Empty) (*writergrpc.GetAgentsResponse, error) {
	agents, err := c.writer.GetAgents(ctx)
	if err != nil {
		return &writergrpc.GetAgentsResponse{}, err
	}
	response := &writergrpc.GetAgentsResponse{
		Agents: make([]*writergrpc.Agent, len(agents)),
	}
	for i, a := range agents {
		response.Agents[i] = c.wrapAgent(&a)
	}
	return response, nil
}

func (c *controller) SaveOperationDuration(ctx context.Context, req *writergrpc.CreateOperDurRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, c.writer.SaveOperationDuration(ctx, req.Name, req.Duration)
}

func (c *controller) GetOperationDurations(ctx context.Context, _ *emptypb.Empty) (*writergrpc.GetOperDurResponse, error) {
	durations, err := c.writer.GetOperationDurations(ctx)
	if err != nil {
		return &writergrpc.GetOperDurResponse{}, err
	}
	response := &writergrpc.GetOperDurResponse{
		OperationDurations: make([]*writergrpc.OperationDuration, len(durations)),
	}
	for i, d := range durations {
		response.OperationDurations[i] = c.wrapOperationDuration(&d)
	}
	return response, nil
}

func (c *controller) SkipAgentSubExpressions(ctx context.Context, req *writergrpc.SkipAgentSubExpressionsRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, c.writer.SkipAgentSubExpressions(ctx, req.AgentName)
}
