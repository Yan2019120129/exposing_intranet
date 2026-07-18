package penetrate

import (
	"my-base/code/message"
	toolsPenetrate "my-base/code/transport"
	"net"
	"testing"
	"time"
)

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
