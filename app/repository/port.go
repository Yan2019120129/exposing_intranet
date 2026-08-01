package repository

import (
	"errors"
	"fmt"

	"my-base/app/models"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var ErrDuplicateKey = errors.New("duplicate key")

type PortRepository struct {
	DB *gorm.DB
}

type PortWithClient struct {
	ID         uint
	Server     string
	Local      string
	Comment    string
	ClientID   uint
	ClientName string
	Symbol     string
}

func NewPortRepository(db *gorm.DB) *PortRepository {
	return &PortRepository{DB: db}
}

func (r *PortRepository) CountServer(server string) (int64, error) {
	if r.DB == nil {
		return 0, gorm.ErrInvalidDB
	}
	var count int64
	err := r.DB.Model(&models.Port{}).Where("server = ?", server).Count(&count).Error
	return count, err
}

func (r *PortRepository) CountClientConflict(clientID uint, server, local string) (int64, error) {
	if r.DB == nil {
		return 0, gorm.ErrInvalidDB
	}
	var count int64
	err := r.DB.Model(&models.Port{}).
		Where("client_id = ?", clientID).
		Where("server = ? OR local = ?", server, local).
		Count(&count).Error
	return count, err
}

func (r *PortRepository) Create(port *models.Port) error {
	if r.DB == nil {
		return gorm.ErrInvalidDB
	}
	err := r.DB.Create(port).Error
	if err == nil {
		return nil
	}
	var mysqlErr *mysqlDriver.MySQLError
	if errors.Is(err, gorm.ErrDuplicatedKey) || (errors.As(err, &mysqlErr) && mysqlErr.Number == 1062) {
		return fmt.Errorf("%w: %v", ErrDuplicateKey, err)
	}
	return err
}

// DeleteByID 物理删除单个端口映射，确保端口可以再次创建。
func (r *PortRepository) DeleteByID(id uint) error {
	if r.DB == nil {
		return gorm.ErrInvalidDB
	}
	return r.DB.Unscoped().Delete(&models.Port{}, id).Error
}

// DeleteByIDs 物理删除多个端口映射，确保端口可以再次创建。
func (r *PortRepository) DeleteByIDs(ids []string) error {
	if r.DB == nil {
		return gorm.ErrInvalidDB
	}
	return r.DB.Unscoped().Where("id IN ?", ids).Delete(&models.Port{}).Error
}

func (r *PortRepository) FindByClientID(clientID uint) ([]models.Port, error) {
	items := make([]models.Port, 0)
	if r.DB == nil {
		return items, gorm.ErrInvalidDB
	}
	err := r.DB.Where("client_id = ?", clientID).Order("id asc").Find(&items).Error
	return items, err
}

func (r *PortRepository) FindByClientAndServer(clientID uint, server string) (models.Port, error) {
	var item models.Port
	if r.DB == nil {
		return item, gorm.ErrInvalidDB
	}
	err := r.DB.Where("client_id = ? AND server = ?", clientID, server).First(&item).Error
	return item, err
}

func (r *PortRepository) FindWithClientByIDs(ids []string) ([]PortWithClient, error) {
	items := make([]PortWithClient, 0)
	if r.DB == nil {
		return items, gorm.ErrInvalidDB
	}
	err := r.DB.Table("port AS p").
		Select("p.id, p.server, p.local, p.comment, p.client_id, c.name AS client_name, c.symbol").
		Joins("LEFT JOIN client AS c ON c.id = p.client_id").
		Where("p.id IN ?", ids).
		Find(&items).Error
	return items, err
}
