package app

import (
	"database/sql"
	"fmt"
	"github.com/VadimGossip/calculator/api/internal/writer"

	_ "github.com/lib/pq"
)

type DBAdapter struct {
	db *sql.DB

	writerRepo writer.Repository
}

func NewDBAdapter() *DBAdapter {
	dba := &DBAdapter{}
	return dba
}

func (d *DBAdapter) Connect() error {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		"postgres", 5432, "postgres", "calculator_db", "disable", "postgres"))
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
