package domain

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/models"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/ports"
)

type ParallelPipelineOrchestrator struct {
	PipelineID uuid.UUID
	UserID     uuid.UUID
	Stages     []Stage
	mu         sync.Mutex
	dbRepo ports.PipelineRepository
}

// NewParallelPipelineOrchestrator initializes a new parallel orchestrator
func NewParallelPipelineOrchestrator(dbRepo ports.PipelineRepository) *ParallelPipelineOrchestrator {
	return &ParallelPipelineOrchestrator{
		PipelineID: uuid.New(),
		dbRepo:     dbRepo,
	}
}

// AddStage adds a new stage to the parallel pipeline
func (p *ParallelPipelineOrchestrator) AddStage(stage Stage) error {
	if stage == nil {
		return errors.New("stage cannot be nil")
	}
	p.mu.Lock()
	p.Stages = append(p.Stages, stage)
	p.mu.Unlock()
	return nil
}

// Execute runs all stages concurrently and logs execution details in the database
func (p *ParallelPipelineOrchestrator) Execute(ctx context.Context, userID uuid.UUID, input interface{}) (uuid.UUID, interface{}, error) {
	// Step 1: Validate user existence
	user, err := p.dbRepo.GetUserByID(userID)
	if err != nil {
		log.Printf("Failed to validate user existence: %v", err)
		return p.PipelineID, nil, err
	}
	if user == nil {
		return p.PipelineID, nil, errors.New("user does not exist")
	}
	// Step 1: Create a new pipeline execution record in DB
	pipelineExecution := &models.PipelineExecution{
		PipelineID: p.PipelineID,
		UserID:     userID,
		Status:     "Running",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := p.dbRepo.SavePipelineExecution(pipelineExecution); err != nil {
		log.Printf("Failed to save pipeline execution: %v", err)
		return p.PipelineID, nil, err
	}

	// Step 2: Execute stages in parallel
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make([]interface{}, 0, len(p.Stages))
	errorsSlice := make([]error, 0, len(p.Stages))

	for _, stage := range p.Stages {
		wg.Add(1)
		go func(stage Stage) {
			defer wg.Done()
			result, err := stage.Execute(ctx, input)

			logEntry := &models.ExecutionLog{
				StageID:    stage.GetID(),
				PipelineID: p.PipelineID,
				Status:     "Completed",
				Timestamp:  time.Now(),
			}

			if err != nil {
				logEntry.Status = "Failed"
				logEntry.ErrorMsg = err.Error()
				mu.Lock()
				errorsSlice = append(errorsSlice, err)
				mu.Unlock()
			} else {
				mu.Lock()
				results = append(results, result)
				mu.Unlock()
			}

			// Save execution log for each stage
			if err := p.dbRepo.SaveExecutionLog(logEntry); err != nil {
				log.Printf("Failed to save execution log: %v", err)
			}
		}(stage)
	}

	wg.Wait()

	// Step 3: Update pipeline execution status in DB
	if len(errorsSlice) > 0 {
		pipelineExecution.Status = "Failed"
		if err := p.dbRepo.UpdatePipelineExecution(pipelineExecution); err != nil {
			log.Printf("Failed to update pipeline execution status: %v", err)
		}
		return p.PipelineID, nil, errors.New("pipeline execution failed")
	}

	pipelineExecution.Status = "Completed"
	if err := p.dbRepo.UpdatePipelineExecution(pipelineExecution); err != nil {
		log.Printf("Failed to update pipeline execution status: %v", err)
	}

	if len(results) == 0 {
		return p.PipelineID, nil, errors.New("no valid results from pipeline stages")
	}
	return p.PipelineID, results, nil
}

// GetStatus retrieves the status of a pipeline from the database
func (p *ParallelPipelineOrchestrator) GetStatus(pipelineID uuid.UUID) (string, error) {
	return p.dbRepo.GetPipelineStatus(pipelineID.String())
}

// Cancel updates the pipeline execution status to "Canceled"
func (p *ParallelPipelineOrchestrator) Cancel(pipelineID uuid.UUID, userID uuid.UUID) error {
	pipelineExecution := &models.PipelineExecution{
		PipelineID: pipelineID,
		UserID:     userID,
		Status:     "Canceled",
		UpdatedAt:  time.Now(),
	}
	return p.dbRepo.UpdatePipelineExecution(pipelineExecution)
}
