package service

import (
	"errors"
	"strings"

	"my-base/app/models"
	"my-base/app/repository"
	orm "my-base/code/gorm"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrClientNotFound      = errors.New("client not found")
	ErrClientNotRegistered = errors.New("client not registered")
	ErrClientDisabled      = errors.New("client disabled")
	ErrInvalidPort         = errors.New("invalid port")
	ErrInvalidAction       = errors.New("invalid port action")
	ErrPortConflict        = errors.New("port conflict")
	ErrPortNotFound        = errors.New("port not found")
)

type ClientRuntime interface {
	DelClientF(symbols ...string) error
}

type ClientService struct {
	DB         *gorm.DB
	Repository *repository.ClientRepository
	Runtime    ClientRuntime
}

type RegisterClientInput struct {
	Username string
	Password string
	Hostname string
}

func NewClientService(db *gorm.DB, runtimes ...ClientRuntime) *ClientService {
	var runtime ClientRuntime
	if len(runtimes) > 0 {
		runtime = runtimes[0]
	}
	return &ClientService{
		DB:         dbOrDefault(db),
		Repository: repository.NewClientRepository(db),
		Runtime:    runtime,
	}
}

func (s *ClientService) Register(input RegisterClientInput) (string, error) {
	user, err := s.Repository.FindAdminUser(strings.TrimSpace(input.Username))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", ErrInvalidCredentials
	}

	name := strings.TrimSpace(input.Hostname)
	if name == "" {
		name = strings.TrimSpace(input.Username)
	}
	client := &models.Client{
		Name:   name,
		Symbol: uuid.NewString(),
		Status: models.StatusOn,
	}
	if err := s.Repository.Create(client); err != nil {
		return "", err
	}
	return client.Symbol, nil
}

func (s *ClientService) FindBySymbol(symbol string) (models.ClientAndPort, error) {
	return s.Repository.FindBySymbol(symbol)
}

func (s *ClientService) UpdateStatus(symbol string, status int) error {
	return s.Repository.UpdateStatusBySymbol(symbol, status)
}

func (s *ClientService) UpdateNameAndStatus(symbol, name string, status int) error {
	return s.Repository.UpdateBySymbol(symbol, &models.Client{Name: name, Status: status})
}

func (s *ClientService) DeleteByIDs(ids []int) error {
	items, err := s.Repository.FindByIDs(ids)
	if err != nil {
		return err
	}
	if s.Runtime != nil {
		symbols := make([]string, 0, len(items))
		for _, item := range items {
			symbols = append(symbols, item.Symbol)
		}
		if err := s.Runtime.DelClientF(symbols...); err != nil {
			return err
		}
	}
	return s.Repository.DeleteByIDs(ids)
}

func (s *ClientService) NameOptions() ([]repository.ClientNameOption, error) {
	return s.Repository.NameOptions()
}

func dbOrDefault(db *gorm.DB) *gorm.DB {
	if db != nil {
		return db
	}
	return orm.DB
}
