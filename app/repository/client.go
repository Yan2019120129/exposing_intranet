package repository

import (
	"my-base/app/models"

	"gorm.io/gorm"
)

// ClientRepository owns persistence operations for clients and their port
// mappings. It deliberately exposes no HTTP or GoAdmin concerns.
type ClientRepository struct {
	DB *gorm.DB
}

type AdminUser struct {
	ID       int    `gorm:"column:id;primaryKey"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
	Name     string `gorm:"column:name"`
}

func (AdminUser) TableName() string { return "goadmin_users" }

type ClientNameOption struct {
	ID   int
	Name string
}

func NewClientRepository(db *gorm.DB) *ClientRepository {
	return &ClientRepository{DB: db}
}

func (r *ClientRepository) FindBySymbol(symbol string) (models.ClientAndPort, error) {
	var item models.ClientAndPort
	if r.DB == nil {
		return item, gorm.ErrInvalidDB
	}
	err := r.DB.Where("symbol = ?", symbol).Preload("PortList").First(&item).Error
	return item, err
}

func (r *ClientRepository) FindByIDs(ids []int) ([]models.Client, error) {
	items := make([]models.Client, 0)
	if r.DB == nil {
		return items, gorm.ErrInvalidDB
	}
	err := r.DB.Where("id IN ?", ids).Find(&items).Error
	return items, err
}

func (r *ClientRepository) FindAdminUser(username string) (AdminUser, error) {
	var user AdminUser
	if r.DB == nil {
		return user, gorm.ErrInvalidDB
	}
	err := r.DB.Where("username = ?", username).First(&user).Error
	return user, err
}

func (r *ClientRepository) Create(client *models.Client) error {
	if r.DB == nil {
		return gorm.ErrInvalidDB
	}
	return r.DB.Create(client).Error
}

func (r *ClientRepository) UpdateBySymbol(symbol string, client *models.Client) error {
	if r.DB == nil {
		return gorm.ErrInvalidDB
	}
	return r.DB.Where("symbol = ?", symbol).Updates(client).Error
}

func (r *ClientRepository) UpdateStatusBySymbol(symbol string, status int) error {
	return r.UpdateBySymbol(symbol, &models.Client{Status: status})
}

func (r *ClientRepository) NameOptions() ([]ClientNameOption, error) {
	items := make([]ClientNameOption, 0)
	if r.DB == nil {
		return items, gorm.ErrInvalidDB
	}
	err := r.DB.Model(&models.Client{}).Select("id", "name").Order("name asc").Find(&items).Error
	return items, err
}

func (r *ClientRepository) DeleteByIDs(ids []int) error {
	if r.DB == nil {
		return gorm.ErrInvalidDB
	}
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("client_id IN ?", ids).Delete(&models.Port{}).Error; err != nil {
			return err
		}
		return tx.Where("id IN ?", ids).Delete(&models.Client{}).Error
	})
}
