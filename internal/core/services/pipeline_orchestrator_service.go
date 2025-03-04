package services

import (
	"context"
	"log"
	"time"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/domain"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/ports"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/models"
)

type PipelineService struct {
	SequentialOrchestrators map[uuid.UUID]*domain.SequentialPipelineOrchestrator
	ParallelOrchestrators   map[uuid.UUID]*domain.ParallelPipelineOrchestrator
	Repository             ports.PipelineRepository
	mu                     sync.RWMutex
}

func NewPipelineService(repo ports.PipelineRepository) *PipelineService {
	return &PipelineService{
		SequentialOrchestrators: make(map[uuid.UUID]*domain.SequentialPipelineOrchestrator),
		ParallelOrchestrators:   make(map[uuid.UUID]*domain.ParallelPipelineOrchestrator),
		Repository:             repo,
	}
}

func (ps *PipelineService) CreatePipeline(userID uuid.UUID, stageCount int, isParallel bool) (uuid.UUID, error) {
	pipelineID := uuid.New()

	ps.mu.Lock()
	var orchestrator domain.PipelineOrchestrator
	if isParallel {
		orchestrator = domain.NewParallelPipelineOrchestrator(pipelineID, ps.Repository)
		ps.ParallelOrchestrators[pipelineID] = orchestrator.(*domain.ParallelPipelineOrchestrator)
	} else {
		orchestrator = domain.NewSequentialPipelineOrchestrator(pipelineID, ps.Repository)
		ps.SequentialOrchestrators[pipelineID] = orchestrator.(*domain.SequentialPipelineOrchestrator)
	}
	ps.mu.Unlock()

	
	for i := 0; i < stageCount; i++ {
		stage := domain.NewBaseStage()
		log.Printf("Adding Stage: %s to Pipeline: %s", stage.GetID(), pipelineID) // Debugging log
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

// ✅ Start pipeline execution based on pipeline ID
func (ps *PipelineService) StartPipeline(ctx context.Context, userID uuid.UUID, pipelineID uuid.UUID, input interface{}, isParallel bool) error {
	ps.mu.RLock()
	var orchestrator domain.PipelineOrchestrator
	if isParallel {
		orchestrator = ps.ParallelOrchestrators[pipelineID]
	} else {
		orchestrator = ps.SequentialOrchestrators[pipelineID]
	}
	ps.mu.RUnlock()

	if orchestrator == nil {
		return errors.New("orchestrator not initialized for this pipeline")
	}

	status, err := ps.Repository.GetPipelineStatus(pipelineID.String())
	if err != nil {
		log.Printf("Failed to get pipeline status: %v", err)
		return err
	}
	if status != "Created" && status != "Paused" {
		return errors.New("invalid pipeline status: " + status)
	}
	// 🚀 **Fix: Ensure stages are present before execution**
	switch o := orchestrator.(type) {
	case *domain.SequentialPipelineOrchestrator:
		if len(o.Stages) == 0 {
			return errors.New("no stages found for this pipeline execution")
		}
	case *domain.ParallelPipelineOrchestrator:
		if len(o.Stages) == 0 {
			return errors.New("no stages found for this pipeline execution")
		}
	default:
		return errors.New("unknown orchestrator type")
	}

	if err := ps.updatePipelineStatus(pipelineID, "Running"); err != nil {
		return err
	}

	stageID, _, err := orchestrator.Execute(ctx, userID, pipelineID, input)
	if err != nil {
		_ = ps.updatePipelineStatus(pipelineID, "Failed")
		ps.logExecutionError(pipelineID, stageID, err.Error())
		return err
	}

	return ps.updatePipelineStatus(pipelineID, "Completed")
}

// ✅ Retrieve pipeline status
func (ps *PipelineService) GetPipelineStatus(pipelineID uuid.UUID, isParallel bool) (string, error) {
	ps.mu.RLock()
	var orchestrator domain.PipelineOrchestrator
	if isParallel {
		orchestrator = ps.ParallelOrchestrators[pipelineID]
	} else {
		orchestrator = ps.SequentialOrchestrators[pipelineID]
	}
	ps.mu.RUnlock()

	if orchestrator == nil {
		return "", errors.New("orchestrator not found for pipeline")
	}

	return orchestrator.GetStatus(pipelineID)
}

// ✅ Cancel pipeline execution
func (ps *PipelineService) CancelPipeline(pipelineID uuid.UUID, userID uuid.UUID, isParallel bool) error {
	ps.mu.RLock()
	var orchestrator domain.PipelineOrchestrator
	if isParallel {
		orchestrator = ps.ParallelOrchestrators[pipelineID]
	} else {
		orchestrator = ps.SequentialOrchestrators[pipelineID]
	}
	ps.mu.RUnlock()

	if orchestrator == nil {
		log.Printf("Orchestrator not found for pipeline: %s", pipelineID)
		return errors.New("orchestrator not initialized for this pipeline")
	}

	log.Printf("Cancelling pipeline: %s by user: %s", pipelineID, userID)

	err := orchestrator.Cancel(pipelineID, userID)
	if err != nil {
		log.Printf("Failed to cancel pipeline: %v", err)
		_ = ps.updatePipelineStatus(pipelineID, "Failed to Cancel")
		return err
	}

	return ps.updatePipelineStatus(pipelineID, "Cancelled")
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

func (ps *PipelineService) GetPipelinesByUser(userID string) ([]models.PipelineExecution, error) {
	return ps.Repository.GetPipelinesByUser(userID)
}
