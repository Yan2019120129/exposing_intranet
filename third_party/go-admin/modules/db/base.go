package db

import (
	"database/sql"
	"errors"
	"sync"

	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/redis"
	"gorm.io/gorm"
)

// Base is a common Connection.
type Base struct {
	DbList   map[string]*sql.DB
	GormList map[string]*gorm.DB
	Once     sync.Once
	Configs  config.DatabaseList
}

// Close implements the method Connection.Close.
func (db *Base) Close() []error {
	errs := make([]error, 0)
	for _, d := range db.DbList {
		errs = append(errs, d.Close())
	}
	if err := redis.Close(); err != nil {
		errs = append(errs, err)
	}
	return errs
}

// GetDB implements the method Connection.GetDB.
func (db *Base) GetDB(key string) *sql.DB {
	return db.DbList[key]
}

func (db *Base) GetGorm(key string) (*gorm.DB, error) {
	if key == "" {
		key = "default"
	}

	if gormDB, ok := db.GormList[key]; ok && gormDB != nil {
		return gormDB, nil
	}
	return nil, errors.New("wrong connection name")
}

func (db *Base) initGorm(conn string, cfg config.Database, sqlDB *sql.DB) {
	if db.GormList == nil {
		db.GormList = make(map[string]*gorm.DB)
	}

	gormDB, err := OpenGormWithSQLDB(cfg, sqlDB)
	if err != nil {
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
		panic(err)
	}
	db.GormList[conn] = gormDB
}

func (db *Base) CreateDB(name string, beans ...interface{}) error {
	cfg := db.GetConfig(name)
	if cfg.Driver == "" {
		return errors.New("wrong connection name")
	}
	gormDB, err := db.GetGorm(name)
	if err != nil {
		return err
	}
	return gormDB.AutoMigrate(beans...)
}

func (db *Base) GetConfig(name string) config.Database {
	return db.Configs[name]
}
