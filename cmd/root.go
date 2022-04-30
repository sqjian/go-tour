package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "root",
	Short: "root",
	Long:  `root cmd`,
}

var (
	CfgFile string
	Verbose bool
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&CfgFile, "config", ".go-tour.yml", "config file")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func initConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigFile(CfgFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Can't read config:", err)
		os.Exit(1)
	}
}
