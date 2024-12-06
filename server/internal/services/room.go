package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	appmodels "server/internal/app_models"
	"server/internal/repository"
	"time"
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

func (r *RoomService) AddClient(ctx context.Context, req appmodels.AddClientReq) error {
	err := r.repo.RoomManager.AddClient(ctx, req)
	if err != nil {
		return err
	}
	return nil
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

func (r *RoomService) AddMessage(ctx context.Context, req appmodels.AddMessageReq) error {
	const op = "service.AddMessage"
	ctx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
	defer cancel()
	err := r.repo.RoomManager.AddMessage(ctx, req)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
