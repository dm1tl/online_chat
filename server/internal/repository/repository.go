package repository

import (
	"context"
	"database/sql"
	appmodels "server/internal/app_models"
)

type Repository struct {
	UserManager
}

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func NewRepository(db DBTX) *Repository {
	return &Repository{
		UserManager: NewUserRepository(db),
	}
}

type UserManager interface {
	Create(ctx context.Context, req appmodels.CreateUserResp) error
}
