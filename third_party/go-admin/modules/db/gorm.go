package db

import (
	"database/sql"
	"errors"

	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func OpenGorm(cfg config.Database) (*gorm.DB, error) {
	dialector, err := gormDialector(cfg)
	if err != nil {
		return nil, err
	}

	return gorm.Open(dialector, &gorm.Config{
		Logger: logger.GormLogger(),
	})
}

func OpenGormWithSQLDB(cfg config.Database, sqlDB *sql.DB) (*gorm.DB, error) {
	dialector, err := gormDialectorWithSQLDB(cfg, sqlDB)
	if err != nil {
		return nil, err
	}

	return gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{ // 命名策略
			SingularTable: true, // 单表去复数s
		},
		QueryFields: true, // 是否全字段映射
		Logger:      logger.GormLogger(),
	})
}

func gormDialector(cfg config.Database) (gorm.Dialector, error) {
	switch cfg.Driver {
	case DriverMysql, DriverOceanBase:
		return mysql.Open(cfg.GetDSN()), nil
	case DriverPostgresql:
		return postgres.Open(cfg.GetDSN()), nil
	case DriverSqlite:
		return sqlite.Open(cfg.GetDSN()), nil
	case DriverMssql:
		return sqlserver.Open(cfg.GetDSN()), nil
	default:
		return nil, errors.New("unsupported gorm driver: " + cfg.Driver)
	}
}

func gormDialectorWithSQLDB(cfg config.Database, sqlDB *sql.DB) (gorm.Dialector, error) {
	if sqlDB == nil {
		return gormDialector(cfg)
	}

	switch cfg.Driver {
	case DriverMysql, DriverOceanBase:
		return mysql.New(mysql.Config{Conn: sqlDB}), nil
	case DriverPostgresql:
		return postgres.New(postgres.Config{Conn: sqlDB}), nil
	case DriverSqlite:
		return sqlite.New(sqlite.Config{Conn: sqlDB}), nil
	case DriverMssql:
		return sqlserver.New(sqlserver.Config{Conn: sqlDB}), nil
	default:
		return nil, errors.New("unsupported gorm driver: " + cfg.Driver)
	}
}
