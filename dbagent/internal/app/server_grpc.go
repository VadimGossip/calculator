package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"
)

type GrpcServer struct {
	port     int
	server   *grpc.Server
	listener net.Listener
}

func NewGrpcServer(port int) *GrpcServer {
	return &GrpcServer{
		port: port,
	}
}

func (s *GrpcServer) Listen(grpcRouter GrpcRouter) error {
	kaep := keepalive.EnforcementPolicy{
		MinTime:             10 * time.Second,
		PermitWithoutStream: true,
	}
	kasp := keepalive.ServerParameters{
		MaxConnectionIdle: 60 * time.Minute,
		Time:              30 * time.Second,
		Timeout:           3 * time.Second,
	}

	serverOptions := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kasp),
	}
	s.server = grpc.NewServer(serverOptions...)
	grpcRouter(s.server)

	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	logrus.Infof("[grpc/server] Starting on port: %d", s.port)

	return s.server.Serve(s.listener)
}
