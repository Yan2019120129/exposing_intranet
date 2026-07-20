package cmd

import (
	"my-base/app"
	"my-base/cmd/app/console"
	"my-base/cmd/app/server"
	"my-base/cmd/client"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:   "base",
	Short: "b",
	Long:  "my base",
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.StartServer()
	},
}

// 初始化项目命令
func init() {
	rootCmd.AddCommand(server.StartCmd)
	rootCmd.AddCommand(server.StopCmd)
	rootCmd.AddCommand(console.CreateCmd)
	rootCmd.AddCommand(console.DeleteCmd)
	rootCmd.AddCommand(console.ResetCmd)
	rootCmd.AddCommand(client.NewCommand())
}

// Execute 初始化命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
