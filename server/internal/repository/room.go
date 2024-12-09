package repository

import (
	"context"
	"database/sql"
	"fmt"
	appmodels "server/internal/app_models"
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

func (r *RoomRepository) GetAllRooms(ctx context.Context) ([]appmodels.BackupRoom, error) {
	op := "repository.GetAllRooms"
	var output []appmodels.BackupRoom
	query := "SELECT id, name, password FROM rooms"
	res, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer res.Close()
	for res.Next() {
		var room appmodels.BackupRoom
		if err := res.Scan(&room.ID, &room.Name, &room.Password); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		output = append(output, room)
	}
	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return output, nil
}
