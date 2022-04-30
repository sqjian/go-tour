package cmd

import (
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
		addr, _ := cmd.Flags().GetString("addr")
		port, _ := cmd.Flags().GetString("port")
		if err := rpc.StartSrv(addr, port); err != nil {
			panic(err)
		}
	},
}

func init() {
	serverCmd.Flags().String("addr", "127.0.0.1", "set addr")
	serverCmd.Flags().String("port", "50051", "set port")
}
