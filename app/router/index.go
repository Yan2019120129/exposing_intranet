package router

import (
	"github.com/gin-gonic/gin"
)

var (
	routerNoCheckRole = make([]func(*gin.RouterGroup), 0)
	routerCheckRole   = make([]func(v1 *gin.RouterGroup), 0)
)

// InitRouter 初始化路由
func InitRouter(r *gin.Engine, authMiddlewares ...gin.HandlerFunc) {
	// 无需认证的路由
	examplesNoCheckRoleRouter(r)
	// 需要认证的路由
	examplesCheckRoleRouter(r, authMiddlewares...)
}

// 无需认证的路由示例
func examplesNoCheckRoleRouter(r *gin.Engine) {
	// 可根据业务需求来设置接口版本
	v1 := r.Group("")
	for _, f := range routerNoCheckRole {
		f(v1)
	}
}

// 需要认证的路由示例
func examplesCheckRoleRouter(r *gin.Engine, authMiddlewares ...gin.HandlerFunc) {
	// 可根据业务需求来设置接口版本
	v1 := r.Group("", authMiddlewares...)
	for _, f := range routerCheckRole {
		f(v1)
	}
}
