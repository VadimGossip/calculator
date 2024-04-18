package app

import (
	"database/sql"
	"fmt"
	"github.com/VadimGossip/calculator/dbagent/internal/domain"
	"github.com/VadimGossip/calculator/dbagent/internal/expression"
	"github.com/VadimGossip/calculator/dbagent/internal/token"
	"github.com/VadimGossip/calculator/dbagent/internal/user"
	_ "github.com/lib/pq"
)

type DBAdapter struct {
	cfg domain.DbCfg
	db  *sql.DB

	exprRepo  expression.Repository
	userRepo  user.Repository
	tokenRepo token.Repository
}

func NewDBAdapter(cfg domain.DbCfg) *DBAdapter {
	return &DBAdapter{cfg: cfg}
}

func (d *DBAdapter) Connect() error {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		d.cfg.Host, d.cfg.Port, d.cfg.Username, d.cfg.Name, d.cfg.SSLMode, d.cfg.Password))
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}

	d.db = db
	d.exprRepo = expression.NewRepository(db)
	d.userRepo = user.NewRepository(db)
	d.tokenRepo = token.NewRepository(db)
	return nil
}

func (d *DBAdapter) Close() error {
	return d.db.Close()
}
