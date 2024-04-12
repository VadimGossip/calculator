package app

import (
	"github.com/VadimGossip/calculator/dbagent/internal/writer"
)

type Factory struct {
	dbAdapter *DBAdapter

	writerService writer.Service
}

var factory *Factory

func newFactory(dbAdapter *DBAdapter) *Factory {
	factory = &Factory{dbAdapter: dbAdapter}
	factory.writerService = writer.NewService(dbAdapter.writerRepo)
	return factory
}
