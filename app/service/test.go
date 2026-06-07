package service

import (
	"errors"
	"my-base/app/models"
	"my-base/app/service/dto"
	"my-base/code/service"
	"strings"

	"gorm.io/gorm"
)

type Test struct {
	service.Service
}

func (e *Test) List(list *[]models.Test) error {
	return e.Orm.Order("id asc").Find(list).Error
}

func (e *Test) Create(payload *dto.TestPayload, item *models.Test) error {
	name := strings.TrimSpace(payload.Name)
	if name == "" {
		return errors.New("name is required")
	}

	item.Name = name
	return e.Orm.Create(item).Error
}

func (e *Test) Get(id int, item *models.Test) error {
	if id <= 0 {
		return errors.New("invalid id")
	}
	return e.Orm.First(item, id).Error
}

func (e *Test) Update(id int, payload *dto.TestPayload, item *models.Test) error {
	if id <= 0 {
		return errors.New("invalid id")
	}

	name := strings.TrimSpace(payload.Name)
	if name == "" {
		return errors.New("name is required")
	}

	if err := e.Orm.First(item, id).Error; err != nil {
		return err
	}
	item.Name = name
	return e.Orm.Save(item).Error
}

func (e *Test) Delete(id int) error {
	if id <= 0 {
		return errors.New("invalid id")
	}

	result := e.Orm.Delete(&models.Test{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
