package services

import (
	"context"
	"server/clients/sso"
	appmodels "server/internal/app_models"
	"server/internal/repository"
)

type UserService struct {
	repo      repository.UserManager
	ssoClient sso.SSOProvider
}

func NewUserService(repo repository.UserManager, ssoCl sso.SSOProvider) *UserService {
	return &UserService{
		repo:      repo,
		ssoClient: ssoCl,
	}
}

func (u *UserService) Create(ctx context.Context, req appmodels.CreateUserReq) error {
	return nil
}

func (u *UserService) Login(ctx context.Context, req appmodels.LoginReq) (appmodels.LoginResp, error) {
	return appmodels.LoginResp{}, nil
}

func (U *UserService) Validate(ctx context.Context, req appmodels.ValidateTokenReq) (appmodels.ValidateTokenResp, error) {
	return appmodels.ValidateTokenResp{}, nil
}
