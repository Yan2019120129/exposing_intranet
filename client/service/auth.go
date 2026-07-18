package service

import (
	"context"
	"fmt"
	"os"
	"strings"

	"my-base/client/apis"
	"my-base/client/configs"
	"my-base/code/contract"
)

// AuthService implements the client authentication use case.
type AuthService struct {
	API    apis.Authenticator
	Config *configs.Config
}

func NewAuthService(api apis.Authenticator, cfg *configs.Config) *AuthService {
	return &AuthService{API: api, Config: cfg}
}

func ParseAuthArg(arg string) (string, string, error) {
	parts := strings.SplitN(arg, ":", 2)
	if len(parts) != 2 || parts[0] == "" {
		return "", "", fmt.Errorf("参数格式错误，应为 username:password")
	}
	return parts[0], parts[1], nil
}

func (s *AuthService) Authenticate(ctx context.Context, authArg string) (string, error) {
	username, password, err := ParseAuthArg(authArg)
	if err != nil {
		return "", err
	}
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	response, err := s.API.Authenticate(ctx, contract.AuthRequest{
		Username: username,
		Password: password,
		Hostname: hostname,
	})
	if err != nil {
		return "", err
	}
	s.Config.SetSymbol(response.Symbol)
	return response.Symbol, nil
}
