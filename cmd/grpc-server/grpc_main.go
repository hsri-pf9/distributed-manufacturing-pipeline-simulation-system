package main

import (
	"log"
	"net"
	"sync"

	proto "github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/api/grpc/proto/auth"
	pipeline_proto "github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/api/grpc/proto/pipeline"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/adapters/primary"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/adapters/secondary"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/services"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func startGRPCServer(authService *services.AuthService, pipelineService *services.PipelineService, wg *sync.WaitGroup) {
	defer wg.Done()

	// Create gRPC server
	grpcServer := grpc.NewServer()
	authServer := &primary.AuthServer{AuthService: authService}
	pipelineServer := &primary.PipelineServer{Service: pipelineService}

	// Register gRPC services
	proto.RegisterAuthServiceServer(grpcServer, authServer)
	pipeline_proto.RegisterPipelineServiceServer(grpcServer, pipelineServer)
	reflection.Register(grpcServer)

	// Start gRPC server
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("❌ Failed to listen on port 50051: %v", err)
	}

	log.Println("🚀 Starting gRPC server on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("❌ Failed to start gRPC server: %v", err)
	}
}

func main() {
	// Initialize database
	secondary.InitDatabase()
	defer secondary.CloseDatabase()

	dbRepo := secondary.NewDatabaseAdapter()

	sseManager := utils.NewSSEManager()

	// Initialize services
	authService := services.NewAuthService(dbRepo)
	pipelineService := services.NewPipelineService(dbRepo, sseManager) // gRPC does not need SSE

	var wg sync.WaitGroup
	wg.Add(1) // Only 1 (gRPC)

	// Start gRPC server
	go startGRPCServer(authService, pipelineService, &wg)

	// Wait for server
	wg.Wait()
}
