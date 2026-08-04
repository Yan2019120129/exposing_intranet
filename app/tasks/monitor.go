package tasks

import (
	"context"
	"time"

	"my-base/app/models"
	fileService "my-base/code/file"
	"my-base/configs"

	"github.com/GoAdminGroup/go-admin/modules/logger"
	"github.com/GoAdminGroup/go-admin/modules/redis"
	"gorm.io/gorm"
)

func heartbeatJob() Job {
	return Job{
		Name:     "heartbeat",
		Interval: time.Minute,
		Run: func(ctx context.Context) error {
			logger.Info("task heartbeat ok")
			return nil
		},
	}
}

func testTableMonitorJob(db *gorm.DB) Job {
	return Job{
		Name:     "test_table_monitor",
		Interval: 5 * time.Minute,
		Run: func(ctx context.Context) error {
			var count int64
			if err := db.WithContext(ctx).Model(&models.Test{}).Count(&count).Error; err != nil {
				return err
			}

			const warnThreshold = 1000
			if count > warnThreshold {
				logger.Infof("test table row count is high: count=%d threshold=%d", count, warnThreshold)
				return nil
			}

			_, err := redis.Do("SET", "some_key", count)
			if err != nil {
				return err
			}

			logger.Infof("test table monitor ok: count=%d", count)
			return nil
		},
	}
}

// tusUploadCleanupJob 定期清理过期的 tus 上传残留文件。
func tusUploadCleanupJob(db *gorm.DB) Job {
	return Job{
		Name:     "tus_upload_cleanup",
		Interval: time.Hour,
		Run: func(ctx context.Context) error {
			if !configs.GetTusUpload().Enabled {
				return nil
			}
			service, err := fileService.NewTusCleanupService(configs.GetTusUpload())
			if err != nil {
				return err
			}
			return service.Cleanup(ctx, db)
		},
	}
}
