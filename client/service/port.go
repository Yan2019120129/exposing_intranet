package service

import (
	"context"
	"fmt"

	"my-base/client/apis"
	"my-base/client/configs"
	"my-base/code/contract"
)

// PortService implements the client port management use case.
type PortService struct {
	API    apis.PortManager
	Config *configs.Config
}

func NewPortService(api apis.PortManager, cfg *configs.Config) *PortService {
	return &PortService{API: api, Config: cfg}
}

func (s *PortService) Manage(ctx context.Context, action, serverPort, localAddr, comment string) (contract.PortResponse, error) {
	symbol := s.Config.GetSymbol()
	if symbol == "" {
		return contract.PortResponse{}, fmt.Errorf("客户端未认证，请先运行: client auth <username:password>")
	}
	return s.API.ManagePort(ctx, contract.PortRequest{
		Symbol: symbol, Action: action, ServerPort: serverPort, LocalAddr: localAddr, Comment: comment,
	})
}
