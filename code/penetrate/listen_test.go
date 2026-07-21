package penetrate

import (
	"errors"
	"my-base/code/message"
	toolsPenetrate "my-base/code/transport"
	"net"
	"sync/atomic"
	"testing"
	"time"
)

var errAcceptFailed = errors.New("accept failed")

type failingListener struct{ closed atomic.Bool }

func (l *failingListener) Accept() (net.Conn, error) { return nil, errAcceptFailed }
func (l *failingListener) Close() error {
	l.closed.Store(true)
	return nil
}
func (l *failingListener) Addr() net.Addr { return &net.TCPAddr{} }

func TestListenBindsSynchronouslyAndStopsIdempotently(t *testing.T) {
	listen := NewListen("127.0.0.1:0", "127.0.0.1:1")
	if err := listen.Bind(); err != nil {
		t.Fatalf("bind: %v", err)
	}
	if listen.listen == nil {
		t.Fatal("listener was not retained after successful bind")
	}
	if err := listen.Stop(); err != nil {
		t.Fatalf("first stop: %v", err)
	}
	if err := listen.Stop(); err != nil {
		t.Fatalf("second stop: %v", err)
	}
}

func TestServeStopsListenerAfterUnexpectedAcceptError(t *testing.T) {
	failing := &failingListener{}
	listen := NewListen("127.0.0.1:0", "127.0.0.1:1")
	listen.listen = failing

	if err := listen.Serve(); !errors.Is(err, errAcceptFailed) {
		t.Fatalf("Serve() error = %v, want %v", err, errAcceptFailed)
	}
	if !failing.closed.Load() {
		t.Fatal("Serve did not close listener after unexpected Accept error")
	}
}

func TestStopRejectsConnectionHandledAfterShutdown(t *testing.T) {
	listen := NewListen("127.0.0.1:0", "127.0.0.1:1")
	public, peer := net.Pipe()
	defer peer.Close()

	if err := listen.Stop(); err != nil {
		t.Fatalf("Stop: %v", err)
	}
	listen.dealWith(public)

	listen.lock.RLock()
	remaining := len(listen.connList)
	listen.lock.RUnlock()
	if remaining != 0 {
		t.Fatalf("late connection remained tracked after Stop: %d", remaining)
	}
	if err := peer.SetReadDeadline(time.Now().Add(time.Second)); err == nil {
		if _, err := peer.Read(make([]byte, 1)); err == nil {
			t.Fatal("late connection remained open after Stop")
		}
	}
}

func TestClosedClientStopsLateRegisteredListener(t *testing.T) {
	client := NewClient(nil)
	if err := client.Close(); err != nil {
		t.Fatalf("client Close: %v", err)
	}
	listen := NewListen("127.0.0.1:0", "127.0.0.1:1")
	if err := listen.Bind(); err != nil {
		t.Fatalf("listen Bind: %v", err)
	}
	client.AddListen("late", listen)

	listen.lock.RLock()
	stopped := listen.stopped
	listener := listen.listen
	listen.lock.RUnlock()
	if !stopped || listener != nil {
		t.Fatal("late listener remained live after client shutdown")
	}
}

func TestListenExpiresUnmatchedConnections(t *testing.T) {
	listen := NewListen("127.0.0.1:0", "127.0.0.1:1").SetTimeouts(30*time.Millisecond, 5*time.Millisecond, 0)
	if err := listen.Bind(); err != nil {
		t.Fatalf("bind: %v", err)
	}
	defer listen.Stop()

	address := listen.listen.Addr().String()
	public, err := net.Dial("tcp", address)
	if err != nil {
		t.Fatalf("dial public listener: %v", err)
	}
	defer public.Close()

	listen.Notify = func(message.Message) error { return nil }
	go listen.Serve()

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		listen.lock.RLock()
		remaining := len(listen.connList)
		listen.lock.RUnlock()
		if remaining == 0 {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("unmatched connection was not reaped")
}

func TestListenKeepsActivatedConnections(t *testing.T) {
	listen := NewListen("127.0.0.1:0", "127.0.0.1:1").SetTimeouts(30*time.Millisecond, 5*time.Millisecond, 0)
	public, peer := net.Pipe()
	defer peer.Close()
	conn := toolsPenetrate.NewConn(public).SetStatusWait()
	listen.AddConn("waiting", conn)
	listen.lock.Lock()
	listen.waitSince["waiting"] = time.Now()
	listen.lock.Unlock()

	if listen.ActivateConn("waiting") != conn {
		t.Fatal("waiting connection was not activated")
	}
	go listen.reapWaiters()
	defer listen.Stop()
	time.Sleep(100 * time.Millisecond)
	if listen.GetConn("waiting") != conn {
		t.Fatal("activated connection was reaped")
	}
}
