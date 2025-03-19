package cmd

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
)

// GetGRPCConnection establishes a connection to the gRPC server
func GetGRPCConnection() (*grpc.ClientConn, context.Context, context.CancelFunc) {
	// Get gRPC server address from environment variable
	grpcServerURL := os.Getenv("DEMOCTL_GRPC_URL")
	if grpcServerURL == "" {
		log.Fatal("❌ DEMOCTL_GRPC_URL is not set. Please export it before running democtl.")
	}

	// Establish gRPC connection
	conn, err := grpc.Dial(grpcServerURL, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("❌ Failed to connect to gRPC server: %v", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	return conn, ctx, cancel
}
