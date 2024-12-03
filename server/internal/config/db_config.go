package config

import (
	"errors"
	"os"
)

const (
	pg_dsn = "PG_DSN"
)

type dbConfig struct {
	dsn string
}

func NewDBConfig() (*dbConfig, error) {
	dsn := os.Getenv(pg_dsn)
	if len(dsn) == 0 {
		return nil, errors.New("pg_dsn is empty")
	}
	return &dbConfig{
		dsn: dsn,
	}, nil
}

func (p *dbConfig) DSN() string {
	return p.dsn
}
