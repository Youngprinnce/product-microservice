package server

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/youngprinnce/product-microservice/config"
	"github.com/youngprinnce/product-microservice/internal/logger"
	"github.com/youngprinnce/product-microservice/internal/postgres"
)

func StartServerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Start the gRPC server",
		Long:  `Start the gRPC server for product and subscription services`,
		Run: func(cmd *cobra.Command, args []string) {
			configFile, _ := cmd.Flags().GetString("config")
			conf := config.LoadConfig(configFile)

			logger.Initialize()

			if err := postgres.Load(conf); err != nil {
				logger.Fatal(fmt.Sprintf("Failed to initialize postgres: %v", err))
			}

			log.WithField("port", conf.Server.Port).Info("Starting gRPC server")

			// Start the gRPC server
			StartGRPCServer(conf)
		},
	}
}
