package router

import (
	"github.com/gin-gonic/gin"
	"my-base/app/admin/apis"
)

func init() {
	routerCheckRole = append(routerCheckRole, Test)
}

// Test 登录接口
func Test(v1 *gin.RouterGroup) {
	v1.GET("/test", apis.Test)
}
