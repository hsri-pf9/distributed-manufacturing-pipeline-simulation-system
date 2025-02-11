package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/domain"
)

// PipelineService orchestrates the execution of both sequential and parallel pipelines
type PipelineService struct {
	SequentialOrchestrator *domain.SequentialPipelineOrchestrator
	ParallelOrchestrator   *domain.ParallelPipelineOrchestrator
}

// NewPipelineService initializes the service
func NewPipelineService(sequential *domain.SequentialPipelineOrchestrator, parallel *domain.ParallelPipelineOrchestrator) *PipelineService {
	return &PipelineService{
		SequentialOrchestrator: sequential,
		ParallelOrchestrator:   parallel,
	}
}

// StartPipeline initializes and executes the pipeline
func (ps *PipelineService) StartPipeline(ctx context.Context, stages []domain.Stage, input interface{}, isParallel bool) (uuid.UUID, error) {
	pipelineID := uuid.New()

	var orchestrator domain.PipelineOrchestrator
	if isParallel {
		orchestrator = ps.ParallelOrchestrator
	} else {
		orchestrator = ps.SequentialOrchestrator
	}

	for _, stage := range stages {
		if err := orchestrator.AddStage(stage); err != nil {
			return pipelineID, err
		}
	}

	_, err := orchestrator.Execute(ctx, input)
	if err != nil {
		ps.setStatus(pipelineID, "Failed", isParallel)
		return pipelineID, err
	}

	ps.setStatus(pipelineID, "Completed", isParallel)
	return pipelineID, nil
}

// setStatus updates the pipeline execution status
func (ps *PipelineService) setStatus(pipelineID uuid.UUID, status string, isParallel bool) {
	if isParallel {
		ps.ParallelOrchestrator.UpdateStatus(pipelineID, status)
	}
}
