package main

import (
	"log"
	"net"

	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/api/grpc/proto"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/api/grpc/server"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/adapters/secondary"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/domain"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/services"
	"google.golang.org/grpc"
)

const (
	grpcPort = ":50051" // gRPC server port
)

func main() {
	// Initialize database repository
	dbRepo := secondary.NewDatabaseAdapter()

	// Initialize pipeline orchestrators
	sequentialOrchestrator := domain.NewSequentialPipelineOrchestrator(dbRepo)
	parallelOrchestrator := domain.NewParallelPipelineOrchestrator(dbRepo)

	// Initialize pipeline service
	pipelineService := services.NewPipelineService(sequentialOrchestrator, parallelOrchestrator, dbRepo)

	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	pipelineGRPCServer := grpcserver.NewPipelineGRPCServer(pipelineService)

	// Register the gRPC service
	proto.RegisterPipelineServiceServer(grpcServer, pipelineGRPCServer)

	// Start listening on the defined port
	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	log.Printf("🚀 gRPC server is running on port %s", grpcPort)

	// Serve gRPC requests
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
