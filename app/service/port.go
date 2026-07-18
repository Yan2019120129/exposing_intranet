package service

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"my-base/app/models"
	"my-base/app/repository"
	"my-base/code/penetrate"

	"gorm.io/gorm"
)

type PortRuntime interface {
	IsExist(clientSymbol string) bool
	IsListenExist(clientSymbol, serverPort string) bool
	NewListen(clientSymbol, localPort, serverPort string) error
	CloseListen(clientSymbol, localPort, serverPort string) error
}

type PortService struct {
	Repository       *repository.PortRepository
	ClientRepository *repository.ClientRepository
	Runtime          PortRuntime
}

type PortCommand struct {
	Symbol     string
	Action     string
	ServerPort string
	LocalAddr  string
	Comment    string
}

type PortMapping struct {
	ServerPort string
	LocalAddr  string
	Comment    string
	Status     string
}

type PortResult struct {
	Mappings []PortMapping
}

func NewPortService(db *gorm.DB, runtimes ...PortRuntime) *PortService {
	runtime := PortRuntime(penetrate.GetServer())
	if len(runtimes) > 0 && runtimes[0] != nil {
		runtime = runtimes[0]
	}
	return &PortService{
		Repository:       repository.NewPortRepository(db),
		ClientRepository: repository.NewClientRepository(db),
		Runtime:          runtime,
	}
}

func (s *PortService) Manage(command PortCommand) (PortResult, error) {
	command.Symbol = strings.TrimSpace(command.Symbol)
	client, err := s.ClientRepository.FindBySymbol(command.Symbol)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return PortResult{}, ErrClientNotRegistered
		}
		return PortResult{}, err
	}

	switch command.Action {
	case "add":
		mapping, err := s.create(client.Id, command.ServerPort, command.LocalAddr, command.Comment, command.Symbol)
		if err != nil {
			return PortResult{}, err
		}
		return PortResult{Mappings: []PortMapping{mapping}}, nil
	case "del":
		mapping, err := s.delete(client.Id, command.ServerPort, command.Symbol)
		if err != nil {
			return PortResult{}, err
		}
		return PortResult{Mappings: []PortMapping{mapping}}, nil
	case "list":
		return PortResult{Mappings: s.list(client)}, nil
	default:
		return PortResult{}, ErrInvalidAction
	}
}

// CreateFromAdmin creates a mapping from a GoAdmin form while keeping all
// validation, bind and persistence behavior in the service layer.
func (s *PortService) CreateFromAdmin(clientID int, serverPort, localAddr, comment string) error {
	client, err := s.ClientRepository.FindByIDs([]int{clientID})
	if err != nil {
		return err
	}
	if len(client) == 0 {
		return ErrClientNotFound
	}
	_, err = s.create(clientID, serverPort, localAddr, comment, client[0].Symbol)
	return err
}

func (s *PortService) DeleteByIDs(ids []string) error {
	items, err := s.Repository.FindWithClientByIDs(ids)
	if err != nil {
		return err
	}
	for _, item := range items {
		if s.Runtime.IsExist(item.Symbol) {
			if err := s.Runtime.CloseListen(item.Symbol, item.Local, item.Server); err != nil {
				return err
			}
		}
	}
	return s.Repository.DeleteByIDs(ids)
}

func (s *PortService) create(clientID int, serverPort, localAddr, comment, symbol string) (PortMapping, error) {
	serverPort, err := normalizeServerPort(serverPort)
	if err != nil {
		return PortMapping{}, err
	}
	if err := validateLocalAddr(localAddr); err != nil {
		return PortMapping{}, err
	}

	count, err := s.Repository.CountServer(serverPort)
	if err != nil {
		return PortMapping{}, err
	}
	if count > 0 {
		return PortMapping{}, ErrPortConflict
	}
	count, err = s.Repository.CountClientConflict(clientID, serverPort, localAddr)
	if err != nil {
		return PortMapping{}, err
	}
	if count > 0 {
		return PortMapping{}, ErrPortConflict
	}

	bound := false
	if s.Runtime.IsExist(symbol) {
		if err := s.Runtime.NewListen(symbol, localAddr, serverPort); err != nil {
			return PortMapping{}, err
		}
		bound = true
	}

	port := &models.Port{ClientId: clientID, Server: serverPort, Local: localAddr, Comment: comment}
	if err := s.Repository.Create(port); err != nil {
		if bound {
			_ = s.Runtime.CloseListen(symbol, localAddr, serverPort)
		}
		return PortMapping{}, err
	}
	return PortMapping{ServerPort: serverPort, LocalAddr: localAddr, Comment: comment, Status: "active"}, nil
}

func (s *PortService) delete(clientID int, serverPort, symbol string) (PortMapping, error) {
	serverPort, err := normalizeServerPort(serverPort)
	if err != nil {
		return PortMapping{}, err
	}
	port, err := s.Repository.FindByClientAndServer(clientID, serverPort)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return PortMapping{}, ErrPortNotFound
		}
		return PortMapping{}, err
	}
	if s.Runtime.IsExist(symbol) {
		if err := s.Runtime.CloseListen(symbol, port.Local, serverPort); err != nil {
			return PortMapping{}, err
		}
	}
	if err := s.Repository.DeleteByID(port.Id); err != nil {
		return PortMapping{}, err
	}
	return PortMapping{ServerPort: serverPort, LocalAddr: port.Local, Comment: port.Comment, Status: "offline"}, nil
}

func (s *PortService) list(client models.ClientAndPort) []PortMapping {
	mappings := make([]PortMapping, 0, len(client.PortList))
	for _, port := range client.PortList {
		status := "offline"
		if s.Runtime.IsListenExist(client.Symbol, port.Server) {
			status = "active"
		}
		mappings = append(mappings, PortMapping{
			ServerPort: port.Server,
			LocalAddr:  port.Local,
			Comment:    port.Comment,
			Status:     status,
		})
	}
	return mappings
}

func normalizeServerPort(value string) (string, error) {
	value = strings.TrimPrefix(strings.TrimSpace(value), ":")
	port, err := strconv.Atoi(value)
	if err != nil || port < 1 || port > 65535 {
		return "", fmt.Errorf("服务端口必须是 1-65535")
	}
	return ":" + strconv.Itoa(port), nil
}

func validateLocalAddr(value string) error {
	value = strings.TrimSpace(value)
	_, port, err := net.SplitHostPort(value)
	if err != nil {
		return fmt.Errorf("本地地址必须是 host:port 格式")
	}
	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		return fmt.Errorf("本地端口必须是 1-65535")
	}
	return nil
}
