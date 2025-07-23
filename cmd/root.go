package cmd

import (
	"github.com/spf13/cobra"
	"github.com/youngprinnce/product-microservice/cmd/server"
)

var rootCmd = &cobra.Command{
	Use:   "product-microservice",
	Short: "Product Microservice API",
	Long:  `A gRPC-based product microservice with subscription management built with clean architecture`,
}

func Execute() {
	rootCmd.PersistentFlags().StringP("config", "c", "etc/config.yaml", "config filename")
	rootCmd.AddCommand(server.StartServerCmd())
	cobra.CheckErr(rootCmd.Execute())
}
