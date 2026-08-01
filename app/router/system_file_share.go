package router

import (
	"my-base/app/apis"

	"github.com/gin-gonic/gin"
)

func init() {
	routerNoCheckRole = append(routerNoCheckRole, fileShareRouter)
}

func fileShareRouter(v1 *gin.RouterGroup) {
	api := apis.File{}
	group := v1.Group("/admin/files", SystemFileAccessMiddleware)
	{
		group.GET("/:id/shares", api.ShareList)
		group.POST("/:id/shares", api.CreateShare)
		group.DELETE("/:id/shares/:shareId", api.RevokeShare)
	}
	v1.GET("/shares/:token/download", api.ShareDownload)
}
