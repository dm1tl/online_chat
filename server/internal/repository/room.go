package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	appmodels "server/internal/app_models"

	"github.com/sirupsen/logrus"
)

type RoomRepository struct {
	db DBTX
}

func NewRoomRepository(db DBTX) *RoomRepository {
	return &RoomRepository{
		db: db,
	}
}

func (r *RoomRepository) CreateRoom(ctx context.Context, req appmodels.CreateRoomReq) (int64, error) {
	op := "repository.CreateRoom"
	query := "INSERT INTO rooms (name, password) VALUES ($1, $2) RETURNING id"
	var id int64

	if err := r.db.QueryRowContext(ctx, query, req.Name, req.Password).Scan(&id); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (r *RoomRepository) AddClient(ctx context.Context, req appmodels.AddClientReq) error {
	op := "repository.AddClient"

	clQuery := "INSERT INTO clients (id, username, room_id) VALUES ($1, $2, $3)"
	row, err := r.db.ExecContext(ctx, clQuery, req.ClientID, req.Username, req.RoomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	res, err := row.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if res != 1 {
		return fmt.Errorf("%s: %w", op, errors.New("error while adding client in db"))
	}
	return nil
}

func (r *RoomRepository) GetRoom(ctx context.Context, req appmodels.AddClientReq) (*appmodels.GetRoomResp, error) {
	const op = "repository.GetRoom"
	query := "SELECT id, name, password FROM rooms WHERE id = $1"
	var output appmodels.GetRoomResp
	if err := r.db.QueryRowContext(ctx, query, req.RoomID).Scan(
		&output.ID,
		&output.Name,
		&output.Password,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: room not found: %w", op, err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &output, nil
}

func (r *RoomRepository) AddMessage(ctx context.Context, req appmodels.AddMessageReq) error {
	op := "repository.AddMessage"

	clQuery := "INSERT INTO messages (client_id, room_id, content) VALUES ($1, $2, $3)"
	_, err := r.db.ExecContext(ctx, clQuery, req.UserID, req.RoomID, req.Content)
	if err != nil {
		logrus.Error(op, err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
