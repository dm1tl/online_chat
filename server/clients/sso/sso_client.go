package sso

import (
	"context"
	"errors"
	"fmt"
	appmodels "server/internal/app_models"
	"server/internal/config"
	"time"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/sirupsen/logrus"

	ssov1 "github.com/dm1tl/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type SSOServiceCLient struct {
	authAPI ssov1.AuthClient
	userAPI ssov1.UserClient
}

func NewSSOServiceClient(log *logrus.Logger, cfg config.SSOConfig) (*SSOServiceCLient, error) {
	const op = "clients.sso.grpc.New()"
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(cfg.RetriesCount),
		grpcretry.WithPerRetryTimeout(cfg.Timeout),
	}
	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}
	cc, err := grpc.NewClient(cfg.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &SSOServiceCLient{
		authAPI: ssov1.NewAuthClient(cc),
		userAPI: ssov1.NewUserClient(cc),
	}, nil
}

func InterceptorLogger(l *logrus.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
		l.Log(logrus.Level(level), msg)
	})
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
	req appmodels.CreateUserReq) error {
	const op = "clients.sso.grpc.Register()"
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := c.authAPI.Register(ctx, &ssov1.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
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
