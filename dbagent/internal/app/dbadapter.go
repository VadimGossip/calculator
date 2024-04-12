package app

import (
	"database/sql"
	"fmt"
	"github.com/VadimGossip/calculator/dbagent/internal/domain"
	"github.com/VadimGossip/calculator/dbagent/internal/writer"

	_ "github.com/lib/pq"
)

type DBAdapter struct {
	cfg domain.DbCfg
	db  *sql.DB

	writerRepo writer.Repository
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
	d.writerRepo = writer.NewRepository(db)
	return nil
}

func (d *DBAdapter) Close() error {
	return d.db.Close()
}
