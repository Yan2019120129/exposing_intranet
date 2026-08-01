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

// List 根据查询条件获取测试记录，名称为空时返回全部记录。
func (e *Test) List(query *dto.TestListQuery, list *[]models.Test) error {
	orm := e.Orm.Order("id asc")
	if query != nil {
		query.Normalize()
		if query.Name != "" {
			orm = orm.Where("name = ?", query.Name)
		}
	}
	return orm.Find(list).Error
}

// Create 创建一条名称非空的测试记录。
func (e *Test) Create(payload *dto.TestPayload, item *models.Test) error {
	name := strings.TrimSpace(payload.Name)
	if name == "" {
		return errors.New("name is required")
	}

	item.Name = name
	return e.Orm.Create(item).Error
}

// Get 获取指定标识的未删除测试记录。
func (e *Test) Get(id uint, item *models.Test) error {
	if id <= 0 {
		return errors.New("invalid id")
	}
	return e.Orm.First(item, id).Error
}

// Update 更新指定未删除测试记录的名称。
func (e *Test) Update(id uint, payload *dto.TestPayload, item *models.Test) error {
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

// Delete 软删除指定测试记录。
func (e *Test) Delete(id uint) error {
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
