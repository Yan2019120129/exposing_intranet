package router

import (
	"my-base/app/apis"

	"github.com/gin-gonic/gin"
)

func init() {
	routerNoCheckRole = append(routerNoCheckRole, systemFileShareRouter)
}

// systemFileShareRouter 注册系统文件分享路由。
func systemFileShareRouter(v1 *gin.RouterGroup) {
	api := apis.SystemFileShare{}
	group := v1.Group("/admin/files", SystemFileAccessMiddleware)
	{
		group.GET("/:id/shares", api.ShareList)
		group.POST("/:id/shares", api.CreateShare)
		group.DELETE("/:id/shares/:shareId", api.RevokeShare)
	}
	v1.GET("/shares/:token/download", api.ShareDownload)
}
