package app

import (
	wc "github.com/VadimGossip/calculator/api/internal/api/client/writer"
	"github.com/VadimGossip/calculator/api/internal/domain"
	"github.com/VadimGossip/calculator/api/internal/expression"
	"github.com/VadimGossip/calculator/api/internal/parser"
	"github.com/VadimGossip/calculator/api/internal/rabbitmq"
	"github.com/VadimGossip/calculator/api/internal/validation"
)

type Factory struct {
	writerClient wc.Client

	rabbitConn     rabbitmq.Connection
	rabbitProducer rabbitmq.Producer
	rabbitService  rabbitmq.Service

	parseService      parser.Service
	validationService validation.Service
	expressionService expression.Service
}

var factory *Factory

func newFactory(cfg *domain.Config) *Factory {
	factory = &Factory{}
	factory.rabbitConn = rabbitmq.NewConnection(cfg.AMPQServerConfig.Url)
	factory.rabbitProducer = rabbitmq.NewProducer(cfg.AMPQStructCfg, factory.rabbitConn)
	factory.rabbitService = rabbitmq.NewService(factory.rabbitConn)

	factory.writerClient = wc.NewClient("calculator-dbagent", 8085) //cfg
	factory.parseService = parser.NewService()
	factory.validationService = validation.NewService(cfg.Expression.MaxLength)
	factory.expressionService = expression.NewService(cfg.Expression, factory.parseService, factory.validationService, factory.writerClient, factory.rabbitProducer)
	return factory
}
