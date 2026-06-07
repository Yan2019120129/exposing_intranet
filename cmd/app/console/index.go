package console

import (
	"car/app/models"

	"github.com/spf13/cobra"
)

var (
	CreateCmd = &cobra.Command{
		Use:   "create",
		Short: "create table",
		Long:  "create table to database,",
		RunE: func(cmd *cobra.Command, args []string) error {
			return models.ModelManage.Create(args...)
		},
	}

	DeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "delete table",
		Long:  "delete table to database,",
		RunE: func(cmd *cobra.Command, args []string) error {
			return models.ModelManage.Delete(args...)
		},
	}

	ResetCmd = &cobra.Command{
		Use:   "reset",
		Short: "reset table",
		Long:  "reset table to database,",
		RunE: func(cmd *cobra.Command, args []string) error {
			return models.ModelManage.Reset(args...)
		},
	}
)
