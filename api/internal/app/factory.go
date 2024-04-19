package app

import (
	"github.com/VadimGossip/calculator/api/internal/api/client/writer"
	"github.com/VadimGossip/calculator/api/internal/auth"
	"github.com/VadimGossip/calculator/api/internal/domain"
	"github.com/VadimGossip/calculator/api/internal/expression"
	"github.com/VadimGossip/calculator/api/internal/parser"
	"github.com/VadimGossip/calculator/api/internal/rabbitmq"
	"github.com/VadimGossip/calculator/api/internal/token"
	"github.com/VadimGossip/calculator/api/internal/user"
	"github.com/VadimGossip/calculator/api/internal/validation"
	"github.com/VadimGossip/calculator/api/pkg/hash"
)

type Factory struct {
	writerClient writer.Client

	rabbitConn     rabbitmq.Connection
	rabbitProducer rabbitmq.Producer
	rabbitService  rabbitmq.Service

	parseService      parser.Service
	validationService validation.Service
	tokenService      token.Service
	userService       user.Service
	authService       auth.Service
	expressionService expression.Service
}

var factory *Factory

func newFactory(cfg *domain.Config) *Factory {
	factory = &Factory{}
	factory.rabbitConn = rabbitmq.NewConnection(cfg.AMPQServerConfig.Url)
	factory.rabbitProducer = rabbitmq.NewProducer(cfg.AMPQStructCfg, factory.rabbitConn)
	factory.rabbitService = rabbitmq.NewService(factory.rabbitConn)

	factory.writerClient = writer.NewClient(cfg.DbAgentGrpc.Host, cfg.DbAgentGrpc.Port)
	factory.parseService = parser.NewService()
	factory.validationService = validation.NewService(cfg.Expression.MaxLength)
	factory.tokenService = token.NewService(factory.writerClient, []byte(cfg.Auth.Secret), cfg.Auth.AccessTokenTTL, cfg.Auth.RefreshTokenTTL)
	factory.userService = user.NewService(factory.writerClient, hash.NewSHA1Hasher(cfg.Auth.Salt))
	factory.authService = auth.NewService(factory.userService, factory.tokenService)
	factory.expressionService = expression.NewService(cfg.Expression, factory.parseService, factory.validationService, factory.writerClient, factory.rabbitProducer)
	return factory
}
