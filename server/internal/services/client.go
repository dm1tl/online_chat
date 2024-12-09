package services

import (
	"context"
	appmodels "server/internal/app_models"
	"server/internal/repository"
)

type ClientService struct {
	repo *repository.Repository
}

func NewClientService(repo *repository.Repository) *ClientService {
	return &ClientService{
		repo: repo,
	}
}

func (r *ClientService) AddClient(ctx context.Context, req appmodels.AddClientReq) error {
	err := r.repo.ClientManager.AddClient(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
