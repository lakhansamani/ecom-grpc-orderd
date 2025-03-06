package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	order "github.com/lakhansamani/ecom-grpc-apis/order/v1"
	user "github.com/lakhansamani/ecom-grpc-apis/user/v1"

	"github.com/lakhansamani/ecom-grpc-orderd/db"
	"github.com/lakhansamani/ecom-grpc-orderd/service"
)

func main() {
	// Read .env file as environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using environment variables")
	}

	// DB URL
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is required")
	}
	// Initialize database
	dbProvider := db.New(dbURL)

	// Get User Service URL
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		log.Fatal("USER_SERVICE_URL is required")
	}

	// Create UserServiceClient using grpc
	grpcConn, err := grpc.NewClient(userServiceURL, grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))
	if err != nil {
		log.Fatalf("Failed to dial UserService: %v", err)
	}
	defer grpcConn.Close()

	userServiceClient := user.NewUserServiceClient(grpcConn)

	// Create a new gRPC server
	server := grpc.NewServer()

	// Register OrderService with gRPC
	orderService := service.New(
		service.Config{},
		service.Dependencies{
			DBProvider:  dbProvider,
			UserService: userServiceClient,
		})
	order.RegisterOrderServiceServer(server, orderService)

	// Start gRPC server
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("gRPC Server is running on port 50052...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
