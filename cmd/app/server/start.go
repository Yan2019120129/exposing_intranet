package server

import (
	"github.com/spf13/cobra"
)

var (
	port     string
	isStart  bool
	StartCmd = &cobra.Command{
		Use:   "start",
		Short: "start server service",
		Long:  "start car server service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if isStart {
			}
			return nil
		},
	}
)

// init 初始化启动命令
func init() {
	StartCmd.Flags().StringVarP(&port, "port", "p", "1010", "set port")
	StartCmd.Flags().BoolVarP(&isStart, "server", "s", false, "start server")
}
