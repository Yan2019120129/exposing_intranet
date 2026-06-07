package api

import (
	"car/code/logger"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	TrafficKey = "X-Request-Id"
	LoggerKey  = "_go-admin-logger-request"
)

// GetRequestLogger 获取上下文提供的日志
func GetRequestLogger(c *gin.Context) *logger.Helper {
	var log *logger.Helper
	l, ok := c.Get("logger")
	if ok {
		ok = false
		log, ok = l.(*logger.Helper)
		if ok {
			return log
		}
	}

	//如果没有在上下文中放入logger
	requestId := GenerateMsgIDFromContext(c)
	log = logger.NewHelper(logger.DefaultLogger).WithFields(map[string]interface{}{
		strings.ToLower(TrafficKey): requestId,
	})
	return log
}

// SetRequestLogger 设置logger中间件
func SetRequestLogger(c *gin.Context) {
	requestId := GenerateMsgIDFromContext(c)
	log := logger.NewHelper(logger.DefaultLogger).WithFields(map[string]interface{}{
		strings.ToLower(TrafficKey): requestId,
	})
	c.Set(LoggerKey, log)
}

// GenerateMsgIDFromContext 生成msgID
func GenerateMsgIDFromContext(c *gin.Context) string {
	requestId := c.GetHeader(TrafficKey)
	if requestId == "" {
		requestId = uuid.New().String()
		c.Header(TrafficKey, requestId)
	}
	return requestId
}
