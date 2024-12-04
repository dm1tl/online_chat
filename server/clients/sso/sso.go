package sso

import (
	"context"
	appmodels "server/internal/app_models"
)

type SSOProvider interface {
	Register(ctx context.Context, req appmodels.CreateUserReq) (*appmodels.CreateUserResp, error)
	Login(ctx context.Context, req appmodels.LoginReq) (*appmodels.LoginResp, error)
	Validate(ctx context.Context, req appmodels.ValidateTokenReq) (*appmodels.ValidateTokenResp, error)
	Delete(ctx context.Context, id int64) error
}

type SSOClientWrapper struct {
	SSOProvider
}

func NewSSOClientWrapper(provider SSOProvider) *SSOClientWrapper {
	return &SSOClientWrapper{
		SSOProvider: provider,
	}
}
