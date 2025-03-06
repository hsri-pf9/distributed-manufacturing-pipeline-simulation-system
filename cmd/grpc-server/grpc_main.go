package main

// import (
// 	"log"
// 	"net"

// 	proto "github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/api/grpc/proto/auth"
// 	pipeline_proto "github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/api/grpc/proto/pipeline"
// 	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/adapters/primary"
// 	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/adapters/secondary"
// 	// "github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/domain"
// 	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/services"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/reflection"
// )

// func main() {
// 	// Initialize database and Supabase client
// 	secondary.InitDatabase()
// 	dbRepo := secondary.NewDatabaseAdapter()
// 	authService := services.NewAuthService(dbRepo)

// 	// Initialize pipeline service
// 	pipelineService := services.NewPipelineService(dbRepo)

// 	// Create gRPC server
// 	grpcServer := grpc.NewServer()
// 	authServer := &primary.AuthServer{AuthService: authService}
// 	pipelineServer := &primary.PipelineServer{Service: pipelineService}

// 	// Register gRPC service
// 	proto.RegisterAuthServiceServer(grpcServer, authServer)
// 	pipeline_proto.RegisterPipelineServiceServer(grpcServer, pipelineServer)
// 	reflection.Register(grpcServer)

// 	// Start gRPC server
// 	listener, err := net.Listen("tcp", ":50051")
// 	if err != nil {
// 		log.Fatalf("Failed to listen on port 50051: %v", err)
// 	}

// 	log.Println("Starting gRPC server on port 50051...")
// 	if err := grpcServer.Serve(listener); err != nil {
// 		log.Fatalf("Failed to start gRPC server: %v", err)
// 	}
// }
