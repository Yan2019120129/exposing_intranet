package transport

import (
	"errors"
	"net"
	"syscall"
	"testing"
	"time"

	"my-base/code/message"
	toolsPenetrate "my-base/code/transport"
)

func TestClientRegistersAndAnswersPing(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if errors.Is(err, syscall.EPERM) {
			t.Skipf("sandbox does not permit local listeners: %v", err)
		}
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()

	serverMessages := make(chan message.Message, 2)
	serverErrors := make(chan error, 1)
	go func() {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			serverErrors <- acceptErr
			return
		}
		defer conn.Close()
		control := toolsPenetrate.NewConn(conn)
		var register message.Message
		if parseErr := control.ParseMsg(&register); parseErr != nil {
			serverErrors <- parseErr
			return
		}
		serverMessages <- register
		if sendErr := control.Send(message.Message{Type: message.MsgTypeRegister, Symbol: "assigned-symbol"}); sendErr != nil {
			serverErrors <- sendErr
			return
		}
		if sendErr := control.Send(message.Message{
			Type: message.MsgTypePing,
			Msg:  message.PingPayload{Seq: 7, SentAtMs: time.Now().UnixMilli()},
		}); sendErr != nil {
			serverErrors <- sendErr
			return
		}
		var pong message.Message
		if parseErr := control.ParseMsg(&pong); parseErr != nil {
			serverErrors <- parseErr
			return
		}
		serverMessages <- pong
	}()

	client := NewClient(listener.Addr().String())
	client.SetKeyFunc(func() string { return "stored-symbol" })
	var assigned string
	client.SetOnKeyFunc(func(symbol string) { assigned = symbol })
	done := make(chan error, 1)
	go func() { done <- client.Start() }()

	select {
	case err := <-serverErrors:
		t.Fatalf("fake server: %v", err)
	case register := <-serverMessages:
		if !register.EqRegister() || register.Symbol != "stored-symbol" {
			t.Fatalf("unexpected register message: %+v", register)
		}
	}

	select {
	case err := <-serverErrors:
		t.Fatalf("fake server: %v", err)
	case pong := <-serverMessages:
		if !pong.EqPong() {
			t.Fatalf("expected pong, got %+v", pong)
		}
	}

	_ = client.Close()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("client did not stop after close")
	}
	if assigned != "assigned-symbol" {
		t.Fatalf("assigned symbol = %q, want assigned-symbol", assigned)
	}
}

func TestClientStopsOnRegistrationRejection(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if errors.Is(err, syscall.EPERM) {
			t.Skipf("sandbox does not permit local listeners: %v", err)
		}
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()

	serverDone := make(chan error, 1)
	go func() {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			serverDone <- acceptErr
			return
		}
		defer conn.Close()
		control := toolsPenetrate.NewConn(conn)
		var register message.Message
		if parseErr := control.ParseMsg(&register); parseErr != nil {
			serverDone <- parseErr
			return
		}
		serverDone <- control.Send(message.Message{
			Type: message.MsgTypeClose,
			Msg:  "Invalid symbol.",
		})
	}()

	client := NewClient(listener.Addr().String())
	client.SetKeyFunc(func() string { return "invalid-symbol" })
	err = client.Start()
	if !errors.Is(err, ErrClientRegistrationRejected) {
		t.Fatalf("Start() error = %v, want ErrClientRegistrationRejected", err)
	}
	if serverErr := <-serverDone; serverErr != nil {
		t.Fatalf("fake server: %v", serverErr)
	}
}

func TestCopyClosesIdleMappingAfterDataTimeout(t *testing.T) {
	targetConn, targetPeer := net.Pipe()
	defer targetPeer.Close()
	mappingConn, mappingPeer := net.Pipe()
	defer mappingPeer.Close()

	client := NewClient("").SetDataTimeout(50 * time.Millisecond)
	done := make(chan struct{})
	go func() {
		client.Copy(targetConn, mappingConn)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Copy did not close an idle mapping after data timeout")
	}
}
