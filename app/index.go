package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"my-base/app/router"
	"my-base/configs"
	orm "my-base/module/gorm"
	"my-base/pages"
	"my-base/tables"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/modules/logger"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/gin-gonic/gin"
)

func StartServer() {
	cfg := configs.GetConfig()
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	r := gin.Default()
	r.LoadHTMLGlob("html/*")

	template.AddComp(chartjs.NewChart())

	eng := engine.Default()

	if err := eng.AddConfig(configs.GetAdmin()).
		AddGenerators(tables.Generators).
		Use(r); err != nil {
		panic(err)
	}

	r.Use(func(c *gin.Context) {
		db := orm.DB
		c.Set("db", db)
	})

	r.Static(cfg.Admin.AssetRootPath, cfg.Admin.Store.Path)

	router.InitRouter(r)

	eng.HTML("GET", cfg.Admin.Prefix(), pages.GetDashBoard)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		fmt.Println("--------- http://localhost" + srv.Addr + cfg.Admin.Prefix() + " ---------")
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			logger.Info("listen:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", err)
	}
	logger.Info("closing database connection")
	eng.MysqlConnection().Close()

	logger.Info("Server exiting")
}
