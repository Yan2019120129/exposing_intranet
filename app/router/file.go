package router

import (
	"my-base/app/apis"

	"github.com/gin-gonic/gin"
)

func init() {
	routerCheckRole = append(routerCheckRole, fileRouter)
}

func fileRouter(v1 *gin.RouterGroup) {
	api := apis.File{}
	group := v1.Group("/admin/files")
	{
		group.GET("/download/:id", api.Download)
	}
}
