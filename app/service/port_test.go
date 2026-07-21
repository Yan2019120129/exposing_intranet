package service

import (
	"errors"
	"path/filepath"
	"sync"
	"testing"

	"my-base/app/models"
	"my-base/app/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type fakePortRuntime struct {
	online        map[string]bool
	listeners     map[string]bool
	newListenHook func()
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
	if f.newListenHook != nil {
		f.newListenHook()
	}
	return nil
}

func (f *fakePortRuntime) CloseListen(symbol, _, serverPort string) error {
	delete(f.listeners, symbol+serverPort)
	return nil
}

func openPortServiceDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := filepath.Join(t.TempDir(), "port.db") + "?_busy_timeout=5000&_journal_mode=WAL"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.Client{}, &models.Port{}); err != nil {
		t.Fatalf("migrate sqlite: %v", err)
	}
	return db
}

func TestPortModelCreatesRequiredUniqueIndexes(t *testing.T) {
	db := openPortServiceDB(t)
	if !db.Migrator().HasIndex(&models.Port{}, "uidx_port_server") {
		t.Fatal("missing global server-port unique index")
	}
	if !db.Migrator().HasIndex(&models.Port{}, "uidx_port_client_local") {
		t.Fatal("missing client/local unique index")
	}
}

func TestPortRepositoryRejectsConcurrentDuplicateServer(t *testing.T) {
	db := openPortServiceDB(t)
	repo := repository.NewPortRepository(db)
	start := make(chan struct{})
	errs := make(chan error, 2)
	var wg sync.WaitGroup
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()
			<-start
			errs <- repo.Create(&models.Port{ClientId: clientID, Server: ":18080", Local: ":8080"})
		}(i)
	}
	close(start)
	wg.Wait()
	close(errs)

	succeeded := 0
	conflicted := 0
	for err := range errs {
		switch {
		case err == nil:
			succeeded++
		case errors.Is(err, repository.ErrDuplicateKey):
			conflicted++
		default:
			t.Fatalf("unexpected create error: %v", err)
		}
	}
	if succeeded != 1 || conflicted != 1 {
		t.Fatalf("successes=%d conflicts=%d, want 1/1", succeeded, conflicted)
	}
}

func TestPortRepositoryRejectsDuplicateClientLocal(t *testing.T) {
	db := openPortServiceDB(t)
	repo := repository.NewPortRepository(db)
	if err := repo.Create(&models.Port{ClientId: 7, Server: ":18081", Local: ":8081"}); err != nil {
		t.Fatalf("create first mapping: %v", err)
	}
	err := repo.Create(&models.Port{ClientId: 7, Server: ":18082", Local: ":8081"})
	if !errors.Is(err, repository.ErrDuplicateKey) {
		t.Fatalf("duplicate client/local error = %v, want ErrDuplicateKey", err)
	}
}

func TestInvalidPortErrorsUseSentinel(t *testing.T) {
	if _, err := normalizeServerPort("70000"); !errors.Is(err, ErrInvalidPort) {
		t.Fatalf("server port error = %v, want ErrInvalidPort", err)
	}
	if err := validateLocalAddr("not-an-address"); !errors.Is(err, ErrInvalidPort) {
		t.Fatalf("local address error = %v, want ErrInvalidPort", err)
	}
}

func TestPortServiceClosesListenerWhenUniqueInsertLosesRace(t *testing.T) {
	db := openPortServiceDB(t)
	client := models.Client{Name: "online", Symbol: "online-symbol", Status: models.StatusActive}
	if err := db.Create(&client).Error; err != nil {
		t.Fatalf("create client: %v", err)
	}
	runtime := newFakePortRuntime()
	runtime.online[client.Symbol] = true
	runtime.newListenHook = func() {
		err := repository.NewPortRepository(db).Create(&models.Port{
			ClientId: 999,
			Server:   ":18083",
			Local:    ":9999",
		})
		if err != nil {
			t.Fatalf("insert competing mapping: %v", err)
		}
		runtime.newListenHook = nil
	}

	_, err := NewPortService(db, runtime).Manage(PortCommand{
		Symbol: client.Symbol, ServerPort: "18083", LocalAddr: ":8083", Action: "add",
	})
	if !errors.Is(err, ErrPortConflict) {
		t.Fatalf("create error = %v, want ErrPortConflict", err)
	}
	if runtime.listeners[client.Symbol+":18083"] {
		t.Fatal("listener from losing insert was not closed")
	}
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
