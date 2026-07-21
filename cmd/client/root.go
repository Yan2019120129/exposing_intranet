package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"my-base/client/apis"
	"my-base/client/configs"
	clientservice "my-base/client/service"
	"my-base/code/contract"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "intranet penetration client",
	RunE: func(cmd *cobra.Command, args []string) error {
		return StartClient()
	},
}

var authCmd = &cobra.Command{
	Use:   "auth <username:password>",
	Short: "Authenticate client with server",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return authenticate(cmd.Context(), args[0])
	},
}

var mapCmd = &cobra.Command{Use: "map", Short: "Manage port mappings"}

var mapAddCmd = &cobra.Command{
	Use:   "add <server_port> <local_addr> [comment]",
	Short: "Add a port mapping",
	Args:  cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		comment := ""
		if len(args) == 3 {
			comment = args[2]
		}
		return mapPort(cmd.Context(), "add", args[0], args[1], comment)
	},
}

var mapDelCmd = &cobra.Command{
	Use:   "del <server_port>",
	Short: "Delete a port mapping",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return mapPort(cmd.Context(), "del", args[0], "", "")
	},
}

var mapListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all port mappings",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return mapPort(cmd.Context(), "list", "", "", "")
	},
}

func init() {
	mapCmd.AddCommand(mapAddCmd, mapDelCmd, mapListCmd)
	rootCmd.AddCommand(authCmd, mapCmd, startCmd, stopCmd)
}

// NewCommand creates the client subcommand for the server root command.
func NewCommand() *cobra.Command { return rootCmd }

// StartClient runs the client in the foreground until it receives a shutdown
// signal.
func StartClient() error {
	cfg := configs.GetConfig()
	runtime := clientservice.NewRuntime(cfg)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	defer runtime.Close()
	if err := runtime.Run(ctx); err != nil && !errors.Is(err, clientservice.ErrClientDeleted) {
		return fmt.Errorf("client runtime: %w", err)
	}
	log.Println("client stopped")
	return nil
}

func authenticate(ctx context.Context, authArg string) error {
	cfg := configs.GetConfig()
	api := apis.NewHTTPAPI(cfg.GetClientConfig())
	service := clientservice.NewAuthService(api, cfg)
	symbol, err := service.Authenticate(ctx, authArg)
	if err != nil {
		return err
	}
	log.Printf("认证成功，symbol: %s", symbol)
	return nil
}

func mapPort(ctx context.Context, action, serverPort, localAddr, comment string) error {
	cfg := configs.GetConfig()
	api := apis.NewHTTPAPI(cfg.GetClientConfig())
	service := clientservice.NewPortService(api, cfg)
	result, err := service.Manage(ctx, action, serverPort, localAddr, comment)
	if err != nil {
		return err
	}
	printPortResult(action, serverPort, localAddr, comment, result)
	return nil
}

func printPortResult(action, serverPort, localAddr, comment string, result contract.PortResponse) {
	switch action {
	case "add":
		if comment != "" {
			log.Printf("✓ 端口映射添加成功: :%s -> %s [%s]", serverPort, localAddr, comment)
		} else {
			log.Printf("✓ 端口映射添加成功: :%s -> %s", serverPort, localAddr)
		}
	case "del":
		log.Printf("✓ 端口映射删除成功: :%s", serverPort)
	case "list":
		if len(result.Data) == 0 {
			log.Println("暂无端口映射")
			return
		}
		log.Println("端口映射列表:")
		for _, mapping := range result.Data {
			log.Printf("- :%s -> %s [%s] %s", mapping.ServerPort, mapping.LocalAddr, mapping.Comment, mapping.Status)
		}
	}
}
