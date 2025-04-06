package cmd

import (
	"github.com/spf13/cobra"
	"my-base/app"
	"os"
)

var rootCmd = &cobra.Command{
	Use:          "my-base",
	Short:        "my-base",
	SilenceUsage: true, // 默认使用
	Run: func(cmd *cobra.Command, args []string) {
		app.InitServer()
	},
}

// Execute 执行命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
