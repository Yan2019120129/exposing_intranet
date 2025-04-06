package apis

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Rbac struct {
}

// Login 登录接口
func (r *Rbac) Login(c *gin.Context) {
	c.String(http.StatusOK, "login")
}

// Register 注册接口
func (r *Rbac) Register(c *gin.Context) {
	c.String(http.StatusOK, "register")
}
