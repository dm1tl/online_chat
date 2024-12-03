package connector

import (
	"context"
	"server/internal/config"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sqlx.DB
}

func NewDatabase(cfg config.DBConfig) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := sqlx.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return &Database{
		db: db,
	}, nil
}

func (d *Database) Close() error {
	if err := d.db.Close(); err != nil {
		return err
	}
	return nil
}

func (d *Database) GetDB() *sqlx.DB {
	return d.db
}
