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
