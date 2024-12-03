package sso

import (
	"context"
)

type SSOProvider interface {
	Register(ctx context.Context, email string, password string) (int64, error)
	Login(ctx context.Context, email string, password string) (string, error)
	Validate(ctx context.Context, token string) (int64, error)
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
