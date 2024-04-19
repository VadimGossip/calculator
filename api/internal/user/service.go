package user

import (
	"context"
	"github.com/VadimGossip/calculator/api/internal/api/client/writer"
	"github.com/VadimGossip/calculator/api/internal/domain"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type service struct {
	writerClient writer.Client
	hasher       PasswordHasher
}

type Service interface {
	Register(ctx context.Context, user *domain.User) error
	Login(ctx context.Context, credentials domain.Credentials) (int64, error)
}

var _ Service = (*service)(nil)

func NewService(writerClient writer.Client, hasher PasswordHasher) *service {
	return &service{writerClient: writerClient,
		hasher: hasher,
	}
}

func (s *service) Register(ctx context.Context, user *domain.User) error {
	password, err := s.hasher.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = password
	err = s.writerClient.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Login(ctx context.Context, credentials domain.Credentials) (int64, error) {
	password, err := s.hasher.Hash(credentials.Password)
	if err != nil {
		return 0, err
	}

	user, err := s.writerClient.GetUserByCredentials(ctx, credentials.Login, password)
	if err != nil {
		return 0, err
	}

	return user.Id, nil
}
