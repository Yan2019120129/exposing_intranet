package router

import (
	"my-base/app/apis"

	"github.com/gin-gonic/gin"
)

func init() {
	routerNoCheckRole = append(routerNoCheckRole, portRouter)
}

func portRouter(v1 *gin.RouterGroup) {
	apiGroup := v1.Group("/api")
	port := apis.Port{}
	apiGroup.POST("/client/port", port.Manage)
}
