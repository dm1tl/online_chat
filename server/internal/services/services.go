package services

import (
	"context"
	"server/clients/sso"
	appmodels "server/internal/app_models"
	"server/internal/repository"
)

type Service struct {
	AuthManager
	RoomManager
	ClientManager
	MessageManager
}

func NewService(repo *repository.Repository, ssoclient *sso.SSOClientWrapper) *Service {
	return &Service{
		AuthManager:    NewUserService(ssoclient),
		RoomManager:    NewRoomService(repo),
		ClientManager:  NewClientService(repo),
		MessageManager: NewMessageService(repo),
	}
}

type AuthManager interface {
	Create(ctx context.Context, req appmodels.CreateUserReq) error
	Login(ctx context.Context, req appmodels.LoginReq) (appmodels.LoginResp, error)
	Validate(ctx context.Context, req appmodels.ValidateTokenReq) (appmodels.ValidateTokenResp, error)
}

type RoomManager interface {
	CreateRoom(ctx context.Context, req appmodels.CreateRoomReq) (int64, error)
	GetRoom(ctx context.Context, req appmodels.AddClientReq) (bool, error)
	GetAllRooms(ctx context.Context) ([]appmodels.BackupRoom, error)
}

type ClientManager interface {
	AddClient(ctx context.Context, req appmodels.AddClientReq) error
}

type MessageManager interface {
	AddMessage(ctx context.Context, req appmodels.AddMessageReq) error
	GetAllMessages(ctx context.Context) (appmodels.BackupMessages, error)
}
