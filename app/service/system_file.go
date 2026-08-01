package service

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"my-base/app/models"
	"my-base/app/service/dto"
	"my-base/code/service"
	"my-base/configs"
)

// SystemFile 提供系统文件查询和存储路径处理能力。
type SystemFile struct {
	service.Service
}

// GetDownloadFile 获取可下载的正常系统文件。
func (e *SystemFile) GetDownloadFile(input *dto.SystemFileDownloadInput, item *models.SystemFile) error {
	if input == nil || input.FileID <= 0 {
		return errors.New("invalid file id")
	}
	if item == nil {
		return errors.New("file is required")
	}
	return e.Orm.Where("id = ? AND status = ?", input.FileID, 1).First(item).Error
}

// ResolveStoragePath 将系统文件存储路径解析为受限的本地绝对路径。
func (e *SystemFile) ResolveStoragePath(storagePath string) (string, error) {
	storagePath = strings.TrimSpace(strings.ReplaceAll(storagePath, "\\", "/"))
	if storagePath == "" {
		return "", errors.New("empty storage path")
	}
	storagePath = strings.TrimPrefix(storagePath, "/")
	prefix := strings.Trim(configs.GetAdmin().Store.Prefix, "/")
	if prefix != "" && strings.HasPrefix(storagePath, prefix+"/") {
		storagePath = strings.TrimPrefix(storagePath, prefix+"/")
	}
	cleaned := filepath.Clean(filepath.FromSlash(storagePath))
	if cleaned == "." || strings.HasPrefix(cleaned, ".."+string(os.PathSeparator)) || filepath.IsAbs(cleaned) {
		return "", errors.New("invalid storage path")
	}
	return filepath.Join(configs.GetAdmin().Store.Path, cleaned), nil
}
