package repository

import (
	"context"
	"errors"
	"fmt"
	appmodels "server/internal/app_models"
)

type UserRepository struct {
	db DBTX
}

func NewUserRepository(db DBTX) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Create(ctx context.Context, req appmodels.CreateUserResp) error {
	const op = "internal.repository.Create()"
	query := "INSERT INTO users (id, username) VALUES ($1, $2)"
	row, err := u.db.ExecContext(ctx, query, req.ID, req.Username)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	res, err := row.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if res == 0 {
		return errors.New("couldn't input user's data")
	}
	return nil
}
