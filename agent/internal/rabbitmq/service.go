package rabbitmq

import (
	"context"
	"github.com/sirupsen/logrus"
)

type Service interface {
	Run(ctx context.Context) error
	Shutdown() error
}

type service struct {
	conn Connection
}

func NewService(conn Connection) *service {
	return &service{conn: conn}
}

func (s *service) Run(ctx context.Context) error {
	if s.conn.Connection() != nil {
		return nil
	}
	logrus.Info("RabbitMQ service. Establishing connection")

	if err := s.conn.Connect(ctx); err != nil {
		return err
	}

	return nil
}

func (s *service) Shutdown() error {
	if s.conn == nil {
		return nil
	}
	logrus.Info("RabbitMQ service. closing connection")

	if err := s.conn.Close(); err != nil {
		return err
	}
	s.conn = nil
	return nil
}
