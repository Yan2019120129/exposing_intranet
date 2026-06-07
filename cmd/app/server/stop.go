package server

import (
	"github.com/spf13/cobra"
)

// StopCmd 关闭服务
var (
	isStop  bool
	StopCmd = &cobra.Command{
		Use:   "stop",
		Short: "stop server service",
		Long:  "stop intranet penetration server service",
		PreRun: func(cmd *cobra.Command, args []string) {
			// 关闭内网穿透服务
			if isStop {
			}
		},
	}
)

// init 初始化关闭服务参数
func init() {
	StopCmd.Flags().BoolVarP(&isStop, "server", "s", false, "stop server")
}
