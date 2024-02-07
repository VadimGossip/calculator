package app

import (
	"github.com/VadimGossip/calculator/api/internal/domain"
	"github.com/VadimGossip/calculator/api/internal/orchestrator"
	"github.com/VadimGossip/calculator/api/internal/rabbitmq"
	"github.com/VadimGossip/calculator/api/internal/writer"
)

type Factory struct {
	dbAdapter *DBAdapter

	writerService writer.Service

	rabbitConn     rabbitmq.Connection
	rabbitProducer rabbitmq.Producer
	rabbitConsumer rabbitmq.Consumer
	rabbitService  rabbitmq.Service

	orchestratorService orchestrator.Service
}

var factory *Factory

func newFactory(cfg *domain.Config, dbAdapter *DBAdapter) *Factory {
	factory = &Factory{dbAdapter: dbAdapter}
	factory.rabbitConn = rabbitmq.NewConnection(cfg.AMPQServerConfig.Url)
	factory.rabbitProducer = rabbitmq.NewProducer(cfg.AMPQStructCfg, factory.rabbitConn)
	factory.rabbitService = rabbitmq.NewService(factory.rabbitConn)
	factory.rabbitConsumer = rabbitmq.NewConsumer(cfg.AMPQStructCfg, factory.rabbitConn)

	factory.writerService = writer.NewService(dbAdapter.writerRepo)
	factory.orchestratorService = orchestrator.NewService(factory.writerService, factory.rabbitProducer)
	return factory
}
