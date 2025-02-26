package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	proto "github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/api/grpc/proto"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/adapters/secondary"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/services"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/adapters/primary"
)

func main() {
	// Initialize database and Supabase client
	secondary.InitDatabase()
	dbRepo := secondary.NewDatabaseAdapter()
	authService := services.NewAuthService(dbRepo)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	authServer := &primary.AuthServer{AuthService: authService}

	// Register gRPC service
	proto.RegisterAuthServiceServer(grpcServer, authServer)
	reflection.Register(grpcServer)

	// Start gRPC server
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	log.Println("Starting gRPC server on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
