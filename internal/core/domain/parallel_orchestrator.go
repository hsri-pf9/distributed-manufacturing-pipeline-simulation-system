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
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/utils"
)

type ParallelPipelineOrchestrator struct {
	PipelineID uuid.UUID
	UserID     uuid.UUID
	Stages     []Stage
	mu         sync.Mutex
	dbRepo ports.PipelineRepository
	SSE        *utils.SSEManager // ✅ Add SSEManager
}

// NewParallelPipelineOrchestrator initializes a new parallel orchestrator
// func NewParallelPipelineOrchestrator(pipelineID uuid.UUID, dbRepo ports.PipelineRepository) *ParallelPipelineOrchestrator {
func NewParallelPipelineOrchestrator(pipelineID uuid.UUID, dbRepo ports.PipelineRepository, sse *utils.SSEManager) *ParallelPipelineOrchestrator {
	return &ParallelPipelineOrchestrator{
		// PipelineID: uuid.New(),
		PipelineID: pipelineID,
		dbRepo:     dbRepo,
		Stages:     []Stage{},
		SSE:        sse,
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

// func (p *ParallelPipelineOrchestrator) Execute(ctx context.Context, userID uuid.UUID, pipelineID uuid.UUID, input interface{}) (uuid.UUID, interface{}, error) {
// 	// ✅ Step 1: Validate user existence
// 	user, err := p.dbRepo.GetUserByID(userID)
// 	if err != nil {
// 		log.Printf("Failed to validate user existence: %v", err)
// 		return pipelineID, nil, err
// 	}
// 	if user == nil {
// 		return pipelineID, nil, errors.New("user does not exist")
// 	}

// 	// ✅ Step 2: Update pipeline execution status instead of inserting a new record
// 	if err := p.dbRepo.UpdatePipelineExecution(&models.PipelineExecution{
// 		PipelineID: pipelineID,
// 		Status:     "Running",
// 		UpdatedAt:  time.Now(),
// 	}); err != nil {
// 		log.Printf("Failed to update pipeline execution status: %v", err)
// 		return pipelineID, nil, err
// 	}

// 	// ✅ Step 3: Execute stages in parallel
// 	var wg sync.WaitGroup
// 	var mu sync.Mutex
// 	results := make([]interface{}, 0, len(p.Stages))
// 	errorsSlice := make([]error, 0, len(p.Stages))

// 	for _, stage := range p.Stages {
// 		wg.Add(1)
// 		go func(stage Stage) {
// 			defer wg.Done()
// 			result, err := stage.Execute(ctx, input)

// 			logEntry := &models.ExecutionLog{
// 				StageID:    stage.GetID(),
// 				PipelineID: pipelineID,
// 				Status:     "Completed",
// 				Timestamp:  time.Now(),
// 			}

// 			if err != nil {
// 				logEntry.Status = "Failed"
// 				logEntry.ErrorMsg = err.Error()
// 				mu.Lock()
// 				errorsSlice = append(errorsSlice, err)
// 				mu.Unlock()
// 			} else {
// 				mu.Lock()
// 				results = append(results, result)
// 				mu.Unlock()
// 			}

// 			// Save execution log for each stage
// 			if err := p.dbRepo.SaveExecutionLog(logEntry); err != nil {
// 				log.Printf("Failed to save execution log: %v", err)
// 			}
// 		}(stage)
// 	}

// 	wg.Wait()

// 	// ✅ Step 4: Update pipeline execution status in DB based on results
// 	finalStatus := "Completed"
// 	if len(errorsSlice) > 0 {
// 		finalStatus = "Failed"
// 	}

// 	if err := p.dbRepo.UpdatePipelineExecution(&models.PipelineExecution{
// 		PipelineID: pipelineID,
// 		Status:     finalStatus,
// 		UpdatedAt:  time.Now(),
// 	}); err != nil {
// 		log.Printf("Failed to update final pipeline execution status: %v", err)
// 	}

// 	if len(results) == 0 {
// 		return pipelineID, nil, errors.New("no valid results from pipeline stages")
// 	}

// 	return pipelineID, results, nil
// }

func (p *ParallelPipelineOrchestrator) Execute(ctx context.Context, userID uuid.UUID, pipelineID uuid.UUID, input interface{}) (uuid.UUID, interface{}, error) {
	// ✅ Step 1: Validate user existence
	user, err := p.dbRepo.GetUserByID(userID)
	if err != nil {
		log.Printf("Failed to validate user existence: %v", err)
		return pipelineID, nil, err
	}
	if user == nil {
		return pipelineID, nil, errors.New("user does not exist")
	}

	// 🔹 Broadcast pipeline start event via SSE
	p.SSE.BroadcastUpdate(map[string]interface{}{
		"type":        "pipeline",
		"pipeline_id": pipelineID.String(),
		"status":      "Running",
	})

	// ✅ Step 2: Update pipeline execution status
	if err := p.dbRepo.UpdatePipelineExecution(&models.PipelineExecution{
		PipelineID: pipelineID,
		Status:     "Running",
		UpdatedAt:  time.Now(),
	}); err != nil {
		log.Printf("Failed to update pipeline execution status: %v", err)
		return pipelineID, nil, err
	}

	// ✅ Step 3: Execute stages in parallel
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make([]interface{}, 0, len(p.Stages))
	errorsSlice := make([]error, 0, len(p.Stages))

	for _, stage := range p.Stages {
		wg.Add(1)
		go func(stage Stage) {
			defer wg.Done()

			// 🔹 Broadcast stage start event via SSE
			p.SSE.BroadcastUpdate(map[string]interface{}{
				"type":        "stage",
				"stage_id":    stage.GetID().String(),
				"pipeline_id": pipelineID.String(),
				"status":      "Running",
			})

			result, err := stage.Execute(ctx, input, p.SSE, pipelineID)

			logEntry := &models.ExecutionLog{
				StageID:    stage.GetID(),
				PipelineID: pipelineID,
				Status:     "Completed",
				Timestamp:  time.Now(),
			}

			if err != nil {
				logEntry.Status = "Failed"
				logEntry.ErrorMsg = err.Error()
				mu.Lock()
				errorsSlice = append(errorsSlice, err)
				mu.Unlock()

				// 🔹 Broadcast stage failure via SSE
				p.SSE.BroadcastUpdate(map[string]interface{}{
					"type":        "stage",
					"stage_id":    stage.GetID().String(),
					"pipeline_id": pipelineID.String(),
					"status":      "Failed",
				})

			} else {
				mu.Lock()
				results = append(results, result)
				mu.Unlock()

				// 🔹 Broadcast stage completion via SSE
				p.SSE.BroadcastUpdate(map[string]interface{}{
					"type":        "stage",
					"stage_id":    stage.GetID().String(),
					"pipeline_id": pipelineID.String(),
					"status":      "Completed",
				})
			}

			// Save execution log for each stage
			if err := p.dbRepo.SaveExecutionLog(logEntry); err != nil {
				log.Printf("Failed to save execution log: %v", err)
			}
		}(stage)
	}

	wg.Wait()

	// ✅ Step 4: Update pipeline execution status
	finalStatus := "Completed"
	if len(errorsSlice) > 0 {
		finalStatus = "Failed"
	}

	if err := p.dbRepo.UpdatePipelineExecution(&models.PipelineExecution{
		PipelineID: pipelineID,
		Status:     finalStatus,
		UpdatedAt:  time.Now(),
	}); err != nil {
		log.Printf("Failed to update final pipeline execution status: %v", err)
	}

	// 🔹 Broadcast pipeline completion via SSE
	p.SSE.BroadcastUpdate(map[string]interface{}{
		"type":        "pipeline",
		"pipeline_id": pipelineID.String(),
		"status":      finalStatus,
	})

	if len(results) == 0 {
		return pipelineID, nil, errors.New("no valid results from pipeline stages")
	}

	return pipelineID, results, nil
}

// GetStatus retrieves the status of a pipeline from the database
func (p *ParallelPipelineOrchestrator) GetStatus(pipelineID uuid.UUID) (string, error) {
	return p.dbRepo.GetPipelineStatus(pipelineID.String())
}

// Cancel updates the pipeline execution status to "Canceled"
func (p *ParallelPipelineOrchestrator) Cancel(pipelineID uuid.UUID, userID uuid.UUID) error {
	log.Printf("Cancelling pipeline: %s for user: %s", pipelineID, userID)

	status, err := p.dbRepo.GetPipelineStatus(pipelineID.String())
	if err != nil {
		log.Printf("Error fetching pipeline status: %v", err)
		return errors.New("pipeline not found")
	}

	if status == "Completed" {
		log.Printf("Pipeline %s is already completed, cannot cancel", pipelineID)
		return errors.New("cannot cancel a completed pipeline")
	}
	log.Printf("Cancelling pipeline %s...", pipelineID)

	err = p.dbRepo.UpdatePipelineExecution(&models.PipelineExecution{
		PipelineID: pipelineID, // ✅ Ensure correct pipeline ID
		Status:     "Cancelled",
		UpdatedAt:  time.Now(),
	})

	if err != nil {
		log.Printf("Failed to update pipeline status: %v", err)
		return errors.New("failed to update pipeline status")
	}

	log.Printf("Pipeline %s successfully cancelled", pipelineID)

	// 🔹 Broadcast cancellation event via SSE
	p.SSE.BroadcastUpdate(map[string]interface{}{
		"type":        "pipeline",
		"pipeline_id": pipelineID.String(),
		"status":      "Cancelled",
	})
	return nil
}
