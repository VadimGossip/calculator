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
	WriterService_Heartbeat_FullMethodName = "/writer.WriterService/Heartbeat"
)

// WriterServiceClient is the client API for WriterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WriterServiceClient interface {
	Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
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

// WriterServiceServer is the server API for WriterService service.
// All implementations should embed UnimplementedWriterServiceServer
// for forward compatibility
type WriterServiceServer interface {
	Heartbeat(context.Context, *HeartbeatRequest) (*emptypb.Empty, error)
}

// UnimplementedWriterServiceServer should be embedded to have forward compatible implementations.
type UnimplementedWriterServiceServer struct {
}

func (UnimplementedWriterServiceServer) Heartbeat(context.Context, *HeartbeatRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Heartbeat not implemented")
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
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dbagent/internal/api/grpcservice/writergrpc/writer.proto",
}
