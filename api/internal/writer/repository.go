package writer

import (
	"context"
	"database/sql"
	"github.com/VadimGossip/calculator/api/internal/domain"
	"time"
)

type Repository interface {
	CreateExpression(ctx context.Context, e *domain.Expression) error
	UpdateExpression(ctx context.Context, e domain.Expression) error
	GetExpressionSummaryBySeId(ctx context.Context, seId int64) (domain.Expression, error)
	GetExpressions(ctx context.Context) ([]domain.Expression, error)
	GetAgent(ctx context.Context, key string) (domain.Agent, error)
	CreateAgent(ctx context.Context, name string) error
	SetAgentHeartbeatAt(ctx context.Context, name string) (bool, error)
	GetAgents(ctx context.Context) ([]domain.Agent, error)
	GetOperationDuration(ctx context.Context, name string) (domain.OperationDuration, error)
	CreateOperationDuration(ctx context.Context, name string, duration uint16) error
	UpdateOperationDuration(ctx context.Context, name string, duration uint16) (bool, error)
	GetOperationDurations(ctx context.Context) ([]domain.OperationDuration, error)
	CreateSubExpression(ctx context.Context, s *domain.SubExpression) error
	StartSubExpressionEval(ctx context.Context, seId int64, agent string) (bool, error)
	StopSubExpressionEval(ctx context.Context, seId int64, result *float64) error
	GetSubExpressionIsLast(ctx context.Context, seId int64) (bool, error)
	DeleteSubExpressions(ctx context.Context, seId int64) error
	GetReadySubExpressions(ctx context.Context, expressionId *int64, skipTimeout time.Duration) ([]domain.SubExpression, error)
	SkipAgentSubExpressions(ctx context.Context, agent string) error
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

func (r *repository) UpdateExpression(ctx context.Context, e domain.Expression) error {
	updStmt := `UPDATE expressions 
                   SET result =$1,
                       state = $2,
                       error_msg =$3,
                       eval_started_at = $4,
                       eval_finished_at = $5
                 WHERE id = $6;`
	_, err := r.db.ExecContext(ctx, updStmt, valPointerToNullVal(e.Result), e.State, e.ErrorMsg, e.EvalStartedAt, e.EvalFinishedAt, e.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetExpressionSummaryBySeId(ctx context.Context, seId int64) (domain.Expression, error) {
	var e domain.Expression
	selectStmt := `SELECT e.id
	    				 ,e.value
					     ,max(case
							   when se.is_last then se.result
					      end) as result
					     ,e.state
					     ,e.created_at
					     ,min(se.eval_started_at) as eval_started_at
					     ,max(se.eval_finished_at) as eval_finished_at
			   	     FROM expressions e
				     JOIN sub_expressions se ON se.expression_id = e.id
				    WHERE exists(select 1
					  			   from sub_expressions se2
							      where se2.id = $1
								    and  se2.expression_id = se.expression_id)
				 group by e.id
					     ,e.value
					     ,e.result
					     ,e.state
					     ,e.error_msg
					     ,e.created_at
					     ,e.eval_started_at
					     ,e.eval_finished_at;`
	return e, r.db.QueryRowContext(ctx, selectStmt, seId).Scan(&e.Id, &e.Value, &e.Result, &e.State, &e.CreatedAt, &e.EvalStartedAt, &e.EvalFinishedAt)
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

func (r *repository) SetAgentHeartbeatAt(ctx context.Context, name string) (bool, error) {
	updStmt := `UPDATE agents 
                   SET last_heartbeat_at = $1
                 WHERE name = $2;`
	result, err := r.db.ExecContext(ctx, updStmt, time.Now(), name)
	if err != nil {
		return false, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows == 1, nil
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

func (r *repository) UpdateOperationDuration(ctx context.Context, name string, duration uint16) (bool, error) {
	updStmt := `UPDATE operation_durations 
                   SET duration = $1
                      ,updated_at = $2
                 WHERE operation_name = $3;`
	result, err := r.db.ExecContext(ctx, updStmt, duration, time.Now(), name)
	if err != nil {
		return false, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows == 1, nil
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
	createStmt := "INSERT INTO sub_expressions(expression_id, val1, val2, sub_expression_id1, sub_expression_id2, operation_name, is_last)" +
		"VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"

	return r.db.QueryRowContext(ctx, createStmt, s.ExpressionId, valPointerToNullVal(s.Val1), valPointerToNullVal(s.Val2), valPointerToNullVal(s.SubExpressionId1), valPointerToNullVal(s.SubExpressionId2), s.Operation, s.IsLast).Scan(&s.Id)
}

func (r *repository) StartSubExpressionEval(ctx context.Context, seId int64, agent string) (bool, error) {
	updStmt := `UPDATE sub_expressions 
                   SET agent_name = $1
                      ,eval_started_at = $2
                 WHERE id = $3
                   AND agent_name is null`
	result, err := r.db.ExecContext(ctx, updStmt, agent, time.Now(), seId)
	if err != nil {
		return false, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows != 1, nil
}

func (r *repository) StopSubExpressionEval(ctx context.Context, seId int64, result *float64) error {
	updStmt := `UPDATE sub_expressions 
                   SET result = $1
                      ,eval_finished_at = $2
                 WHERE id = $3
                   AND eval_finished_at is null;`
	_, err := r.db.ExecContext(ctx, updStmt, valPointerToNullVal(result), time.Now(), seId)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetSubExpressionIsLast(ctx context.Context, seId int64) (bool, error) {
	var isLast bool
	selectStmt := `SELECT is_last
                     FROM sub_expressions
                    WHERE id = $1;`
	return isLast, r.db.QueryRowContext(ctx, selectStmt, seId).Scan(&isLast)
}

func (r *repository) DeleteSubExpressions(ctx context.Context, seId int64) error {
	deleteStmt := `delete from sub_expressions se
      where exists (select 1
                      from sub_expressions p
                     where p.expression_id = se.expression_id
                       and p.id = $1);`
	_, err := r.db.ExecContext(ctx, deleteStmt, seId)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) SkipStartSubExpressionEval(ctx context.Context, seId int64) error {
	updStmt := `UPDATE sub_expressions 
                   SET agent_name = null
                      ,eval_started_at = null
                 WHERE id = $1
                   AND eval_finished_at is null;`
	_, err := r.db.ExecContext(ctx, updStmt, seId)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetReadySubExpressions(ctx context.Context, expressionId *int64, skipTimeout time.Duration) ([]domain.SubExpression, error) {
	selectStmt := `select se.id
			     		  ,coalesce(se.val1, se1.result) as val1
				    	  ,coalesce(se.val2, se2.result) as val2
						  ,se.operation_name
     					  ,coalesce(d.duration, 0) as operation_duration
						  ,se.eval_started_at						  
					 from sub_expressions se
				left join sub_expressions se1 on se.sub_expression_id1 = se1.id
				left join sub_expressions se2 on se.sub_expression_id2 = se2.id
				left join operation_durations d on d.operation_name = se.operation_name
					where se.expression_id = coalesce($1, se.expression_id)
					  and coalesce(se.val1, se1.result) is not null
					  and coalesce(se.val2, se2.result) is not null
					  and se.result is null
                      and se.eval_finished_at is null;`
	rows, err := r.db.QueryContext(ctx, selectStmt, valPointerToNullVal(expressionId))
	if err != nil {
		return nil, err
	}
	result := make([]domain.SubExpression, 0)
	for rows.Next() {
		var se domain.SubExpression
		var evalStartedAt sql.NullTime
		if err = rows.Scan(&se.Id, &se.Val1, &se.Val2, &se.Operation, &se.OperationDuration, &evalStartedAt); err != nil {
			return nil, err
		}
		if evalStartedAt.Valid {
			se.EvalStartedAt = evalStartedAt.Time
		}
		isHung := time.Since(se.EvalStartedAt) > skipTimeout
		if isHung {
			if err = r.SkipStartSubExpressionEval(ctx, se.Id); err != nil {
				return nil, err
			}
		}

		if se.EvalStartedAt.Equal(time.Time{}) || isHung {
			result = append(result, se)
		}
	}
	return result, nil
}

func (r *repository) SkipAgentSubExpressions(ctx context.Context, agent string) error {
	updStmt := `UPDATE sub_expressions 
                   SET agent_name = null
                      ,eval_started_at = null
                 WHERE agent_name = $1
                   and eval_started_at is not null
                   and eval_finished_at is null;`
	_, err := r.db.ExecContext(ctx, updStmt, agent)
	if err != nil {
		return err
	}
	return nil
}
