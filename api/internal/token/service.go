package token

import (
	"context"
	"errors"
	"fmt"
	"github.com/VadimGossip/calculator/api/internal/api/client/writer"
	"github.com/VadimGossip/calculator/api/internal/domain"
	"github.com/VadimGossip/calculator/api/pkg/util"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

type service struct {
	writerClient writer.Client
	hmacSecret   []byte
	accessTTL    time.Duration
	refreshTTL   time.Duration
}

type Service interface {
	GenerateTokens(ctx context.Context, userId int64) (string, string, error)
	ParseToken(token string) (int64, error)
	GetRefreshTokenTTL() time.Duration
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
}

var _ Service = (*service)(nil)

func NewService(writerClient writer.Client, hmacSecret []byte, accessTTL, refreshTTL time.Duration) *service {
	return &service{writerClient: writerClient,
		hmacSecret: hmacSecret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *service) GenerateTokens(ctx context.Context, userId int64) (string, string, error) {
	now := time.Now()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userId,
		"exp":    now.Add(s.accessTTL).Unix(),
		"iat":    now.Unix(),
	})

	accessTokenStr, err := t.SignedString(s.hmacSecret)
	if err != nil {
		return "", "", err
	}

	refreshTokenStr, err := util.NewRandString(32)
	if err != nil {
		return "", "", err
	}

	if err = s.writerClient.CreateToken(ctx, &domain.Token{
		UserId:    userId,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(s.refreshTTL),
	}); err != nil {
		return "", "", err
	}

	return accessTokenStr, refreshTokenStr, nil
}

func (s *service) ParseToken(token string) (int64, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.hmacSecret, nil
	})
	if err != nil {
		return 0, err
	}

	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	userVal, ok := claims["id"].(string)
	if !ok {
		return 0, errors.New("invalid user_id")
	}
	userId, err := strconv.ParseInt(userVal, 10, 64)
	if !ok {
		return 0, errors.New("invalid user_id")
	}

	return userId, nil
}

func (s *service) GetRefreshTokenTTL() time.Duration {
	return s.refreshTTL
}

func (s *service) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	token, err := s.writerClient.GetToken(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	if token.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", domain.ErrRefreshTokenExpired
	}

	return s.GenerateTokens(ctx, token.UserId)
}
