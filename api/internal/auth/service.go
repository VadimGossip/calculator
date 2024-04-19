package auth

import (
	"context"
	"github.com/VadimGossip/calculator/api/internal/domain"
	"github.com/VadimGossip/calculator/api/internal/token"
	"github.com/VadimGossip/calculator/api/internal/user"
	"time"
)

type service struct {
	userService  user.Service
	tokenService token.Service
}

type Service interface {
	Register(ctx context.Context, user *domain.User) error
	Login(ctx context.Context, credentials domain.Credentials) (string, string, error)
	ParseToken(token string) (int64, error)
	GetRefreshTokenTTL() time.Duration
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
}

var _ Service = (*service)(nil)

func NewService(userService user.Service, tokenService token.Service) *service {
	return &service{userService: userService,
		tokenService: tokenService,
	}
}

func (s *service) Register(ctx context.Context, user *domain.User) error {
	return s.userService.Register(ctx, user)
}

func (s *service) Login(ctx context.Context, credentials domain.Credentials) (string, string, error) {
	userId, err := s.userService.Login(ctx, credentials)
	if err != nil {
		return "", "", err
	}

	return s.tokenService.GenerateTokens(ctx, userId)
}

func (s *service) ParseToken(token string) (int64, error) {
	return s.tokenService.ParseToken(token)
}

func (s *service) GetRefreshTokenTTL() time.Duration {
	return s.tokenService.GetRefreshTokenTTL()
}

func (s *service) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	return s.tokenService.RefreshTokens(ctx, refreshToken)
}
