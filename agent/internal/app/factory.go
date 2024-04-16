package app

import (
	"github.com/VadimGossip/calculator/agent/internal/api/client/writer"
	"github.com/VadimGossip/calculator/agent/internal/domain"
	"github.com/VadimGossip/calculator/agent/internal/rabbitmq"
	"github.com/VadimGossip/calculator/agent/internal/worker"
)

type Factory struct {
	rabbitConn     rabbitmq.Connection
	rabbitConsumer rabbitmq.Consumer
	rabbitService  rabbitmq.Service

	writerClient  writer.Client
	workerService worker.Service
}

var factory *Factory

func newFactory(cfg *domain.Config) *Factory {
	factory = &Factory{}
	factory.rabbitConn = rabbitmq.NewConnection(cfg.AMPQServerConfig.Url)
	factory.rabbitService = rabbitmq.NewService(factory.rabbitConn)

	factory.writerClient = writer.NewClient("calculator-dbagent", 8085) //cfg

	factory.workerService = worker.NewService(cfg.Agent, factory.writerClient)
	factory.rabbitConsumer = rabbitmq.NewConsumer(cfg.AMPQStructCfg, factory.rabbitConn, factory.workerService)

	return factory
}
