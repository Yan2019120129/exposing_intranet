package transport

import (
	"errors"
	"net"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"my-base/code/message"
	toolsPenetrate "my-base/code/transport"
)

func TestHeartbeatTimeoutMatchesFailureSchedule(t *testing.T) {
	got := heartbeatTimeout(30*time.Second, 90*time.Second, 3)
	if want := 150 * time.Second; got != want {
		t.Fatalf("heartbeat timeout = %v, want %v", got, want)
	}
}

func TestHeartbeatWatcherClosesOnlyItsSession(t *testing.T) {
	oldConn, oldPeer := net.Pipe()
	defer oldPeer.Close()
	newConn, newPeer := net.Pipe()
	defer newPeer.Close()

	oldControl := toolsPenetrate.NewConnWithOptions(oldConn, ConnOptions{})
	newControl := toolsPenetrate.NewConnWithOptions(newConn, ConnOptions{})
	client := NewClient("")
	client.conn = newControl

	lastSeen := &atomic.Int64{}
	lastSeen.Store(time.Now().Add(-time.Second).UnixNano())
	done := make(chan struct{})
	go func() {
		client.watchHeartbeat(oldControl, lastSeen, 5*time.Millisecond, 10*time.Millisecond, 1, make(chan struct{}))
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("heartbeat watcher did not stop")
	}
	if !oldControl.IsClose() {
		t.Fatal("watcher did not close its own session")
	}
	if newControl.IsClose() || client.currentConn() != newControl {
		t.Fatal("stale watcher closed or replaced the current session")
	}
}

func TestClientHeartbeatClosesSilentControlConnection(t *testing.T) {
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
		var next message.Message
		serverDone <- control.ParseMsg(&next)
	}()

	client := NewClient(listener.Addr().String()).SetKeyFunc(func() string { return "silent" })
	client.SetConnOptions(ConnOptions{ReadTimeout: 0, WriteTimeout: time.Second})
	started := time.Now()
	err = client.StartWithHeartbeat(10*time.Millisecond, 20*time.Millisecond, 2)
	if err == nil {
		t.Fatal("expected silent control connection to be closed")
	}
	if elapsed := time.Since(started); elapsed > 500*time.Millisecond {
		t.Fatalf("heartbeat watcher took too long: %v", elapsed)
	}
	select {
	case <-serverDone:
	case <-time.After(time.Second):
		t.Fatal("server read was not released by heartbeat close")
	}
}

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
