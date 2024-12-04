package sso

import (
	"context"
	"errors"
	"fmt"
	appmodels "server/internal/app_models"
	"time"

	ssov1 "github.com/dm1tl/protos/gen/go/sso"
)

type SSOServiceCLient struct {
	authAPI ssov1.AuthClient
	userAPI ssov1.UserClient
}

func (c *SSOServiceCLient) Login(ctx context.Context,
	req appmodels.LoginReq) (*appmodels.LoginResp, error) {
	const op = "clients.sso.grpc.Login()"
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	resp, err := c.authAPI.Login(ctx, &ssov1.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return &appmodels.LoginResp{}, fmt.Errorf("%s: %w", op, err)
	}
	return &appmodels.LoginResp{
		Token: resp.Token,
	}, nil
}

func (c *SSOServiceCLient) Register(ctx context.Context,
	req appmodels.CreateUserReq) (*appmodels.CreateUserResp, error) {
	const op = "clients.sso.grpc.Register()"
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	resp, err := c.authAPI.Register(ctx, &ssov1.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return &appmodels.CreateUserResp{}, fmt.Errorf("%s: %w", op, err)
	}
	return &appmodels.CreateUserResp{
		ID:       resp.UserId,
		Username: req.Username,
	}, nil
}

func (c *SSOServiceCLient) Validate(ctx context.Context,
	req appmodels.ValidateTokenReq) (*appmodels.ValidateTokenResp, error) {
	const op = "clients.sso.grpc.ValidateToken()"
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	resp, err := c.authAPI.ValidateToken(ctx, &ssov1.ValidateTokenRequest{
		Token: req.Token,
	})
	if err != nil {
		return &appmodels.ValidateTokenResp{}, fmt.Errorf("%s: %w", op, err)
	}
	return &appmodels.ValidateTokenResp{
		ID: resp.Id,
	}, nil
}

func (c *SSOServiceCLient) Delete(ctx context.Context,
	id int64) (err error) {
	const op = "clients.sso.grpc.Delete()"
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	resp, err := c.userAPI.Delete(ctx, &ssov1.DeleteRequest{
		Id: id,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if resp.ErrorMessage != "success" {
		return errors.New("couldn't delete user from grpc db")
	}
	return nil
}
