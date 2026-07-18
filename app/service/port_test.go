package service

import (
	"testing"

	"my-base/app/models"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type fakePortRuntime struct {
	online    map[string]bool
	listeners map[string]bool
}

func newFakePortRuntime() *fakePortRuntime {
	return &fakePortRuntime{online: make(map[string]bool), listeners: make(map[string]bool)}
}

func (f *fakePortRuntime) IsExist(symbol string) bool { return f.online[symbol] }

func (f *fakePortRuntime) IsListenExist(symbol, serverPort string) bool {
	return f.listeners[symbol+serverPort]
}

func (f *fakePortRuntime) NewListen(symbol, _, serverPort string) error {
	f.listeners[symbol+serverPort] = true
	return nil
}

func (f *fakePortRuntime) CloseListen(symbol, _, serverPort string) error {
	delete(f.listeners, symbol+serverPort)
	return nil
}

func openPortServiceDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+uuid.NewString()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.Client{}, &models.Port{}); err != nil {
		t.Fatalf("migrate sqlite: %v", err)
	}
	return db
}

func TestPortServiceManageWithInjectedRuntime(t *testing.T) {
	db := openPortServiceDB(t)
	client := models.Client{Name: "test", Symbol: "client-1", Status: models.StatusActive}
	if err := db.Create(&client).Error; err != nil {
		t.Fatalf("create client: %v", err)
	}

	runtime := newFakePortRuntime()
	runtime.online[client.Symbol] = true
	portService := NewPortService(db, runtime)

	add, err := portService.Manage(PortCommand{
		Symbol:     client.Symbol,
		Action:     "add",
		ServerPort: "12345",
		LocalAddr:  "127.0.0.1:23456",
		Comment:    "test",
	})
	if err != nil {
		t.Fatalf("add mapping: %v", err)
	}
	if len(add.Mappings) != 1 || add.Mappings[0].Status != "active" {
		t.Fatalf("unexpected add result: %+v", add)
	}

	list, err := portService.Manage(PortCommand{Symbol: client.Symbol, Action: "list"})
	if err != nil || len(list.Mappings) != 1 || list.Mappings[0].Status != "active" {
		t.Fatalf("unexpected list result: %+v, %v", list, err)
	}

	if _, err := portService.Manage(PortCommand{
		Symbol:     client.Symbol,
		Action:     "add",
		ServerPort: "12345",
		LocalAddr:  "127.0.0.1:23457",
	}); err != ErrPortConflict {
		t.Fatalf("expected port conflict, got %v", err)
	}

	if _, err := portService.Manage(PortCommand{Symbol: client.Symbol, Action: "del", ServerPort: "12345"}); err != nil {
		t.Fatalf("delete mapping: %v", err)
	}
	if runtime.listeners[client.Symbol+":12345"] {
		t.Fatal("listener should be closed after delete")
	}
}
