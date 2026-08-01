package router

import (
	"my-base/app/apis"

	"github.com/gin-gonic/gin"
)

func init() {
	routerCheckRole = append(routerCheckRole, systemFileRouter)
}

// systemFileRouter 注册系统文件下载路由。
func systemFileRouter(v1 *gin.RouterGroup) {
	api := apis.SystemFile{}
	group := v1.Group("/admin/files")
	{
		group.GET("/download/:id", api.Download)
	}
}
