package client

import "github.com/spf13/cobra"

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start client service",
	RunE: func(cmd *cobra.Command, args []string) error {
		return StartClient()
	},
}
