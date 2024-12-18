package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	appmodels "server/internal/app_models"
	"server/internal/repository"
)

type RoomService struct {
	repo *repository.Repository
}

func NewRoomService(repo *repository.Repository) *RoomService {
	return &RoomService{
		repo: repo,
	}
}

func (r *RoomService) CreateRoom(ctx context.Context, req appmodels.CreateRoomReq) (int64, error) {
	id, err := r.repo.RoomManager.CreateRoom(ctx, req)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *RoomService) GetRoom(ctx context.Context, req appmodels.AddClientReq) (bool, error) {
	const op = "service.GetRoom"

	resp, err := r.repo.RoomManager.GetRoom(ctx, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: room not found: %w", op, err)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if resp.Password != req.Password {
		return false, nil
	}

	return true, nil
}

func (r *RoomService) GetAllRooms(ctx context.Context) ([]appmodels.BackupRoom, error) {
	data, err := r.repo.RoomManager.GetAllRooms(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}
