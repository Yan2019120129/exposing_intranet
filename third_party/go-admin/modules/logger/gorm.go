package logger

import (
	stdctx "context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const defaultGormSlowThreshold = 200 * time.Millisecond

type GormLoggerConfig struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	LogLevel                  gormlogger.LogLevel
}

type gormLogger struct {
	config GormLoggerConfig
}

func GormLogger(configs ...GormLoggerConfig) gormlogger.Interface {
	cfg := GormLoggerConfig{
		SlowThreshold:             defaultGormSlowThreshold,
		IgnoreRecordNotFoundError: true,
		LogLevel:                  gormlogger.Info,
	}
	if len(configs) > 0 {
		cfg = configs[0]
		if cfg.SlowThreshold == 0 {
			cfg.SlowThreshold = defaultGormSlowThreshold
		}
	}
	return &gormLogger{config: cfg}
}

func (l *gormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.config.LogLevel = level
	return &newLogger
}

func (l *gormLogger) Info(ctx stdctx.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= gormlogger.Info {
		Infof(msg, data...)
	}
}

func (l *gormLogger) Warn(ctx stdctx.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= gormlogger.Warn {
		Warnf(msg, data...)
	}
}

func (l *gormLogger) Error(ctx stdctx.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= gormlogger.Error {
		Errorf(msg, data...)
	}
}

func (l *gormLogger) Trace(ctx stdctx.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.config.LogLevel <= gormlogger.Silent || !logger.sqlLogOpen {
		return
	}

	elapsed := time.Since(begin)
	elapsedMS := float64(elapsed.Nanoseconds()) / 1e6

	switch {
	case err != nil && l.config.LogLevel >= gormlogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.config.IgnoreRecordNotFoundError):
		sql, rows := fc()
		Errorf("[GORM] err=%v elapsed=%.3fms rows=%s sql=%s", err, elapsedMS, formatGormRows(rows), sql)
	case l.config.SlowThreshold != 0 && elapsed > l.config.SlowThreshold && l.config.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		Warnf("[GORM] slow sql >= %v elapsed=%.3fms rows=%s sql=%s", l.config.SlowThreshold, elapsedMS, formatGormRows(rows), sql)
	case l.config.LogLevel == gormlogger.Info:
		sql, rows := fc()
		Infof("[GORM] elapsed=%.3fms rows=%s sql=%s", elapsedMS, formatGormRows(rows), sql)
	}
}

func (l *gormLogger) ParamsFilter(ctx stdctx.Context, sql string, params ...interface{}) (string, []interface{}) {
	return sql, params
}

func formatGormRows(rows int64) string {
	if rows == -1 {
		return "-"
	}
	return fmt.Sprintf("%d", rows)
}
