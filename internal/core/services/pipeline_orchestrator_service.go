package services

import (
	"context"
	"log"
	"time"

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

func (ps *PipelineService) CreatePipeline(userID uuid.UUID, stageCount int, isParallel bool) (uuid.UUID, error) {
	pipelineID := uuid.New()

	var orchestrator domain.PipelineOrchestrator
	if isParallel {
		orchestrator = ps.ParallelOrchestrator
	} else {
		orchestrator = ps.SequentialOrchestrator
	}

	
	for i := 0; i < stageCount; i++ {
		stage := domain.NewBaseStage()
		if err := orchestrator.AddStage(stage); err != nil {
			return uuid.Nil, err
		}
	}

	err := ps.Repository.SavePipelineExecution(&models.PipelineExecution{
		PipelineID: pipelineID,
		UserID:     userID,
		Status:     "Created",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
	if err != nil {
		return uuid.Nil, err
	}

	return pipelineID, nil
}

func (ps *PipelineService) StartPipeline(ctx context.Context, userID uuid.UUID, pipelineID uuid.UUID, input interface{}, isParallel bool) error {
	var orchestrator domain.PipelineOrchestrator
	if isParallel {
		orchestrator = ps.ParallelOrchestrator
	} else {
		orchestrator = ps.SequentialOrchestrator
	}

	if err := ps.updatePipelineStatus(pipelineID, "Running"); err != nil {
		return err
	}

	stageID, _, err := orchestrator.Execute(ctx, userID, input)
	if err != nil {
		_ = ps.updatePipelineStatus(pipelineID, "Failed")
		ps.logExecutionError(pipelineID, stageID, err.Error())
		return err
	}

	return ps.updatePipelineStatus(pipelineID, "Completed")
}

// func (ps *PipelineService) GetPipelineStatus(pipelineID uuid.UUID) (string, error) {
// 	return ps.Repository.GetPipelineStatus(pipelineID.String())
// }

func (ps *PipelineService) GetPipelineStatus(pipelineID uuid.UUID, isParallel bool) (string, error) {
	var orchestrator domain.PipelineOrchestrator
	if isParallel {
		orchestrator = ps.ParallelOrchestrator
	} else {
		orchestrator = ps.SequentialOrchestrator
	}

	return orchestrator.GetStatus(pipelineID)
}

// func (ps *PipelineService) CancelPipeline(pipelineID uuid.UUID) error {
// 	return ps.updatePipelineStatus(pipelineID, "Canceled")
// }
func (ps *PipelineService) CancelPipeline(pipelineID uuid.UUID, userID uuid.UUID, isParallel bool) error {
	var orchestrator domain.PipelineOrchestrator
	if isParallel {
		orchestrator = ps.ParallelOrchestrator
	} else {
		orchestrator = ps.SequentialOrchestrator
	}
	log.Printf("Cancelling pipeline: %s by user: %s", pipelineID, userID)

	// return orchestrator.Cancel(pipelineID, userID)
	err := orchestrator.Cancel(pipelineID, userID)
    if err != nil {
        log.Printf("Failed to cancel pipeline: %v", err)
    }

    return err
}

func (ps *PipelineService) updatePipelineStatus(pipelineID uuid.UUID, status string) error {
	return ps.Repository.UpdatePipelineExecution(&models.PipelineExecution{
		PipelineID: pipelineID,
		Status:     status,
		UpdatedAt:  time.Now(),
	})
}

func (ps *PipelineService) logExecutionError(pipelineID uuid.UUID, stageID uuid.UUID, errorMsg string) {
	logErr := ps.Repository.SaveExecutionLog(&models.ExecutionLog{
		StageID:    stageID,
		PipelineID: pipelineID,
		Status:     "Error",
		ErrorMsg:   errorMsg,
		Timestamp:  time.Now(),
	})
	if logErr != nil {
		log.Printf("Failed to log execution error: %v", logErr)
	}
}
