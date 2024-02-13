package app

import (
	"github.com/VadimGossip/calculator/api/internal/domain"
	"github.com/VadimGossip/calculator/api/internal/manager"
	"github.com/VadimGossip/calculator/api/internal/parser"
	"github.com/VadimGossip/calculator/api/internal/rabbitmq"
	"github.com/VadimGossip/calculator/api/internal/writer"
)

type Factory struct {
	dbAdapter *DBAdapter

	writerService writer.Service

	rabbitConn     rabbitmq.Connection
	rabbitProducer rabbitmq.Producer
	rabbitService  rabbitmq.Service

	parseService   parser.Service
	managerService manager.Service
}

var factory *Factory

func newFactory(cfg *domain.Config, dbAdapter *DBAdapter) *Factory {
	factory = &Factory{dbAdapter: dbAdapter}
	factory.rabbitConn = rabbitmq.NewConnection(cfg.AMPQServerConfig.Url)
	factory.rabbitProducer = rabbitmq.NewProducer(cfg.AMPQStructCfg, factory.rabbitConn)
	factory.rabbitService = rabbitmq.NewService(factory.rabbitConn)

	factory.writerService = writer.NewService(dbAdapter.writerRepo)
	factory.parseService = parser.NewService()
	factory.managerService = manager.NewService(factory.parseService, factory.writerService, factory.rabbitProducer)
	return factory
}
