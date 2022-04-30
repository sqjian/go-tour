package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sqjian/go-tour/internal/rpc"
)

func init() {
	rootCmd.AddCommand(clientCmd)
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "client",
	Long:  `start client`,
	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("addr")
		port, _ := cmd.Flags().GetString("port")
		if err := rpc.StartCli(addr, port); err != nil {
			panic(err)
		}
	},
}

func init() {
	clientCmd.Flags().String("addr", "127.0.0.1", "set addr")
	clientCmd.Flags().String("port", "50051", "set port")
}
