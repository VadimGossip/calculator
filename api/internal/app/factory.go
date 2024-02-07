package app

import (
	"github.com/VadimGossip/calculator/api/internal/domain"
	"github.com/VadimGossip/calculator/api/internal/manager"
	"github.com/VadimGossip/calculator/api/internal/rabbitmq"
	"github.com/VadimGossip/calculator/api/internal/writer"
)

type Factory struct {
	dbAdapter *DBAdapter

	writerService writer.Service

	rabbitConn     rabbitmq.Connection
	rabbitProducer rabbitmq.Producer
	rabbitService  rabbitmq.Service

	managerService manager.Service
}

var factory *Factory

func newFactory(cfg *domain.Config, dbAdapter *DBAdapter) *Factory {
	factory = &Factory{dbAdapter: dbAdapter}
	factory.rabbitConn = rabbitmq.NewConnection(cfg.AMPQServerConfig.Url)
	factory.rabbitProducer = rabbitmq.NewProducer(cfg.AMPQStructCfg, factory.rabbitConn)
	factory.rabbitService = rabbitmq.NewService(factory.rabbitConn)

	factory.writerService = writer.NewService(dbAdapter.writerRepo)
	factory.managerService = manager.NewService(factory.writerService, factory.rabbitProducer)
	return factory
}
