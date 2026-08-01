package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"my-base/app/models"
	"my-base/code/service"

	"gorm.io/gorm"
)

const (
	fileShareOneDayHours     = 24
	fileShareSevenDaysHours  = 24 * 7
	fileShareThirtyDaysHours = 24 * 30
)

var (
	// ErrInvalidShareDuration 表示分享有效期不在允许范围内。
	ErrInvalidShareDuration = errors.New("invalid share duration")
	// ErrInvalidFileID 表示文件标识无效。
	ErrInvalidFileID = errors.New("invalid file id")
	// ErrInvalidShareID 表示分享标识无效。
	ErrInvalidShareID = errors.New("invalid share id")
	// ErrInvalidShareExpiresAt 表示分享链接到期时间无效。
	ErrInvalidShareExpiresAt = errors.New("invalid share expires at")
	// ErrShareRevoked 表示分享链接已被撤销，不能再修改。
	ErrShareRevoked = errors.New("share revoked")
)

// FileShare 提供系统文件分享记录的管理能力。
type FileShare struct {
	service.Service
}

// IsValidDuration 判断分享有效期是否为允许的固定选项。
func (e *FileShare) IsValidDuration(hours int) bool {
	return hours == fileShareOneDayHours || hours == fileShareSevenDaysHours || hours == fileShareThirtyDaysHours
}

// GetManageableFile 获取当前后台用户可管理的正常文件。
func (e *FileShare) GetManageableFile(fileID int, isSuperAdmin bool, item *models.SystemFile) error {
	if fileID <= 0 {
		return ErrInvalidFileID
	}
	if item == nil {
		return errors.New("file is required")
	}

	orm := e.Orm.Where("id = ? AND status = ?", fileID, 1)
	if !isSuperAdmin {
		orm = orm.Where("is_public = ?", true)
	}
	return orm.First(item).Error
}

// Create 创建指定有效期的文件分享记录。
func (e *FileShare) Create(fileID, durationHours int, createdBy int64, isSuperAdmin bool, item *models.SystemFile) (*models.SystemFileShare, error) {
	if item == nil {
		return nil, errors.New("file is required")
	}
	if err := e.GetManageableFile(fileID, isSuperAdmin, item); err != nil {
		return nil, err
	}
	if !e.IsValidDuration(durationHours) {
		return nil, ErrInvalidShareDuration
	}

	token, err := e.generateToken()
	if err != nil {
		return nil, err
	}
	share := &models.SystemFileShare{
		SystemFileID: fileID,
		Token:        token,
		ExpiresAt:    time.Now().Add(time.Duration(durationHours) * time.Hour),
		CreatedBy:    createdBy,
	}
	if err := e.Orm.Create(share).Error; err != nil {
		return nil, err
	}
	return share, nil
}

// List 返回指定文件的全部分享记录，按创建时间倒序排列。
func (e *FileShare) List(fileID int, shares *[]models.SystemFileShare) error {
	if fileID <= 0 {
		return ErrInvalidFileID
	}
	return e.Orm.Where("system_file_id = ?", fileID).Order("id desc").Find(shares).Error
}

// Revoke 撤销指定文件的分享记录。
func (e *FileShare) Revoke(fileID, shareID int) error {
	if fileID <= 0 || shareID <= 0 {
		return ErrInvalidShareID
	}
	now := time.Now()
	result := e.Orm.Model(&models.SystemFileShare{}).
		Where("id = ? AND system_file_id = ? AND revoked_at IS NULL", shareID, fileID).
		Update("revoked_at", now)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// UpdateExpiresAt 更新可管理且未撤销分享链接的到期时间。
func (e *FileShare) UpdateExpiresAt(shareID int, expiresAt time.Time, isSuperAdmin bool) error {
	if shareID <= 0 {
		return ErrInvalidShareID
	}
	if !expiresAt.After(time.Now()) {
		return ErrInvalidShareExpiresAt
	}

	share := models.SystemFileShare{}
	if err := e.Orm.Where("id = ?", shareID).First(&share).Error; err != nil {
		return err
	}
	if share.RevokedAt != nil {
		return ErrShareRevoked
	}

	item := models.SystemFile{}
	if err := e.GetManageableFile(share.SystemFileID, isSuperAdmin, &item); err != nil {
		return err
	}
	result := e.Orm.Model(&models.SystemFileShare{}).
		Where("id = ? AND revoked_at IS NULL", shareID).
		Update("expires_at", expiresAt)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrShareRevoked
	}
	return nil
}

// RevokeByID 撤销当前后台用户可管理的分享链接。
func (e *FileShare) RevokeByID(shareID int, isSuperAdmin bool) error {
	if shareID <= 0 {
		return ErrInvalidShareID
	}

	share := models.SystemFileShare{}
	if err := e.Orm.Where("id = ?", shareID).First(&share).Error; err != nil {
		return err
	}
	item := models.SystemFile{}
	if err := e.GetManageableFile(share.SystemFileID, isSuperAdmin, &item); err != nil {
		return err
	}
	return e.Revoke(share.SystemFileID, shareID)
}

// GetDownloadFile 根据有效分享令牌获取可下载文件。
func (e *FileShare) GetDownloadFile(token string, item *models.SystemFile) error {
	if token == "" {
		return gorm.ErrRecordNotFound
	}

	share := models.SystemFileShare{}
	if err := e.Orm.Where("token = ? AND revoked_at IS NULL AND expires_at > ?", token, time.Now()).First(&share).Error; err != nil {
		return err
	}
	return e.Orm.Where("id = ? AND status = ?", share.SystemFileID, 1).First(item).Error
}

// generateToken 生成不可预测的 URL 安全分享令牌。
func (e *FileShare) generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
