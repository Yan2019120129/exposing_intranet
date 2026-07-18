package router

import (
	"my-base/app/apis"

	"github.com/gin-gonic/gin"
)

func init() {
	routerNoCheckRole = append(routerNoCheckRole, clientRouter)
}

func clientRouter(v1 *gin.RouterGroup) {
	apiGroup := v1.Group("/api")
	client := apis.Client{}
	apiGroup.POST("/client/register", client.Register)
}
