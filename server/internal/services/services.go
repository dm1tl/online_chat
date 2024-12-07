package services

import (
	"context"
	"server/clients/sso"
	appmodels "server/internal/app_models"
	"server/internal/repository"
)

type Service struct {
	UserManager
	RoomManager
}

func NewService(repo *repository.Repository, ssoclient *sso.SSOClientWrapper) *Service {
	return &Service{
		UserManager: NewUserService(ssoclient),
		RoomManager: NewRoomService(repo),
	}
}

type UserManager interface {
	Create(ctx context.Context, req appmodels.CreateUserReq) error
	Login(ctx context.Context, req appmodels.LoginReq) (appmodels.LoginResp, error)
	Validate(ctx context.Context, req appmodels.ValidateTokenReq) (appmodels.ValidateTokenResp, error)
}

type RoomManager interface {
	CreateRoom(ctx context.Context, req appmodels.CreateRoomReq) (int64, error)
	AddClient(ctx context.Context, req appmodels.AddClientReq) error
	GetRoom(ctx context.Context, req appmodels.AddClientReq) (bool, error)
	AddMessage(ctx context.Context, req appmodels.AddMessageReq) error
}
