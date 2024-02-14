package app

import (
	"github.com/VadimGossip/calculator/agent/internal/domain"
	"github.com/VadimGossip/calculator/agent/internal/rabbitmq"
)

type Factory struct {
	rabbitConn     rabbitmq.Connection
	rabbitConsumer rabbitmq.Consumer
	rabbitService  rabbitmq.Service
}

var factory *Factory

func newFactory(cfg *domain.Config) *Factory {
	factory = &Factory{}
	factory.rabbitConn = rabbitmq.NewConnection(cfg.AMPQServerConfig.Url)
	factory.rabbitService = rabbitmq.NewService(factory.rabbitConn)
	factory.rabbitConsumer = rabbitmq.NewConsumer(cfg.AMPQStructCfg, factory.rabbitConn)

	return factory
}
