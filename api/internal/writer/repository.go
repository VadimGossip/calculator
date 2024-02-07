package writer

import (
	"context"
	"database/sql"
	"github.com/VadimGossip/calculator/api/internal/domain"
	"time"
)

type Repository interface {
	CreateExpression(ctx context.Context, e *domain.Expression) error
	SaveExpressionResult(ctx context.Context, id int64, result int) error
	GetExpressions(ctx context.Context) ([]domain.Expression, error)
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

func (r *repository) CreateExpression(ctx context.Context, e *domain.Expression) error {
	createStmt := "INSERT INTO expressions(value, state, created_at)" +
		"VALUES ($1, $2, $3) RETURNING id"

	return r.db.QueryRowContext(ctx, createStmt,
		e.Value, domain.ExpressionStateNew, time.Now()).
		Scan(&e.Id)
}

func (r *repository) SaveExpressionResult(ctx context.Context, id int64, result int) error {
	updStmt := `UPDATE expressions 
                   SET result = $1, 
                       state = $2
                 WHERE id =$3;`
	_, err := r.db.ExecContext(ctx, updStmt, result, domain.ExpressionStateOK, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetExpressions(ctx context.Context) ([]domain.Expression, error) {
	selectStmt := `SELECT id 
                         ,value
                         ,result
                         ,state
                         ,created_at
                         ,eval_started_at
                         ,eval_finished_at
                    FROM expressions;`
	rows, err := r.db.QueryContext(ctx, selectStmt)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Expression, 0)
	for rows.Next() {
		var e domain.Expression
		if err = rows.Scan(&e.Id, &e.Value, &e.Result, &e.State, &e.CreatedAt, &e.EvalStartedAt, &e.EvalFinishedAt); err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}
