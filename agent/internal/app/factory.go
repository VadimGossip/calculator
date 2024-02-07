package app

import (
	"github.com/VadimGossip/calculator/agent/internal/domain"
	"github.com/VadimGossip/calculator/agent/internal/rabbitmq"
	"github.com/VadimGossip/calculator/agent/internal/writer"
)

type Factory struct {
	dbAdapter *DBAdapter

	writerService writer.Service

	rabbitConn     rabbitmq.Connection
	rabbitConsumer rabbitmq.Consumer
	rabbitService  rabbitmq.Service
}

var factory *Factory

func newFactory(cfg *domain.Config, dbAdapter *DBAdapter) *Factory {
	factory = &Factory{dbAdapter: dbAdapter}
	factory.rabbitConn = rabbitmq.NewConnection(cfg.AMPQServerConfig.Url)
	factory.rabbitService = rabbitmq.NewService(factory.rabbitConn)
	factory.rabbitConsumer = rabbitmq.NewConsumer(cfg.AMPQStructCfg, factory.rabbitConn)

	factory.writerService = writer.NewService(dbAdapter.writerRepo)
	return factory
}
