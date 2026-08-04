package file

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"my-base/app/models"
	"my-base/configs"

	"gorm.io/gorm"
)

// TusCleanupService 清理过期的 tus 暂存文件和无数据库引用的正式文件。
type TusCleanupService struct {
	stagingRoot string
	finalRoot   string
	retention   time.Duration
	now         func() time.Time
}

// NewTusCleanupService 创建 tus 残留文件清理服务。
func NewTusCleanupService(tusConfig *configs.TusUploadConfig) (*TusCleanupService, error) {
	if tusConfig == nil || strings.TrimSpace(tusConfig.StagingPath) == "" {
		return nil, errors.New("tus staging path is required")
	}
	adminConfig := configs.GetAdmin()
	if adminConfig == nil || strings.TrimSpace(adminConfig.Store.Path) == "" {
		return nil, errors.New("admin file store path is required")
	}
	stagingRoot, err := filepath.Abs(tusConfig.StagingPath)
	if err != nil {
		return nil, fmt.Errorf("resolve tus staging path: %w", err)
	}
	finalRoot, err := filepath.Abs(filepath.Join(adminConfig.Store.Path, "tus"))
	if err != nil {
		return nil, fmt.Errorf("resolve tus final path: %w", err)
	}
	return &TusCleanupService{
		stagingRoot: stagingRoot,
		finalRoot:   finalRoot,
		retention:   tusConfig.RetentionDuration(),
		now:         time.Now,
	}, nil
}

// Cleanup 清理超过保留时长且未被使用的 tus 文件。
func (s *TusCleanupService) Cleanup(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return errors.New("database is required")
	}
	cutoff := s.now().Add(-s.retention)
	return errors.Join(s.cleanupStaging(cutoff), s.cleanupFinal(ctx, db, cutoff))
}

// cleanupStaging 清理无活动锁且已经过期的暂存文件集合。
func (s *TusCleanupService) cleanupStaging(cutoff time.Time) error {
	entries, err := os.ReadDir(s.stagingRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var cleanupErrors []error
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".lock") || strings.HasSuffix(entry.Name(), ".stop") {
			baseName := strings.TrimSuffix(strings.TrimSuffix(entry.Name(), ".lock"), ".stop")
			if fileExists(filepath.Join(s.stagingRoot, baseName)) || fileExists(filepath.Join(s.stagingRoot, baseName+".info")) {
				continue
			}
			info, infoErr := entry.Info()
			if infoErr != nil {
				cleanupErrors = append(cleanupErrors, infoErr)
				continue
			}
			if info.ModTime().Before(cutoff) {
				if removeErr := os.Remove(filepath.Join(s.stagingRoot, entry.Name())); removeErr != nil && !os.IsNotExist(removeErr) {
					cleanupErrors = append(cleanupErrors, removeErr)
				}
			}
			continue
		}
		name := strings.TrimSuffix(entry.Name(), ".info")
		if entry.Name() != name+".info" && fileExists(filepath.Join(s.stagingRoot, name+".info")) {
			continue
		}
		if fileExists(filepath.Join(s.stagingRoot, name+".lock")) {
			continue
		}
		paths := []string{
			filepath.Join(s.stagingRoot, name),
			filepath.Join(s.stagingRoot, name+".info"),
			filepath.Join(s.stagingRoot, name+".stop"),
		}
		latest, exists, statErr := latestModification(paths)
		if statErr != nil {
			cleanupErrors = append(cleanupErrors, statErr)
			continue
		}
		if !exists || !latest.Before(cutoff) {
			continue
		}
		for _, candidate := range paths {
			if removeErr := os.Remove(candidate); removeErr != nil && !os.IsNotExist(removeErr) {
				cleanupErrors = append(cleanupErrors, removeErr)
			}
		}
	}
	return errors.Join(cleanupErrors...)
}

// cleanupFinal 清理正式目录中超过保留期且没有数据库记录的文件。
func (s *TusCleanupService) cleanupFinal(ctx context.Context, db *gorm.DB, cutoff time.Time) error {
	referencedPaths := make([]string, 0)
	if err := db.WithContext(ctx).Unscoped().Model(&models.SystemFile{}).
		Where("storage_path LIKE ?", "tus/%").
		Pluck("storage_path", &referencedPaths).Error; err != nil {
		return err
	}
	referenced := make(map[string]struct{}, len(referencedPaths))
	for _, value := range referencedPaths {
		referenced[filepath.Clean(filepath.FromSlash(value))] = struct{}{}
	}

	return filepath.WalkDir(s.finalRoot, func(filePath string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			if os.IsNotExist(walkErr) {
				return nil
			}
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.ModTime().Before(cutoff) {
			return nil
		}
		relative, err := filepath.Rel(filepath.Dir(s.finalRoot), filePath)
		if err != nil {
			return err
		}
		if _, ok := referenced[filepath.Clean(relative)]; ok {
			return nil
		}
		return os.Remove(filePath)
	})
}

// latestModification 返回一组可能存在文件中的最新修改时间。
func latestModification(paths []string) (time.Time, bool, error) {
	var latest time.Time
	exists := false
	for _, candidate := range paths {
		info, err := os.Stat(candidate)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return time.Time{}, false, err
		}
		exists = true
		if info.ModTime().After(latest) {
			latest = info.ModTime()
		}
	}
	return latest, exists, nil
}

// fileExists 判断普通路径是否存在。
func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
