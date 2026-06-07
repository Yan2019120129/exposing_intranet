package console

import (
	"my-base/app/models"
	"my-base/configs"

	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var (
	CreateCmd = &cobra.Command{
		Use:   "create",
		Short: "create table",
		Long:  "create table to database,",
		RunE: func(cmd *cobra.Command, args []string) error {
			return withEngineGorm(func(db *gorm.DB) error {
				return models.ModelManage.Create(db, args...)
			})
		},
	}

	DeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "delete table",
		Long:  "delete table to database,",
		RunE: func(cmd *cobra.Command, args []string) error {
			return withEngineGorm(func(db *gorm.DB) error {
				return models.ModelManage.Delete(db, args...)
			})
		},
	}

	ResetCmd = &cobra.Command{
		Use:   "reset",
		Short: "reset table",
		Long:  "reset table to database,",
		RunE: func(cmd *cobra.Command, args []string) error {
			return withEngineGorm(func(db *gorm.DB) error {
				return models.ModelManage.Reset(db, args...)
			})
		},
	}
)

func withEngineGorm(fn func(*gorm.DB) error) error {
	eng := engine.Default().AddConfig(configs.GetAdmin())
	defer eng.DefaultConnection().Close()
	db, err := eng.DefaultConnection().GetGorm("default")
	if err != nil {
		return err
	}
	return fn(db)
}
