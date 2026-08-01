package service

import (
	"errors"
	"testing"

	"my-base/app/models"
	"my-base/app/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openClientServiceDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+uuid.NewString()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.Client{}, &models.Port{}, &repository.AdminUser{}); err != nil {
		t.Fatalf("migrate sqlite: %v", err)
	}
	return db
}

func TestClientServiceRegister(t *testing.T) {
	db := openClientServiceDB(t)
	hash, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	if err := db.Create(&repository.AdminUser{Username: "admin", Password: string(hash)}).Error; err != nil {
		t.Fatalf("create admin: %v", err)
	}

	clientService := NewClientService(db)
	symbol, err := clientService.Register(RegisterClientInput{
		Username: "admin",
		Password: "secret",
		Hostname: "workstation",
	})
	if err != nil {
		t.Fatalf("register client: %v", err)
	}
	if symbol == "" {
		t.Fatal("expected generated client symbol")
	}

	var client models.Client
	if err := db.Where("symbol = ?", symbol).First(&client).Error; err != nil {
		t.Fatalf("load registered client: %v", err)
	}
	if client.Name != "workstation" || client.Status != models.StatusOn {
		t.Fatalf("unexpected registered client: %+v", client)
	}

	_, err = clientService.Register(RegisterClientInput{Username: "admin", Password: "wrong"})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

type observingDeleteRuntime struct {
	db        *gorm.DB
	symbols   []string
	rowGone   bool
	portsGone bool
}

func (r *observingDeleteRuntime) DelClientF(symbols ...string) error {
	r.symbols = append(r.symbols, symbols...)
	var clients int64
	var ports int64
	if err := r.db.Model(&models.Client{}).Where("symbol IN ?", symbols).Count(&clients).Error; err != nil {
		return err
	}
	if err := r.db.Model(&models.Port{}).Count(&ports).Error; err != nil {
		return err
	}
	r.rowGone = clients == 0
	r.portsGone = ports == 0
	return nil
}

func TestClientDeleteCommitsDatabaseBeforeRuntimeCleanup(t *testing.T) {
	db := openClientServiceDB(t)
	client := models.Client{Name: "delete-me", Symbol: "delete-symbol", Status: models.StatusActive}
	if err := db.Create(&client).Error; err != nil {
		t.Fatalf("create client: %v", err)
	}
	if err := db.Create(&models.Port{ClientId: client.ID, Server: ":18090", Local: ":8090"}).Error; err != nil {
		t.Fatalf("create port: %v", err)
	}

	runtime := &observingDeleteRuntime{db: db}
	if err := NewClientService(db, runtime).DeleteByIDs([]int{int(client.ID)}); err != nil {
		t.Fatalf("delete client: %v", err)
	}
	if len(runtime.symbols) != 1 || runtime.symbols[0] != client.Symbol {
		t.Fatalf("runtime symbols = %v", runtime.symbols)
	}
	if !runtime.rowGone || !runtime.portsGone {
		t.Fatalf("runtime cleanup ran before commit: clientGone=%v portsGone=%v", runtime.rowGone, runtime.portsGone)
	}
}
