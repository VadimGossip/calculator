package user

import (
	"context"
	"github.com/VadimGossip/calculator/dbagent/internal/domain"
	"github.com/VadimGossip/calculator/dbagent/internal/token"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type service struct {
	repo         Repository
	hasher       PasswordHasher
	tokenService token.Service
}

type Service interface {
	Register(ctx context.Context, user *domain.User) error
	Login(ctx context.Context, credentials domain.Credentials) (string, string, error)
}

var _ Service = (*service)(nil)

func NewService(repo Repository, hasher PasswordHasher, tokenService token.Service) *service {
	return &service{repo: repo,
		hasher:       hasher,
		tokenService: tokenService,
	}
}

func (s *service) Register(ctx context.Context, user *domain.User) error {
	password, err := s.hasher.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = password
	err = s.repo.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Login(ctx context.Context, credentials domain.Credentials) (string, string, error) {
	password, err := s.hasher.Hash(credentials.Password)
	if err != nil {
		return "", "", err
	}

	user, err := s.repo.GetByCredentials(ctx, credentials.Login, password)
	if err != nil {
		return "", "", err
	}

	return s.tokenService.GenerateTokens(ctx, user.Id)
}
