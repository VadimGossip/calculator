package writer

import (
	"context"
	"database/sql"
	"time"
)

type Repository interface {
	CreateExpression(ctx context.Context, e *Expression) error
	//UpdateExpressionState(id int64, state string) error
	//GetExpressions() ([]Expression, error)
}

type repository struct {
	db *sql.DB
}

var _ Repository = (*repository)(nil)

func NewRepository(db *sql.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateExpression(ctx context.Context, e *Expression) error {
	createStmt := "INSERT INTO expressions(value, state, created_at)" +
		"VALUES ($1, $2, $3) RETURNING id"

	return r.db.QueryRowContext(ctx, createStmt,
		e.Value, e.State, time.Now()).
		Scan(&e.Id)
}
