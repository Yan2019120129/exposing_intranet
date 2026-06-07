package cmd

import (
	"car/app"
	"car/cmd/app/console"
	"car/cmd/app/server"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:   "base",
	Short: "b",
	Long:  "my base",
	Run: func(cmd *cobra.Command, args []string) {
		app.StartServer()
	},
}

// 初始化项目命令
func init() {
	rootCmd.AddCommand(server.StartCmd)
	rootCmd.AddCommand(server.StopCmd)
	rootCmd.AddCommand(console.CreateCmd)
	rootCmd.AddCommand(console.DeleteCmd)
	rootCmd.AddCommand(console.ResetCmd)
}

// Execute 初始化命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
