// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.5
// source: dbagent/internal/api/grpcservice/writergrpc/writer.proto

package writergrpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	WriterService_Heartbeat_FullMethodName               = "/writer.WriterService/Heartbeat"
	WriterService_StartEval_FullMethodName               = "/writer.WriterService/StartEval"
	WriterService_StopEval_FullMethodName                = "/writer.WriterService/StopEval"
	WriterService_GetReadySubExpressions_FullMethodName  = "/writer.WriterService/GetReadySubExpressions"
	WriterService_GetExpressionByReqUid_FullMethodName   = "/writer.WriterService/GetExpressionByReqUid"
	WriterService_CreateExpression_FullMethodName        = "/writer.WriterService/CreateExpression"
	WriterService_CreateSubExpression_FullMethodName     = "/writer.WriterService/CreateSubExpression"
	WriterService_GetExpressions_FullMethodName          = "/writer.WriterService/GetExpressions"
	WriterService_GetAgents_FullMethodName               = "/writer.WriterService/GetAgents"
	WriterService_SaveOperationDuration_FullMethodName   = "/writer.WriterService/SaveOperationDuration"
	WriterService_GetOperationDurations_FullMethodName   = "/writer.WriterService/GetOperationDurations"
	WriterService_SkipAgentSubExpressions_FullMethodName = "/writer.WriterService/SkipAgentSubExpressions"
	WriterService_CreateUser_FullMethodName              = "/writer.WriterService/CreateUser"
	WriterService_GetUserByCred_FullMethodName           = "/writer.WriterService/GetUserByCred"
	WriterService_CreateToken_FullMethodName             = "/writer.WriterService/CreateToken"
	WriterService_GetToken_FullMethodName                = "/writer.WriterService/GetToken"
)

// WriterServiceClient is the client API for WriterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WriterServiceClient interface {
	Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	StartEval(ctx context.Context, in *StartEvalRequest, opts ...grpc.CallOption) (*StartEvalResponse, error)
	StopEval(ctx context.Context, in *StopEvalRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetReadySubExpressions(ctx context.Context, in *ReadySubExpressionsRequest, opts ...grpc.CallOption) (*ReadySubExpressionsResponse, error)
	GetExpressionByReqUid(ctx context.Context, in *ExpressionByReqUidRequest, opts ...grpc.CallOption) (*Expression, error)
	CreateExpression(ctx context.Context, in *CreateExpressionRequest, opts ...grpc.CallOption) (*CreateExpressionResponse, error)
	CreateSubExpression(ctx context.Context, in *CreateSubExpressionRequest, opts ...grpc.CallOption) (*CreateSubExpressionResponse, error)
	GetExpressions(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetExpressionsResponse, error)
	GetAgents(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetAgentsResponse, error)
	SaveOperationDuration(ctx context.Context, in *CreateOperDurRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetOperationDurations(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetOperDurResponse, error)
	SkipAgentSubExpressions(ctx context.Context, in *SkipAgentSubExpressionsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*User, error)
	GetUserByCred(ctx context.Context, in *GetUserByCredRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
	CreateToken(ctx context.Context, in *CreateTokenRequest, opts ...grpc.CallOption) (*CreateTokenResponse, error)
	GetToken(ctx context.Context, in *GetTokenRequest, opts ...grpc.CallOption) (*GetTokenResponse, error)
}

type writerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWriterServiceClient(cc grpc.ClientConnInterface) WriterServiceClient {
	return &writerServiceClient{cc}
}

func (c *writerServiceClient) Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, WriterService_Heartbeat_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writerServiceClient) StartEval(ctx context.Context, in *StartEvalRequest, opts ...grpc.CallOption) (*StartEvalResponse, error) {
	out := new(StartEvalResponse)
	err := c.cc.Invoke(ctx, WriterService_StartEval_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writerServiceClient) StopEval(ctx context.Context, in *StopEvalRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, WriterService_StopEval_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writerServiceClient) GetReadySubExpressions(ctx context.Context, in *ReadySubExpressionsRequest, opts ...grpc.CallOption) (*ReadySubExpressionsResponse, error) {
	out := new(ReadySubExpressionsResponse)
	err := c.cc.Invoke(ctx, WriterService_GetReadySubExpressions_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writerServiceClient) GetExpressionByReqUid(ctx context.Context, in *ExpressionByReqUidRequest, opts ...grpc.CallOption) (*Expression, error) {
	out := new(Expression)
	err := c.cc.Invoke(ctx, WriterService_GetExpressionByReqUid_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writerServiceClient) CreateExpression(ctx context.Context, in *CreateExpressionRequest, opts ...grpc.CallOption) (*CreateExpressionResponse, error) {
	out := new(CreateExpressionResponse)
	err := c.cc.Invoke(ctx, WriterService_CreateExpression_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writerServiceClient) CreateSubExpression(ctx context.Context, in *CreateSubExpressionRequest, opts ...grpc.CallOption) (*CreateSubExpressionResponse, error) {
	out := new(CreateSubExpressionResponse)
	err := c.cc.Invoke(ctx, WriterService_CreateSubExpression_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writerServiceClient) GetExpressions(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetExpressionsResponse, error) {
	out := new(GetExpressionsResponse)
	err := c.cc.Invoke(ctx, WriterService_GetExpressions_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writerServiceClient) GetAgents(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetAgentsResponse, error) {
	out := new(GetAgentsResponse)
	err := c.cc.Invoke(ctx, WriterService_GetAgents_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writerServiceClient) SaveOperationDuration(ctx context.Context, in *CreateOperDurRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, WriterService_SaveOperationDuration_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writerServiceClient) GetOperationDurations(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetOperDurResponse, error) {
	out := new(GetOperDurResponse)
	err := c.cc.Invoke(ctx, WriterService_GetOperationDurations_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writerServiceClient) SkipAgentSubExpressions(ctx context.Context, in *SkipAgentSubExpressionsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, WriterService_SkipAgentSubExpressions_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writerServiceClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, WriterService_CreateUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writerServiceClient) GetUserByCred(ctx context.Context, in *GetUserByCredRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, WriterService_GetUserByCred_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writerServiceClient) CreateToken(ctx context.Context, in *CreateTokenRequest, opts ...grpc.CallOption) (*CreateTokenResponse, error) {
	out := new(CreateTokenResponse)
	err := c.cc.Invoke(ctx, WriterService_CreateToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writerServiceClient) GetToken(ctx context.Context, in *GetTokenRequest, opts ...grpc.CallOption) (*GetTokenResponse, error) {
	out := new(GetTokenResponse)
	err := c.cc.Invoke(ctx, WriterService_GetToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WriterServiceServer is the server API for WriterService service.
// All implementations should embed UnimplementedWriterServiceServer
// for forward compatibility
type WriterServiceServer interface {
	Heartbeat(context.Context, *HeartbeatRequest) (*emptypb.Empty, error)
	StartEval(context.Context, *StartEvalRequest) (*StartEvalResponse, error)
	StopEval(context.Context, *StopEvalRequest) (*emptypb.Empty, error)
	GetReadySubExpressions(context.Context, *ReadySubExpressionsRequest) (*ReadySubExpressionsResponse, error)
	GetExpressionByReqUid(context.Context, *ExpressionByReqUidRequest) (*Expression, error)
	CreateExpression(context.Context, *CreateExpressionRequest) (*CreateExpressionResponse, error)
	CreateSubExpression(context.Context, *CreateSubExpressionRequest) (*CreateSubExpressionResponse, error)
	GetExpressions(context.Context, *emptypb.Empty) (*GetExpressionsResponse, error)
	GetAgents(context.Context, *emptypb.Empty) (*GetAgentsResponse, error)
	SaveOperationDuration(context.Context, *CreateOperDurRequest) (*emptypb.Empty, error)
	GetOperationDurations(context.Context, *emptypb.Empty) (*GetOperDurResponse, error)
	SkipAgentSubExpressions(context.Context, *SkipAgentSubExpressionsRequest) (*emptypb.Empty, error)
	CreateUser(context.Context, *CreateUserRequest) (*User, error)
	GetUserByCred(context.Context, *GetUserByCredRequest) (*CreateUserResponse, error)
	CreateToken(context.Context, *CreateTokenRequest) (*CreateTokenResponse, error)
	GetToken(context.Context, *GetTokenRequest) (*GetTokenResponse, error)
}

// UnimplementedWriterServiceServer should be embedded to have forward compatible implementations.
type UnimplementedWriterServiceServer struct {
}

func (UnimplementedWriterServiceServer) Heartbeat(context.Context, *HeartbeatRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Heartbeat not implemented")
}
func (UnimplementedWriterServiceServer) StartEval(context.Context, *StartEvalRequest) (*StartEvalResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartEval not implemented")
}
func (UnimplementedWriterServiceServer) StopEval(context.Context, *StopEvalRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopEval not implemented")
}
func (UnimplementedWriterServiceServer) GetReadySubExpressions(context.Context, *ReadySubExpressionsRequest) (*ReadySubExpressionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReadySubExpressions not implemented")
}
func (UnimplementedWriterServiceServer) GetExpressionByReqUid(context.Context, *ExpressionByReqUidRequest) (*Expression, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetExpressionByReqUid not implemented")
}
func (UnimplementedWriterServiceServer) CreateExpression(context.Context, *CreateExpressionRequest) (*CreateExpressionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateExpression not implemented")
}
func (UnimplementedWriterServiceServer) CreateSubExpression(context.Context, *CreateSubExpressionRequest) (*CreateSubExpressionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSubExpression not implemented")
}
func (UnimplementedWriterServiceServer) GetExpressions(context.Context, *emptypb.Empty) (*GetExpressionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetExpressions not implemented")
}
func (UnimplementedWriterServiceServer) GetAgents(context.Context, *emptypb.Empty) (*GetAgentsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAgents not implemented")
}
func (UnimplementedWriterServiceServer) SaveOperationDuration(context.Context, *CreateOperDurRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveOperationDuration not implemented")
}
func (UnimplementedWriterServiceServer) GetOperationDurations(context.Context, *emptypb.Empty) (*GetOperDurResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOperationDurations not implemented")
}
func (UnimplementedWriterServiceServer) SkipAgentSubExpressions(context.Context, *SkipAgentSubExpressionsRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SkipAgentSubExpressions not implemented")
}
func (UnimplementedWriterServiceServer) CreateUser(context.Context, *CreateUserRequest) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedWriterServiceServer) GetUserByCred(context.Context, *GetUserByCredRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserByCred not implemented")
}
func (UnimplementedWriterServiceServer) CreateToken(context.Context, *CreateTokenRequest) (*CreateTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateToken not implemented")
}
func (UnimplementedWriterServiceServer) GetToken(context.Context, *GetTokenRequest) (*GetTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetToken not implemented")
}

// UnsafeWriterServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WriterServiceServer will
// result in compilation errors.
type UnsafeWriterServiceServer interface {
	mustEmbedUnimplementedWriterServiceServer()
}

func RegisterWriterServiceServer(s grpc.ServiceRegistrar, srv WriterServiceServer) {
	s.RegisterService(&WriterService_ServiceDesc, srv)
}

func _WriterService_Heartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartbeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriterServiceServer).Heartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WriterService_Heartbeat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriterServiceServer).Heartbeat(ctx, req.(*HeartbeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WriterService_StartEval_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartEvalRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriterServiceServer).StartEval(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WriterService_StartEval_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriterServiceServer).StartEval(ctx, req.(*StartEvalRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WriterService_StopEval_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopEvalRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriterServiceServer).StopEval(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WriterService_StopEval_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriterServiceServer).StopEval(ctx, req.(*StopEvalRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WriterService_GetReadySubExpressions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadySubExpressionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriterServiceServer).GetReadySubExpressions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WriterService_GetReadySubExpressions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriterServiceServer).GetReadySubExpressions(ctx, req.(*ReadySubExpressionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WriterService_GetExpressionByReqUid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExpressionByReqUidRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriterServiceServer).GetExpressionByReqUid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WriterService_GetExpressionByReqUid_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriterServiceServer).GetExpressionByReqUid(ctx, req.(*ExpressionByReqUidRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WriterService_CreateExpression_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateExpressionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriterServiceServer).CreateExpression(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WriterService_CreateExpression_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriterServiceServer).CreateExpression(ctx, req.(*CreateExpressionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WriterService_CreateSubExpression_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSubExpressionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriterServiceServer).CreateSubExpression(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WriterService_CreateSubExpression_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriterServiceServer).CreateSubExpression(ctx, req.(*CreateSubExpressionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WriterService_GetExpressions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriterServiceServer).GetExpressions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WriterService_GetExpressions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriterServiceServer).GetExpressions(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _WriterService_GetAgents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriterServiceServer).GetAgents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WriterService_GetAgents_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriterServiceServer).GetAgents(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _WriterService_SaveOperationDuration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOperDurRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriterServiceServer).SaveOperationDuration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WriterService_SaveOperationDuration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriterServiceServer).SaveOperationDuration(ctx, req.(*CreateOperDurRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WriterService_GetOperationDurations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriterServiceServer).GetOperationDurations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WriterService_GetOperationDurations_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriterServiceServer).GetOperationDurations(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _WriterService_SkipAgentSubExpressions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SkipAgentSubExpressionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriterServiceServer).SkipAgentSubExpressions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WriterService_SkipAgentSubExpressions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriterServiceServer).SkipAgentSubExpressions(ctx, req.(*SkipAgentSubExpressionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WriterService_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriterServiceServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WriterService_CreateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriterServiceServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WriterService_GetUserByCred_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserByCredRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriterServiceServer).GetUserByCred(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WriterService_GetUserByCred_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriterServiceServer).GetUserByCred(ctx, req.(*GetUserByCredRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WriterService_CreateToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriterServiceServer).CreateToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WriterService_CreateToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriterServiceServer).CreateToken(ctx, req.(*CreateTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WriterService_GetToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriterServiceServer).GetToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WriterService_GetToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriterServiceServer).GetToken(ctx, req.(*GetTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// WriterService_ServiceDesc is the grpc.ServiceDesc for WriterService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WriterService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "writer.WriterService",
	HandlerType: (*WriterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Heartbeat",
			Handler:    _WriterService_Heartbeat_Handler,
		},
		{
			MethodName: "StartEval",
			Handler:    _WriterService_StartEval_Handler,
		},
		{
			MethodName: "StopEval",
			Handler:    _WriterService_StopEval_Handler,
		},
		{
			MethodName: "GetReadySubExpressions",
			Handler:    _WriterService_GetReadySubExpressions_Handler,
		},
		{
			MethodName: "GetExpressionByReqUid",
			Handler:    _WriterService_GetExpressionByReqUid_Handler,
		},
		{
			MethodName: "CreateExpression",
			Handler:    _WriterService_CreateExpression_Handler,
		},
		{
			MethodName: "CreateSubExpression",
			Handler:    _WriterService_CreateSubExpression_Handler,
		},
		{
			MethodName: "GetExpressions",
			Handler:    _WriterService_GetExpressions_Handler,
		},
		{
			MethodName: "GetAgents",
			Handler:    _WriterService_GetAgents_Handler,
		},
		{
			MethodName: "SaveOperationDuration",
			Handler:    _WriterService_SaveOperationDuration_Handler,
		},
		{
			MethodName: "GetOperationDurations",
			Handler:    _WriterService_GetOperationDurations_Handler,
		},
		{
			MethodName: "SkipAgentSubExpressions",
			Handler:    _WriterService_SkipAgentSubExpressions_Handler,
		},
		{
			MethodName: "CreateUser",
			Handler:    _WriterService_CreateUser_Handler,
		},
		{
			MethodName: "GetUserByCred",
			Handler:    _WriterService_GetUserByCred_Handler,
		},
		{
			MethodName: "CreateToken",
			Handler:    _WriterService_CreateToken_Handler,
		},
		{
			MethodName: "GetToken",
			Handler:    _WriterService_GetToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dbagent/internal/api/grpcservice/writergrpc/writer.proto",
}
