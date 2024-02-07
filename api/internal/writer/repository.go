package writer

import (
	"context"
	"database/sql"
	"time"
)

type Repository interface {
	CreateExpression(ctx context.Context, e *Expression) error
	SaveExpressionResult(ctx context.Context, id int64, result int) error
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
		e.Value, NewState, time.Now()).
		Scan(&e.Id)
}

func (r *repository) SaveExpressionResult(ctx context.Context, id int64, result int) error {
	updStmt := `UPDATE expressions 
                   SET result = $1, 
                       state = $2
                 WHERE id =$3;`
	_, err := r.db.ExecContext(ctx, updStmt, result, OkState, id)
	if err != nil {
		return err
	}
	return nil
}
