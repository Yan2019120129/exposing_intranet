package app

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"my-base/app/admin/router"
	"my-base/config"
	"my-base/core/gorm/orm"
)

// InitServer 启动项目服务
func InitServer() {
	cfg := config.GetGin()

	r := gin.Default()

	r.Use(func(context *gin.Context) {
		context.Set("db", orm.DB)
		context.Next()
	})
	router.InitRouter(r)
	//views.InitView(r)

	log.Printf("server start http://localhost%s%s", cfg.Port, "/login")
	if err := r.Run(cfg.Port); err != nil {
		zap.L().Error(err.Error())
	}
}
