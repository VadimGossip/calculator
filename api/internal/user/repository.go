package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/VadimGossip/calculator/api/internal/domain"
	"time"
)

type Repository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByCredentials(ctx context.Context, login, password string) (*domain.User, error)
}

type repository struct {
	db *sql.DB
}

var _ Repository = (*repository)(nil)

func NewRepository(db *sql.DB) *repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, user *domain.User) error {
	createStmt := "insert into users(login, password, admin)" +
		"values ($1, $2, $3) returning id, created_at"

	return r.db.QueryRowContext(ctx, createStmt, user.Login, user.Password, user.Admin, time.Now()).
		Scan(&user.Id, &user.RegisteredAt)
}

func (r *repository) GetByCredentials(ctx context.Context, login, password string) (*domain.User, error) {
	var user domain.User
	selectStmt := `select u.id
                         ,u.login
                         ,u.password
                         ,u.admin                    
                         ,u.registered_at
			   	     from users u
                    where u.login = $1
                      and u.password = $2`

	if err := r.db.QueryRowContext(ctx, selectStmt, login, password).Scan(&user.Id, &user.Login, &user.Password, &user.Admin, &user.RegisteredAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}
