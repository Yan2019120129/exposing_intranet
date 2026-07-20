package service

import (
	"context"
	"errors"
	"net"
	"syscall"
	"testing"
	"time"

	"my-base/client/configs"
	toolsPenetrate "my-base/code/transport"
)

func TestRuntimeCancellationClosesBlockingControlRead(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if errors.Is(err, syscall.EPERM) {
			t.Skipf("sandbox does not permit local listeners: %v", err)
		}
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()

	connected := make(chan struct{})
	serverDone := make(chan struct{})
	go func() {
		defer close(serverDone)
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			return
		}
		defer conn.Close()
		control := toolsPenetrate.NewConn(conn)
		var register any
		if control.ParseMsg(&register) == nil {
			close(connected)
			_, _ = control.GetConn().Read(make([]byte, 1))
		}
	}()

	host, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatalf("split listener address: %v", err)
	}
	runtime := NewRuntime(&configs.Config{Client: &configs.ClientConfig{
		ServerAddr: host,
		JobPort:    port,
	}})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- runtime.Run(ctx) }()

	select {
	case <-connected:
	case <-time.After(2 * time.Second):
		t.Fatal("client did not establish the control connection")
	}
	cancel()

	select {
	case runErr := <-done:
		if runErr != nil {
			t.Fatalf("Run() error = %v, want nil", runErr)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not stop after context cancellation")
	}
	<-serverDone
}
