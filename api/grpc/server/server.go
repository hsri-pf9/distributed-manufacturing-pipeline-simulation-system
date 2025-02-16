package grpcserver

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	pb "github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/api/grpc/proto"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/services"

	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// PipelineGRPCServer implements pb.PipelineServiceServer
type PipelineGRPCServer struct {
	pb.UnimplementedPipelineServiceServer
	pipelineService *services.PipelineService
}

// NewPipelineGRPCServer initializes a new gRPC server instance
func NewPipelineGRPCServer(pipelineService *services.PipelineService) *PipelineGRPCServer {
	return &PipelineGRPCServer{pipelineService: pipelineService}
}

// CreatePipeline handles creating a new pipeline
func (s *PipelineGRPCServer) CreatePipeline(ctx context.Context, req *pb.CreatePipelineRequest) (*pb.CreatePipelineResponse, error) {
	log.Println("Received CreatePipeline request")

	if len(req.Stages) == 0 {
		return nil, grpcstatus.Error(codes.InvalidArgument, "Pipeline must have at least one stage")
	}

	pipelineID, err := s.pipelineService.CreatePipeline(len(req.Stages), false) // Assuming sequential execution
	if err != nil {
		return nil, grpcstatus.Errorf(codes.Internal, "Failed to create pipeline: %v", err)
	}

	return &pb.CreatePipelineResponse{PipelineId: pipelineID.String()}, nil
}

// ExecutePipeline starts execution of an existing pipeline
func (s *PipelineGRPCServer) ExecutePipeline(ctx context.Context, req *pb.ExecutePipelineRequest) (*pb.ExecutePipelineResponse, error) {
	log.Printf("Received ExecutePipeline request for PipelineID: %s", req.PipelineId)

	pipelineUUID, err := uuid.Parse(req.PipelineId)
	if err != nil {
		return nil, grpcstatus.Error(codes.InvalidArgument, "Invalid Pipeline ID format")
	}

	err = s.pipelineService.StartPipeline(ctx, pipelineUUID, req.InputData, false) // Assuming sequential execution
	if err != nil {
		return nil, grpcstatus.Errorf(codes.Internal, "Failed to execute pipeline: %v", err)
	}

	return &pb.ExecutePipelineResponse{
		PipelineId: req.PipelineId,
		OutputData: "Execution started successfully",
	}, nil
}

// GetPipelineStatus returns the status of a pipeline
func (s *PipelineGRPCServer) GetPipelineStatus(ctx context.Context, req *pb.GetPipelineStatusRequest) (*pb.GetPipelineStatusResponse, error) {
	log.Printf("Received GetPipelineStatus request for PipelineID: %s", req.PipelineId)

	pipelineUUID, err := uuid.Parse(req.PipelineId)
	if err != nil {
		return nil, grpcstatus.Error(codes.InvalidArgument, "Invalid Pipeline ID format")
	}

	status, err := s.pipelineService.GetPipelineStatus(pipelineUUID)
	if err != nil {
		return nil, grpcstatus.Errorf(codes.Internal, "Failed to get pipeline status: %v", err)
	}

	return &pb.GetPipelineStatusResponse{
		PipelineId: req.PipelineId,
		Status:     status,
		UpdatedAt:  timestamppb.New(time.Now()), // Assuming latest update
	}, nil
}

// CancelPipeline stops an active pipeline
func (s *PipelineGRPCServer) CancelPipeline(ctx context.Context, req *pb.CancelPipelineRequest) (*emptypb.Empty, error) {
	log.Printf("Received CancelPipeline request for PipelineID: %s", req.PipelineId)

	pipelineUUID, err := uuid.Parse(req.PipelineId)
	if err != nil {
		return nil, grpcstatus.Error(codes.InvalidArgument, "Invalid Pipeline ID format")
	}

	err = s.pipelineService.CancelPipeline(pipelineUUID)
	if err != nil {
		return nil, grpcstatus.Errorf(codes.Internal, "Failed to cancel pipeline: %v", err)
	}

	return &emptypb.Empty{}, nil
}
