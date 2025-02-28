package primary

import (
	"context"
	"log"

	"github.com/google/uuid"
	proto "github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/api/grpc/proto/pipeline"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	// "google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// PipelineServer implements the gRPC pipeline service
type PipelineServer struct {
	proto.UnimplementedPipelineServiceServer
	Service *services.PipelineService
}

// CreatePipeline handles gRPC pipeline creation
func (s *PipelineServer) CreatePipeline(ctx context.Context, req *proto.CreatePipelineRequest) (*proto.CreatePipelineResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user ID: %v", err)
	}

	pipelineID, err := s.Service.CreatePipeline(userID, int(req.Stages), req.IsParallel)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create pipeline: %v", err)
	}

	return &proto.CreatePipelineResponse{PipelineId: pipelineID.String()}, nil
}

// func (s *PipelineServer) StartPipeline(ctx context.Context, req *proto.StartPipelineRequest) (*proto.StartPipelineResponse, error) {
// 	// Parse pipeline ID
// 	pipelineID, err := uuid.Parse(req.PipelineId)
// 	if err != nil {
// 		return nil, status.Errorf(codes.InvalidArgument, "Invalid pipeline ID: %v", err)
// 	}

// 	// Parse user ID
// 	userID, err := uuid.Parse(req.UserId)
// 	if err != nil {
// 		return nil, status.Errorf(codes.InvalidArgument, "Invalid user ID: %v", err)
// 	}

// 	// Ensure user ID is not empty
// 	if userID == uuid.Nil {
// 		return nil, status.Error(codes.InvalidArgument, "User ID is required")
// 	}

// 	// Process input data correctly
// 	var input interface{}
// 	if req.Input != nil {
// 		// Try to unpack as a StringValue
// 		stringWrapper := &wrapperspb.StringValue{}
// 		if err := req.Input.UnmarshalTo(stringWrapper); err == nil {
// 			input = stringWrapper.Value
// 		} else {
// 			// Try to unpack as a struct (JSON object)
// 			structValue := &structpb.Struct{}
// 			if err := req.Input.UnmarshalTo(structValue); err == nil {
// 				input = structValue.AsMap()
// 			} else {
// 				log.Printf("Failed to unpack input: %v", err)
// 				return nil, status.Errorf(codes.InvalidArgument, "Invalid input format: %v", err)
// 			}
// 		}
// 	}

// 	// Run the pipeline asynchronously
// 	go func() {
// 		err := s.Service.StartPipeline(context.Background(), userID, pipelineID, input, req.IsParallel)
// 		if err != nil {
// 			log.Printf("Failed to start pipeline %s: %v", pipelineID, err)
// 		}
// 	}()

// 	// Return response
// 	return &proto.StartPipelineResponse{
// 		Message: "Pipeline execution started",
// 	}, nil
// }

func (s *PipelineServer) StartPipeline(ctx context.Context, req *proto.StartPipelineRequest) (*proto.StartPipelineResponse, error) {
	log.Println("[GRPC] Received StartPipeline request...")

	// Parse pipeline ID
	pipelineID, err := uuid.Parse(req.PipelineId)
	if err != nil {
		log.Printf("[ERROR] Invalid pipeline ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid pipeline ID: %v", err)
	}

	// Parse user ID
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		log.Printf("[ERROR] Invalid user ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user ID: %v", err)
	}

	// Ensure user ID is not empty
	if userID == uuid.Nil {
		log.Println("[ERROR] User ID is required but received nil")
		return nil, status.Error(codes.InvalidArgument, "User ID is required")
	}

	// Process input data
	var input interface{}
	if req.Input != nil {
		log.Println("[DEBUG] Processing input data...")

		// Try to unpack as StringValue
		stringWrapper := &wrapperspb.StringValue{}
		if err := req.Input.UnmarshalTo(stringWrapper); err == nil {
			input = stringWrapper.Value
			log.Printf("[DEBUG] Parsed input as string: %v", input)
		} else {
			// Try to unpack as Struct (JSON object)
			structValue := &structpb.Struct{}
			if err := req.Input.UnmarshalTo(structValue); err == nil {
				input = structValue.AsMap()
				log.Printf("[DEBUG] Parsed input as JSON object: %v", input)
			} else {
				log.Printf("[ERROR] Failed to unpack input: %v", err)
				return nil, status.Errorf(codes.InvalidArgument, "Invalid input format: %v", err)
			}
		}
	} else {
		log.Println("[DEBUG] No input data provided, using nil")
	}

	// Run the pipeline asynchronously
	go func() {
		log.Printf("[INFO] Starting pipeline execution: %s", pipelineID)
		err := s.Service.StartPipeline(context.Background(), userID, pipelineID, input, req.IsParallel)
		if err != nil {
			log.Printf("[ERROR] Pipeline execution failed for %s: %v", pipelineID, err)
		} else {
			log.Printf("[INFO] Pipeline execution completed successfully: %s", pipelineID)
		}
	}()

	// Return response
	return &proto.StartPipelineResponse{
		Message: "Pipeline execution started",
	}, nil
}


// GetPipelineStatus retrieves the current pipeline status
func (s *PipelineServer) GetPipelineStatus(ctx context.Context, req *proto.GetPipelineStatusRequest) (*proto.GetPipelineStatusResponse, error) {
	pipelineID, err := uuid.Parse(req.PipelineId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid pipeline ID: %v", err)
	}

	stat, err := s.Service.GetPipelineStatus(pipelineID, req.IsParallel)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get pipeline status: %v", err)
	}

	return &proto.GetPipelineStatusResponse{
		PipelineId: pipelineID.String(),
		Status:     stat,
	}, nil
}

// CancelPipeline cancels an ongoing pipeline execution
func (s *PipelineServer) CancelPipeline(ctx context.Context, req *proto.CancelPipelineRequest) (*proto.CancelPipelineResponse, error) {
	pipelineID, err := uuid.Parse(req.PipelineId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid pipeline ID: %v", err)
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user ID: %v", err)
	}

	err = s.Service.CancelPipeline(pipelineID, userID, req.IsParallel)
	if err != nil {
		log.Printf("Error cancelling pipeline %s: %v", pipelineID, err)
		return nil, status.Errorf(codes.Internal, "Failed to cancel pipeline: %v", err)
	}

	return &proto.CancelPipelineResponse{Message: "Pipeline cancelled"}, nil
}
