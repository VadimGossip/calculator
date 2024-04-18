package app

import (
	"github.com/VadimGossip/calculator/dbagent/internal/expression"
	"github.com/VadimGossip/calculator/dbagent/internal/token"
	"github.com/VadimGossip/calculator/dbagent/internal/user"
	"github.com/VadimGossip/calculator/dbagent/internal/writer"
)

type Factory struct {
	dbAdapter *DBAdapter

	exprService   expression.Service
	userService   user.Service
	tokenService  token.Service
	writerService writer.Service
}

var factory *Factory

func newFactory(dbAdapter *DBAdapter) *Factory {
	factory = &Factory{dbAdapter: dbAdapter}
	factory.exprService = expression.NewService(dbAdapter.exprRepo)
	factory.userService = user.NewService(dbAdapter.userRepo)
	factory.tokenService = token.NewService(dbAdapter.tokenRepo)
	factory.writerService = writer.NewService(factory.exprService, factory.userService, factory.tokenService)
	return factory
}
