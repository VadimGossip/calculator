package writer

import (
	"context"
	"database/sql"
	"github.com/VadimGossip/calculator/api/internal/domain"
	"github.com/pkg/errors"
	"time"
)

type Repository interface {
	CreateExpression(ctx context.Context, e *domain.Expression) error
	StartExpressionEval(ctx context.Context, id int64) error
	SaveExpressionResult(ctx context.Context, id int64, result int) error
	GetExpressions(ctx context.Context) ([]domain.Expression, error)
	GetAgent(ctx context.Context, key string) (domain.Agent, error)
	CreateAgent(ctx context.Context, name string) error
	SetAgentHeartbeatAt(ctx context.Context, name string) error
	GetAgents(ctx context.Context) ([]domain.Agent, error)
	GetOperationDuration(ctx context.Context, name string) (domain.OperationDuration, error)
	CreateOperationDuration(ctx context.Context, name string, duration uint16) error
	UpdateOperationDuration(ctx context.Context, name string, duration uint16) error
	GetOperationDurations(ctx context.Context) ([]domain.OperationDuration, error)
	CreateSubExpression(ctx context.Context, s *domain.SubExpression) error
}

type repository struct {
	db *sql.DB
}

var _ Repository = (*repository)(nil)

func NewRepository(db *sql.DB) *repository {
	return &repository{db}
}

func valPointerToNullVal(val any) any {
	switch v := val.(type) {
	case *int:
		if val == (*int)(nil) {
			return sql.NullInt32{}
		}
		return sql.NullInt32{
			Int32: int32(*v),
			Valid: true,
		}
	case *int32:
		if val == (*int32)(nil) {
			return sql.NullInt32{}
		}
		return sql.NullInt32{
			Int32: *v,
			Valid: true,
		}
	case *int64:
		if val == (*int64)(nil) {
			return sql.NullInt64{}
		}
		return sql.NullInt64{
			Int64: *v,
			Valid: true,
		}
	case *float64:
		if val == (*float64)(nil) {
			return sql.NullFloat64{}
		}
		return sql.NullFloat64{
			Float64: *v,
			Valid:   true,
		}
	default:
		return nil
	}
}

func (r *repository) CreateExpression(ctx context.Context, e *domain.Expression) error {
	createStmt := "INSERT INTO expressions(value, state, created_at)" +
		"VALUES ($1, $2, $3) RETURNING id"

	return r.db.QueryRowContext(ctx, createStmt,
		e.Value, domain.ExpressionStateNew, time.Now()).
		Scan(&e.Id)
}

func (r *repository) StartExpressionEval(ctx context.Context, id int64) error {
	updStmt := `UPDATE expressions 
                   SET state = $1,
                       eval_started_at = $2
                 WHERE id =$3;`
	_, err := r.db.ExecContext(ctx, updStmt, domain.ExpressionStateInProgress, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) SaveExpressionResult(ctx context.Context, id int64, result int) error {
	updStmt := `UPDATE expressions 
                   SET result = $1, 
                       state = $2,
                       eval_finished_at = $3
                 WHERE id =$4;`
	_, err := r.db.ExecContext(ctx, updStmt, result, domain.ExpressionStateOK, time.Now(), id)
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

func (r *repository) GetAgent(ctx context.Context, key string) (domain.Agent, error) {
	var a domain.Agent
	selectStmt := `SELECT name
                         ,created_at
                         ,last_heartbeat_at
                     FROM agents
                    WHERE name = $1;`
	return a, r.db.QueryRowContext(ctx, selectStmt, key).Scan(&a.Name, &a.CreatedAt, &a.LastHeartbeatAt)
}

func (r *repository) CreateAgent(ctx context.Context, name string) error {
	createStmt := `INSERT 
                     INTO agents(name, created_at, last_heartbeat_at)
		           VALUES ($1, $2, $3);`

	_, err := r.db.ExecContext(ctx, createStmt, name, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) SetAgentHeartbeatAt(ctx context.Context, name string) error {
	updStmt := `UPDATE agents 
                   SET last_heartbeat_at = $1
                 WHERE name =$1;`
	_, err := r.db.ExecContext(ctx, updStmt, time.Now(), name)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetAgents(ctx context.Context) ([]domain.Agent, error) {
	selectStmt := `SELECT name
                         ,created_at
                         ,last_heartbeat_at
                     FROM agents;`
	rows, err := r.db.QueryContext(ctx, selectStmt)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Agent, 0)
	for rows.Next() {
		var a domain.Agent
		if err = rows.Scan(&a.Name, &a.CreatedAt, &a.LastHeartbeatAt); err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	return result, nil
}

func (r *repository) GetOperationDuration(ctx context.Context, name string) (domain.OperationDuration, error) {
	var d domain.OperationDuration
	selectStmt := `SELECT operation_name
                         ,duration
                         ,created_at
                         ,updated_at
                     FROM operation_durations
                    WHERE operation_name = $1;`
	return d, r.db.QueryRowContext(ctx, selectStmt, name).Scan(&d.Name, &d.Duration, &d.CreatedAt, &d.UpdatedAt)
}

func (r *repository) CreateOperationDuration(ctx context.Context, name string, duration uint16) error {
	createStmt := `INSERT 
                     INTO operation_durations(operation_name, duration, created_at, updated_at)
		           VALUES ($1, $2, $3, $4);`

	_, err := r.db.ExecContext(ctx, createStmt, name, duration, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateOperationDuration(ctx context.Context, name string, duration uint16) error {
	updStmt := `UPDATE operation_durations 
                   SET duration = $1
                      ,updated_at = $2
                 WHERE operation_name = $3;`
	_, err := r.db.ExecContext(ctx, updStmt, time.Now(), duration, time.Now(), name)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetOperationDurations(ctx context.Context) ([]domain.OperationDuration, error) {
	selectStmt := `SELECT operation_name
                         ,duration
                         ,created_at
                         ,updated_at
                     FROM operation_durations;`
	rows, err := r.db.QueryContext(ctx, selectStmt)
	if err != nil {
		return nil, err
	}
	result := make([]domain.OperationDuration, 0)
	for rows.Next() {
		var d domain.OperationDuration
		if err = rows.Scan(&d.Name, &d.Duration, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, nil
}

func (r *repository) CreateSubExpression(ctx context.Context, s *domain.SubExpression) error {
	createStmt := `INSERT 
                     INTO sub_expressions(expression_id, val1, val2, sub_expression_id1, sub_expression_id2, operation_name, is_last) 
                   VALUES ($1, $2, $3, $4, $5, $6, $7)
                RETURNING id into $1;`
	_, err := r.db.ExecContext(ctx, createStmt, s.ExpressionsId, valPointerToNullVal(s.Val1), valPointerToNullVal(s.Val2), valPointerToNullVal(s.SubExpressionId1), valPointerToNullVal(s.SubExpressionId2), s.OperationName, s.IsLast, s.Id)
	if err != nil {
		return errors.Wrap(err, "failed to create sub expression")
	}
	return nil
}
