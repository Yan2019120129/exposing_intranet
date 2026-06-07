package tasks

import (
	"context"
	"time"

	"my-base/app/models"

	"github.com/GoAdminGroup/go-admin/modules/logger"
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

			logger.Infof("test table monitor ok: count=%d", count)
			return nil
		},
	}
}
