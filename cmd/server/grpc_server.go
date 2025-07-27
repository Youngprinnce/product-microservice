package server

import (
	"fmt"
	"log"
	"net"

	"github.com/youngprinnce/product-microservice/config"
	"github.com/youngprinnce/product-microservice/internal/auth"
	"github.com/youngprinnce/product-microservice/internal/grpc/handlers"
	"github.com/youngprinnce/product-microservice/internal/postgres"
	"github.com/youngprinnce/product-microservice/internal/service/product"
	"github.com/youngprinnce/product-microservice/internal/service/subscription"
	pb "github.com/youngprinnce/product-microservice/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGRPCServer(cfg *config.Config) {
	// Initialize database
	err := postgres.Load(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db := postgres.GetSession()

	// Auto-migrate database schema
	err = db.AutoMigrate(&product.Product{}, &subscription.SubscriptionPlan{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	// Initialize repositories
	productRepo := product.NewProductRepo(db)
	subscriptionRepo := subscription.NewSubscriptionRepo(db)

	// Initialize services
	productService := product.NewProductService(productRepo)
	subscriptionService := subscription.NewSubscriptionService(subscriptionRepo)

	// Initialize gRPC handlers
	productHandler := handlers.NewProductHandler(productService)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)

	// Initialize authentication
	authenticator := auth.NewAuthenticator()
	log.Printf("Basic authentication enabled. Available users: admin, client, test")

	// Create gRPC server with authentication interceptors
	server := grpc.NewServer(
		grpc.UnaryInterceptor(authenticator.UnaryInterceptor()),
		grpc.StreamInterceptor(authenticator.StreamInterceptor()),
	)

	// Register services
	pb.RegisterProductServiceServer(server, productHandler)
	pb.RegisterSubscriptionServiceServer(server, subscriptionHandler)

	// Enable reflection for grpcurl and other tools
	reflection.Register(server)

	// Create listener
	port := cfg.Server.Port
	if port == "" {
		port = "50051"
	}

	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	log.Printf("gRPC server starting on port %s", port)
	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
