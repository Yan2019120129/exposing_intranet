package service

import (
	"context"
	"testing"

	"my-base/client/configs"
	"my-base/code/contract"
)

type fakeAuthenticator struct {
	request contract.AuthRequest
}

func (f *fakeAuthenticator) Authenticate(_ context.Context, request contract.AuthRequest) (contract.AuthResponse, error) {
	f.request = request
	return contract.AuthResponse{Success: true, Symbol: "test-symbol"}, nil
}

func TestParseAuthArg(t *testing.T) {
	username, password, err := ParseAuthArg("admin:p:a:ss")
	if err != nil {
		t.Fatalf("parse auth argument: %v", err)
	}
	if username != "admin" || password != "p:a:ss" {
		t.Fatalf("unexpected auth argument: %q:%q", username, password)
	}
	if _, _, err := ParseAuthArg("invalid"); err == nil {
		t.Fatal("expected invalid auth argument error")
	}
}

func TestAuthServicePersistsSymbol(t *testing.T) {
	t.Setenv("EXPOSING_INTRANET_SYMBOL_PATH", t.TempDir()+"/public/key.init")
	fake := &fakeAuthenticator{}
	service := NewAuthService(fake, &configs.Config{Client: &configs.ClientConfig{}})

	symbol, err := service.Authenticate(context.Background(), "admin:secret")
	if err != nil {
		t.Fatalf("authenticate: %v", err)
	}
	if symbol != "test-symbol" || fake.request.Username != "admin" || fake.request.Password != "secret" {
		t.Fatalf("unexpected authentication result: symbol=%q request=%+v", symbol, fake.request)
	}
	if got := service.Config.GetSymbol(); got != symbol {
		t.Fatalf("persisted symbol = %q, want %q", got, symbol)
	}
}
