package orm

import (
	"car/configs"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DB 定义全局数据库对象
var DB *gorm.DB

func init() {
	// 获取配置文件数据
	cfg := configs.GetGorm()
	adminCfg := configs.GetAdmin()
	db := adminCfg.Databases.GetDefault()
	var err error
	if DB, err = gorm.Open(mysql.Open(fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", db.User, db.Pwd, db.Host, db.Port, db.Name)), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{ // 命名策略
			SingularTable: cfg.SingularTable, // 单表去复数s
		},
		QueryFields: cfg.QueryFields,                     // 是否全字段映射
		Logger:      logger.Default.LogMode(logger.Info), // 日志级别
	}); err != nil {
		log.Println("gorm", err)
	}
}
