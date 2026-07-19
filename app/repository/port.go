package repository

import (
	"my-base/app/models"

	"gorm.io/gorm"
)

type PortRepository struct {
	DB *gorm.DB
}

type PortWithClient struct {
	ID         int
	Server     string
	Local      string
	Comment    string
	ClientID   int
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

func (r *PortRepository) CountClientConflict(clientID int, server, local string) (int64, error) {
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
	return r.DB.Create(port).Error
}

func (r *PortRepository) DeleteByID(id int) error {
	if r.DB == nil {
		return gorm.ErrInvalidDB
	}
	return r.DB.Delete(&models.Port{}, id).Error
}

func (r *PortRepository) DeleteByIDs(ids []string) error {
	if r.DB == nil {
		return gorm.ErrInvalidDB
	}
	return r.DB.Where("id IN ?", ids).Delete(&models.Port{}).Error
}

func (r *PortRepository) FindByClientID(clientID int) ([]models.Port, error) {
	items := make([]models.Port, 0)
	if r.DB == nil {
		return items, gorm.ErrInvalidDB
	}
	err := r.DB.Where("client_id = ?", clientID).Order("id asc").Find(&items).Error
	return items, err
}

func (r *PortRepository) FindByClientAndServer(clientID int, server string) (models.Port, error) {
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
