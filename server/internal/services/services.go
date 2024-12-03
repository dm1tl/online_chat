package services

import (
	"context"
	"server/clients/sso"
	appmodels "server/internal/app_models"
	"server/internal/repository"
)

type Service struct {
	UserManager
}

func NewService(repo *repository.Repository, ssoclient *sso.SSOClientWrapper) *Service {
	return &Service{
		UserManager: NewUserService(repo, ssoclient),
	}
}

type UserManager interface {
	Create(ctx context.Context, req appmodels.CreateUserReq) error
	Login(ctx context.Context, req appmodels.LoginReq) (appmodels.LoginResp, error)
	Validate(ctx context.Context, req appmodels.ValidateTokenReq) (appmodels.ValidateTokenResp, error)
}
