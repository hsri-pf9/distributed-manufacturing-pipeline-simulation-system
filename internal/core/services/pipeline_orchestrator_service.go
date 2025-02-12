package services

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/domain"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/ports"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/models"
)

type PipelineService struct {
	SequentialOrchestrator *domain.SequentialPipelineOrchestrator
	ParallelOrchestrator   *domain.ParallelPipelineOrchestrator
	Repository             ports.PipelineRepository
}

func NewPipelineService(sequential *domain.SequentialPipelineOrchestrator, parallel *domain.ParallelPipelineOrchestrator, repo ports.PipelineRepository) *PipelineService {
	return &PipelineService{
		SequentialOrchestrator: sequential,
		ParallelOrchestrator:   parallel,
		Repository:             repo,
	}
}

func (ps *PipelineService) CreatePipeline(stageCount int, isParallel bool) (uuid.UUID, error) {
	pipelineID := uuid.New()

	var orchestrator domain.PipelineOrchestrator
	if isParallel {
		orchestrator = ps.ParallelOrchestrator
	} else {
		orchestrator = ps.SequentialOrchestrator
	}

	
	for i := 0; i < stageCount; i++ {
		stage := &domain.BaseStage{ID: uuid.New()}
		if err := orchestrator.AddStage(stage); err != nil {
			return uuid.Nil, err
		}
	}

	err := ps.Repository.SavePipelineExecution(&models.PipelineExecution{
		ID:         pipelineID,
		PipelineID: pipelineID,
		Status:     "Created",
	})
	if err != nil {
		return uuid.Nil, err
	}

	return pipelineID, nil
}

func (ps *PipelineService) StartPipeline(ctx context.Context, pipelineID uuid.UUID, input interface{}, isParallel bool) error {
	var orchestrator domain.PipelineOrchestrator
	if isParallel {
		orchestrator = ps.ParallelOrchestrator
	} else {
		orchestrator = ps.SequentialOrchestrator
	}

	if err := ps.updatePipelineStatus(pipelineID, "Running"); err != nil {
		return err
	}

	_, err := orchestrator.Execute(ctx, input)
	if err != nil {
		_ = ps.updatePipelineStatus(pipelineID, "Failed")
		ps.logExecutionError(pipelineID, err.Error())
		return err
	}

	return ps.updatePipelineStatus(pipelineID, "Completed")
}

func (ps *PipelineService) GetPipelineStatus(pipelineID uuid.UUID) (string, error) {
	return ps.Repository.GetPipelineStatus(pipelineID.String())
}

func (ps *PipelineService) CancelPipeline(pipelineID uuid.UUID) error {
	return ps.updatePipelineStatus(pipelineID, "Canceled")
}

func (ps *PipelineService) updatePipelineStatus(pipelineID uuid.UUID, status string) error {
	return ps.Repository.UpdatePipelineExecution(&models.PipelineExecution{
		ID:     pipelineID,
		Status: status,
	})
}

func (ps *PipelineService) logExecutionError(pipelineID uuid.UUID, errorMsg string) {
	logErr := ps.Repository.SaveExecutionLog(&models.ExecutionLog{
		ID:         uuid.New(),
		PipelineID: pipelineID,
		Status:     "Error",
		ErrorMsg:   errorMsg,
	})
	if logErr != nil {
		log.Printf("Failed to log execution error: %v", logErr)
	}
}
