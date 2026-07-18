package apis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"my-base/client/configs"
	"my-base/code/contract"
)

// Authenticator is the client-side outbound authentication boundary.
type Authenticator interface {
	Authenticate(context.Context, contract.AuthRequest) (contract.AuthResponse, error)
}

// HTTPAPI calls the server's public client management endpoints.
type HTTPAPI struct {
	Config *configs.ClientConfig
	Client *http.Client
}

func NewHTTPAPI(cfg *configs.ClientConfig) *HTTPAPI {
	return &HTTPAPI{Config: cfg, Client: http.DefaultClient}
}

func (a *HTTPAPI) Authenticate(ctx context.Context, req contract.AuthRequest) (contract.AuthResponse, error) {
	var result contract.AuthResponse
	if err := a.post(ctx, "/api/client/register", req, &result); err != nil {
		return result, err
	}
	if !result.Success {
		return result, fmt.Errorf("认证失败: %s", result.Message)
	}
	return result, nil
}

func (a *HTTPAPI) post(ctx context.Context, path string, request any, response any) error {
	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("编码请求失败: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.url(path), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := a.Client
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("连接服务端失败: %w", err)
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}
	return nil
}

func (a *HTTPAPI) url(path string) string {
	return fmt.Sprintf("http://%s:%s%s", a.Config.GetServerAddr(), a.Config.GetPort(), path)
}
