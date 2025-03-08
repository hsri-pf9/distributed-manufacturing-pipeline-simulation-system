package main

import (
	"context"
	"fmt"
	"log"
	"time"

	proto "github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/api/grpc/proto/pipeline"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"google.golang.org/protobuf/types/known/anypb"
)

// Root pipeline command
var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Manage pipelines",
}

// ✅ Create Pipeline Command
var createPipelineCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new pipeline",
	Run: func(cmd *cobra.Command, args []string) {
		userID, _ := cmd.Flags().GetString("user")
		stages, _ := cmd.Flags().GetInt("stages")
		isParallel, _ := cmd.Flags().GetBool("parallel")

		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to gRPC server: %v", err)
		}
		defer conn.Close()

		client := proto.NewPipelineServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		resp, err := client.CreatePipeline(ctx, &proto.CreatePipelineRequest{
			UserId:    userID,
			Stages:    int32(stages),
			IsParallel: isParallel,
		})
		if err != nil {
			log.Fatalf("Pipeline creation failed: %v", err)
		}

		fmt.Printf("✅ Pipeline created successfully! Pipeline ID: %s\n", resp.PipelineId)
	},
}

// ✅ Start Pipeline Command
var startPipelineCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a pipeline execution",
	Run: func(cmd *cobra.Command, args []string) {
		pipelineID, _ := cmd.Flags().GetString("pipeline-id")
		userID, _ := cmd.Flags().GetString("user-id")
		inputValue, _ := cmd.Flags().GetString("input") // ✅ Get input from CLI
		isParallel, _ := cmd.Flags().GetBool("parallel")

		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to gRPC server: %v", err)
		}
		defer conn.Close()

		client := proto.NewPipelineServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		// ✅ Convert input string to Google Protobuf `StringValue`
		stringValue := &wrapperspb.StringValue{Value: inputValue}

		// ✅ Convert `StringValue` to `Any`
		anyInput, err := anypb.New(stringValue)
		if err != nil {
			log.Fatalf("Failed to wrap input in Any: %v", err)
		}

		resp, err := client.StartPipeline(ctx, &proto.StartPipelineRequest{
			PipelineId: pipelineID,
			UserId:     userID,
			Input:      anyInput, // ✅ Correctly passing as Any
			IsParallel: isParallel,
		})
		if err != nil {
			log.Fatalf("Failed to start pipeline: %v", err)
		}

		fmt.Printf("🚀 Pipeline execution started successfully! Message: %s\n", resp.Message)
	},
}



// ✅ Cancel Pipeline Command
var cancelPipelineCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel a running pipeline",
	Run: func(cmd *cobra.Command, args []string) {
		pipelineID, _ := cmd.Flags().GetString("pipeline-id")
		userID, _ := cmd.Flags().GetString("user-id")
		isParallel, _ := cmd.Flags().GetBool("parallel")

		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to gRPC server: %v", err)
		}
		defer conn.Close()

		client := proto.NewPipelineServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		resp, err := client.CancelPipeline(ctx, &proto.CancelPipelineRequest{
			PipelineId: pipelineID,
			UserId:     userID,
			IsParallel: isParallel,
		})
		if err != nil {
			log.Fatalf("Failed to cancel pipeline: %v", err)
		}

		fmt.Printf("❌ Pipeline cancelled successfully! Message: %s\n", resp.Message)
	},
}

// ✅ Get Pipeline Status Command
var getPipelineStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of a pipeline",
	Run: func(cmd *cobra.Command, args []string) {
		pipelineID, _ := cmd.Flags().GetString("pipeline-id")
		isParallel, _ := cmd.Flags().GetBool("parallel")

		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to gRPC server: %v", err)
		}
		defer conn.Close()

		client := proto.NewPipelineServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		resp, err := client.GetPipelineStatus(ctx, &proto.GetPipelineStatusRequest{
			PipelineId: pipelineID,
			IsParallel: isParallel,
		})
		if err != nil {
			log.Fatalf("Failed to get pipeline status: %v", err)
		}

		fmt.Printf("📊 Pipeline Status: %s\n", resp.Status)
	},
}

func init() {
	// Add all commands under `pipeline`
	pipelineCmd.AddCommand(createPipelineCmd)
	pipelineCmd.AddCommand(startPipelineCmd)
	pipelineCmd.AddCommand(cancelPipelineCmd)
	pipelineCmd.AddCommand(getPipelineStatusCmd)

	// Flags for create pipeline
	createPipelineCmd.Flags().String("user", "", "User ID")
	createPipelineCmd.Flags().Int("stages", 3, "Number of stages")
	createPipelineCmd.Flags().Bool("parallel", false, "Parallel execution")
	createPipelineCmd.MarkFlagRequired("user")

	// Flags for start pipeline
	startPipelineCmd.Flags().String("pipeline-id", "", "Pipeline ID")
	startPipelineCmd.Flags().String("user-id", "", "User ID")
	startPipelineCmd.Flags().String("input", "", "Input for pipeline")
	startPipelineCmd.Flags().Bool("parallel", false, "Run in parallel mode")
	startPipelineCmd.MarkFlagRequired("pipeline-id")
	startPipelineCmd.MarkFlagRequired("user-id")

	// Flags for cancel pipeline
	cancelPipelineCmd.Flags().String("pipeline-id", "", "Pipeline ID")
	cancelPipelineCmd.Flags().String("user-id", "", "User ID")
	cancelPipelineCmd.Flags().Bool("parallel", false, "Cancel parallel pipeline")
	cancelPipelineCmd.MarkFlagRequired("pipeline-id")
	cancelPipelineCmd.MarkFlagRequired("user-id")

	// Flags for get pipeline status
	getPipelineStatusCmd.Flags().String("pipeline-id", "", "Pipeline ID")
	getPipelineStatusCmd.Flags().Bool("parallel", false, "Check parallel pipeline status")
	getPipelineStatusCmd.MarkFlagRequired("pipeline-id")
}
