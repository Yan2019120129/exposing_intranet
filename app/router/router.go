package router

import (
	"car/app/apis"

	"github.com/gin-gonic/gin"
)

func init() {
	routerNoCheckRole = append(routerNoCheckRole, RegisterTestRouter)
}

// RegisterTestRouter initializes test routes.
func RegisterTestRouter(r *gin.RouterGroup) {
	test := apis.Test{}

	r.GET("/test", test.Ping)
	r.GET("/test-page", test.Page)
	r.GET("/tests", test.List)
	r.POST("/tests", test.Create)
	r.GET("/tests/:id", test.Get)
	r.PUT("/tests/:id", test.Update)
	r.DELETE("/tests/:id", test.Delete)
}
