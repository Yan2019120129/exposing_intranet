package penetrate

import (
	"net"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"my-base/app/models"
	"my-base/code/message"
	transport "my-base/code/transport"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestDelClientFForcesClosedResourcesWhenPeerDoesNotRead(t *testing.T) {
	controlConn, controlPeer := net.Pipe()
	defer controlPeer.Close()
	client := NewClient(transport.NewConnWithOptions(controlConn, transport.ConnOptions{})).SetSymbol("deleted")

	listen := NewListen("127.0.0.1:0", ":8080")
	if err := listen.Bind(); err != nil {
		t.Fatalf("bind listener: %v", err)
	}
	activeConn, activePeer := net.Pipe()
	defer activePeer.Close()
	listen.AddConn("active", transport.NewConn(activeConn))
	client.AddListen(":0", listen)

	server := &Server{
		client: map[string]*Client{"deleted": client},
		status: map[string]ClientStatus{"deleted": {Symbol: "deleted"}},
	}
	started := time.Now()
	err := server.DelClientF("deleted")
	if err == nil {
		t.Fatal("expected the unread deletion notification to time out")
	}
	if elapsed := time.Since(started); elapsed > 1500*time.Millisecond {
		t.Fatalf("forced deletion exceeded its bounded notification timeout: %v", elapsed)
	}
	if server.IsExist("deleted") {
		t.Fatal("deleted client remains in the online map")
	}
	if _, ok := server.GetClientStatus("deleted"); ok {
		t.Fatal("deleted client status remains cached")
	}

	client.lock.RLock()
	closed := client.closed
	remaining := len(client.Listens)
	client.lock.RUnlock()
	if !closed || remaining != 0 {
		t.Fatalf("client resources not closed: closed=%v listeners=%d", closed, remaining)
	}
	listen.lock.RLock()
	stopped := listen.stopped
	listen.lock.RUnlock()
	if !stopped {
		t.Fatal("deleted client's listener was not stopped")
	}
	if err := activePeer.SetReadDeadline(time.Now().Add(100 * time.Millisecond)); err == nil {
		if _, readErr := activePeer.Read(make([]byte, 1)); readErr == nil {
			t.Fatal("active tunnel remained open after client deletion")
		}
	}
}

func TestRegisterRevalidationRejectsClientDeletedInFlight(t *testing.T) {
	dsn := filepath.Join(t.TempDir(), "register.db") + "?_busy_timeout=5000&_journal_mode=WAL"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(&models.Client{}, &models.Port{}); err != nil {
		t.Fatalf("migrate database: %v", err)
	}
	clientRow := models.Client{Name: "in-flight", Symbol: "in-flight-symbol", Status: models.StatusOn}
	if err := db.Create(&clientRow).Error; err != nil {
		t.Fatalf("create client row: %v", err)
	}

	var clientQueries atomic.Int32
	if err := db.Callback().Query().Before("gorm:query").Register("test:delete_before_revalidation", func(tx *gorm.DB) {
		if tx.Statement.Table != "client" || clientQueries.Add(1) != 2 {
			return
		}
		if deleteErr := db.Exec("DELETE FROM client WHERE symbol = ?", clientRow.Symbol).Error; deleteErr != nil {
			t.Errorf("delete client during registration: %v", deleteErr)
		}
	}); err != nil {
		t.Fatalf("register query callback: %v", err)
	}

	serverConn, peerConn := net.Pipe()
	defer peerConn.Close()
	server := &Server{
		client: make(map[string]*Client),
		status: make(map[string]ClientStatus),
		db:     db,
	}
	received := make(chan message.Message, 1)
	go func() {
		var msg message.Message
		_ = transport.NewConn(peerConn).ParseMsg(&msg)
		received <- msg
	}()

	server.Register(message.Message{Symbol: clientRow.Symbol, Type: message.MsgTypeRegister}, transport.NewConnWithOptions(serverConn, transport.ConnOptions{}))
	if server.IsExist(clientRow.Symbol) {
		t.Fatal("client deleted during registration remained online")
	}
	select {
	case msg := <-received:
		if !msg.EqDel() {
			t.Fatalf("registration response type = %v, want delete", msg.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("deleted in-flight registration did not receive a response")
	}
}
