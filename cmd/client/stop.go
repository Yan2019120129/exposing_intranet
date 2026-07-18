package client

import (
	"log"

	"github.com/spf13/cobra"
)

// A foreground client must receive SIGTERM to stop its runtime.
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop client",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("client is a foreground process; send SIGTERM to stop it")
	},
}
