package services

import (
	"context"
	"fmt"
	"server/clients/sso"
	appmodels "server/internal/app_models"
)

type UserService struct {
	ssoClient sso.SSOProvider
}

func NewUserService(ssoCl sso.SSOProvider) *UserService {
	return &UserService{
		ssoClient: ssoCl,
	}
}

func (u *UserService) Create(ctx context.Context, req appmodels.CreateUserReq) error {
	const op = "internal.services.Create()"
	err := u.ssoClient.Register(ctx, req)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (u *UserService) Login(ctx context.Context, req appmodels.LoginReq) (appmodels.LoginResp, error) {
	var resp appmodels.LoginResp
	const op = "internal.services.Login()"
	ssoResp, err := u.ssoClient.Login(ctx, req)
	if err != nil {
		return resp, fmt.Errorf("%s: %w", op, err)
	}
	return appmodels.LoginResp{
		Token: ssoResp.Token,
	}, nil
}

func (u *UserService) Validate(ctx context.Context, req appmodels.ValidateTokenReq) (appmodels.ValidateTokenResp, error) {
	var resp appmodels.ValidateTokenResp
	const op = "internal.service.Validate()"
	ssoResp, err := u.ssoClient.Validate(ctx, req)
	if err != nil {
		return resp, fmt.Errorf("%s: %w", op, err)
	}
	return appmodels.ValidateTokenResp{
		ID: ssoResp.ID,
	}, nil
}
