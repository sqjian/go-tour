package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sqjian/go-tour/internal/rpc"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server",
	Long:  `start server`,
	Run: func(cmd *cobra.Command, args []string) {
		grpcAddr, _ := cmd.Flags().GetString("grpc")
		gatewayAddr, _ := cmd.Flags().GetString("gateway")
		if err := rpc.StartSrv(grpcAddr, gatewayAddr); err != nil {
			panic(err)
		}
	},
}

const (
	grpcAddr    = "127.0.0.1:50051"
	gatewayAddr = "127.0.0.1:50052"
)

func init() {
	serverCmd.Flags().String("grpc", grpcAddr, fmt.Sprintf("set grpc addr"))
	serverCmd.Flags().String("gateway", gatewayAddr, fmt.Sprintf("set gateway addr"))
}
