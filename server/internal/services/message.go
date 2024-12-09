package services

import (
	"context"
	"fmt"
	appmodels "server/internal/app_models"
	"server/internal/repository"
	"time"
)

type MessageService struct {
	repo *repository.Repository
}

func NewMessageService(repo *repository.Repository) *MessageService {
	return &MessageService{
		repo: repo,
	}
}

func (r *MessageService) AddMessage(ctx context.Context, req appmodels.AddMessageReq) error {
	const op = "service.AddMessage"
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := r.repo.MessageManager.AddMessage(ctx, req)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
