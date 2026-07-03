package code

import (
	"errors"
	"log"
	"runtime"
	"strconv"

	adminDB "github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	TrafficKey                   = "X-Request-Id"
	LoggerKey                    = "_go-admin-logger-request"
	ContextDBKey                 = "db"
	DefaultGoAdminConnectionName = "default"
)

func CompareHashAndPassword(e string, p string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(e), []byte(p))
	if err != nil {
		return false, err
	}
	return true, nil
}

// Assert 条件断言
// 当断言条件为 假 时触发 panic
// 对于当前请求不会再执行接下来的代码，并且返回指定格式的错误信息和错误码
func Assert(condition bool, msg string, code ...int) {
	if !condition {
		statusCode := 200
		if len(code) > 0 {
			statusCode = code[0]
		}
		panic("CustomError#" + strconv.Itoa(statusCode) + "#" + msg)
	}
}

// HasError 错误断言
// 当 error 不为 nil 时触发 panic
// 对于当前请求不会再执行接下来的代码，并且返回指定格式的错误信息和错误码
// 若 msg 为空，则默认为 error 中的内容
func HasError(err error, msg string, code ...int) {
	if err != nil {
		statusCode := 200
		if len(code) > 0 {
			statusCode = code[0]
		}
		if msg == "" {
			msg = err.Error()
		}
		_, file, line, _ := runtime.Caller(1)
		log.Printf("%s:%v error: %#v", file, line, err)
		panic("CustomError#" + strconv.Itoa(statusCode) + "#" + msg)
	}
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

// GetOrm 获取orm连接
func GetOrm(c *gin.Context) (*gorm.DB, error) {
	con, err := GetGoAdminConnection(c)
	if err != nil {
		return nil, err
	}
	return con.GetGorm(DefaultGoAdminConnectionName)
}

// GetGoAdminConnection 获取 GoAdmin 数据库连接。
func GetGoAdminConnection(c *gin.Context) (adminDB.Connection, error) {
	idb, exist := c.Get(ContextDBKey)
	if !exist {
		return nil, errors.New("db connect not exist")
	}
	conn, ok := idb.(adminDB.Connection)
	if !ok {
		return nil, errors.New("go-admin db connection not exist")
	}
	return conn, nil
}
