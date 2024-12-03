package sso

import (
	"context"
	"server/internal/config"
)

type MockSSOServiceClient struct{}

func NewMockSSOServiceClient(cfg config.SSOConfig) *MockSSOServiceClient {
	return &MockSSOServiceClient{}
}

func (m *MockSSOServiceClient) Register(ctx context.Context, email string, password string) (int64, error) {
	return 1, nil
}

func (m *MockSSOServiceClient) Login(ctx context.Context, email string, password string) (string, error) {
	return "mock-tocken", nil
}

func (m *MockSSOServiceClient) Validate(ctx context.Context, token string) (int64, error) {
	return 1, nil
}

func (m *MockSSOServiceClient) Delete(ctx context.Context, id int64) error {
	return nil
}
