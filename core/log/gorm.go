package log

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	utils2 "my-base/utils"
	"time"
)

var LeverGorm = map[string]logger.LogLevel{
	"debug":  logger.Info,
	"info":   logger.Info,
	"warn":   logger.Warn,
	"error":  logger.Error,
	"dpanic": logger.Silent,
	"panic":  logger.Silent,
	"fatal":  logger.Silent,
}

var Instance = &zapLogger{
	Config: logger.Config{
		SlowThreshold:             500 * time.Millisecond,
		Colorful:                  true,
		IgnoreRecordNotFoundError: false,
		//ParameterizedQueries:      false,
		LogLevel: logger.Warn,
	},
}

// ZapLogger zap日志配置
type zapLogger struct {
	logger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// LogMode 配置日志模式
func (l *zapLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 配置info日志
func (l *zapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		Logger.WithOptions(zap.WithCaller(false)).
			Info(SetLogFormat(
				utils.FileWithLineNum(),
				utils2.SetRed(msg),
				SetLogFormat(data),
			))
	}
}

// Warn 配置Warn日志
func (l *zapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		Logger.WithOptions(zap.WithCaller(false)).
			Warn(SetLogFormat(
				utils.FileWithLineNum(),
				utils2.SetRed(msg),
				SetLogFormat(data),
			))
	}
}

// Error 配置Error日志
func (l *zapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		Logger.WithOptions(zap.WithCaller(false)).
			Error(SetLogFormat(
				utils.FileWithLineNum(),
				utils2.SetRed(msg),
				SetLogFormat(data),
			))
	}
}

// Trace 配置Trace日志
func (l *zapLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			Logger.WithOptions(zap.WithCaller(false)).
				Error(SetLogFormat(
					utils.FileWithLineNum(),
					utils2.SetRed(err.Error()),
					"\n"+l.getNanoseconds(elapsed),
					GetRows(rows),
					sql,
				))
		} else {
			Logger.WithOptions(zap.WithCaller(false)).
				Error(SetLogFormat(
					utils.FileWithLineNum(),
					utils2.SetRed(err.Error()),
					"\n"+l.getNanoseconds(elapsed),
					GetRows(rows),
					sql,
				))
		}

	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			Logger.WithOptions(zap.WithCaller(false)).
				Warn(SetLogFormat(
					utils.FileWithLineNum(),
					utils2.SetYellow(slowLog),
					"\n"+l.getNanoseconds(elapsed),
					GetRows(rows),
					sql,
				))
		} else {
			Logger.WithOptions(zap.WithCaller(false)).
				Warn(SetLogFormat(
					utils.FileWithLineNum(),
					utils2.SetYellow(slowLog),
					"\n"+l.getNanoseconds(elapsed),
					GetRows(rows),
					sql,
				))
		}

	case l.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			Logger.WithOptions(zap.WithCaller(false)).
				Info(SetLogFormat(
					utils.FileWithLineNum(),
					"\n"+l.getNanoseconds(elapsed),
					GetRows(rows),
					sql,
				))
		} else {
			Logger.WithOptions(zap.WithCaller(false)).
				Info(SetLogFormat(
					utils.FileWithLineNum(),
					"\n"+l.getNanoseconds(elapsed),
					GetRows(rows),
					sql,
				))
		}
	}
}

// GetNanoseconds 获取拼接的执行时间
func (l *zapLogger) getNanoseconds(nanoseconds time.Duration) string {
	nanosecondsTemp := float64(nanoseconds.Nanoseconds()) / 1e6
	msg := ""
	switch {
	// 大于1秒，小于500毫秒秒为黄色
	case nanoseconds > l.SlowThreshold:
		msg = utils2.SetYellow(nanosecondsTemp)
	// 大于1秒为红色
	case nanoseconds > 2*l.SlowThreshold:
		msg = utils2.SetRed(nanosecondsTemp)
	// 默认绿色
	default:
		msg = utils2.SetGreen(nanosecondsTemp)
	}
	return fmt.Sprintf("spend:[%v]", msg)
}

// GetRows 获取拼接的行数
func GetRows(rows int64) string {
	return fmt.Sprintf("rows:[%v]", rows)
}

// SetLogFormat 设置日志格式
func SetLogFormat(msg ...interface{}) string {
	length := len(msg)
	if length == 0 {
		return ""
	}
	message := fmt.Sprintf("%v", msg[0])
	for i := 1; i < length; i++ {
		message += fmt.Sprintf(" %v", msg[i])
	}
	return message
}
