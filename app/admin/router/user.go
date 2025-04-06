package router

import (
	"github.com/gin-gonic/gin"
	"my-base/app/admin/apis"
)

func init() {
	routerCheckRole = append(routerCheckRole, User)
}

// User 用户接口
func User(v1 *gin.RouterGroup) {
	api := apis.User{}
	r := v1.Group("/user")
	{
		r.GET("", api.GetPage)
		r.GET("/:id", api.Get)
	}
}
