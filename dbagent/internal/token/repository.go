package token

import (
	"context"
	"database/sql"
	"github.com/VadimGossip/calculator/dbagent/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, token *domain.Token) error
	Get(ctx context.Context, tokenValue string) (*domain.Token, error)
}

type repository struct {
	db *sql.DB
}

var _ Repository = (*repository)(nil)

func NewRepository(db *sql.DB) *repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, token *domain.Token) error {
	createStmt := "insert into refresh_tokens(user_id, token, expires_at)" +
		"values ($1, $2, $3) returning id"

	return r.db.QueryRowContext(ctx, createStmt, token.UserId, token.Token, token.ExpiresAt).
		Scan(&token.Id)
}

func (r *repository) Get(ctx context.Context, tokenValue string) (*domain.Token, error) {
	var t domain.Token
	selectStmt := `select t.id
                         ,t.user_id
                         ,t.token
                         ,t.expires_at
			   	     from refresh_tokens t
                    where t.token = $1`

	if err := r.db.QueryRowContext(ctx, selectStmt, tokenValue).Scan(&t.Id, &t.UserId, &t.Token, &t.ExpiresAt); err != nil {
		return nil, err
	}
	return &t, nil
}
