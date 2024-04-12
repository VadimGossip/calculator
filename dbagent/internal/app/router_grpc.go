package app

import (
	"github.com/VadimGossip/calculator/dbagent/internal/api/grpcservice/writergrpc"
	"github.com/VadimGossip/calculator/dbagent/internal/api/server/writer"
	"google.golang.org/grpc"
)

type GrpcRouter func(*grpc.Server)

func initGrpcRouter(app *App) GrpcRouter {
	return func(s *grpc.Server) {
		c := writer.NewController(app.writerService)
		writergrpc.RegisterWriterServiceServer(s, c)
	}
}
