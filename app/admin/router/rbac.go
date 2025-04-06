package router

import (
	"github.com/gin-gonic/gin"
	"my-base/app/admin/apis"
)

func init() {
	routerCheckRole = append(routerCheckRole, registersRbac)
}

// Rbac 权限控制
type Rbac struct {
}

func registersRbac(v1 *gin.RouterGroup) {
	api := apis.Rbac{}
	// Login 登录接口
	v1.POST("/login", api.Login)

	// Register 登录接口
	v1.POST("/register", api.Register)
}
