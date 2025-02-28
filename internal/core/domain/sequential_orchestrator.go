package domain

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/models"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/ports"
)

type SequentialPipelineOrchestrator struct {
	ID        uuid.UUID
	Stages    []Stage
	Status    map[uuid.UUID]string
	DBAdapter ports.PipelineRepository
}

func NewSequentialPipelineOrchestrator(pipelineID uuid.UUID, dbAdapter ports.PipelineRepository) *SequentialPipelineOrchestrator {
	if pipelineID == uuid.Nil {
		log.Fatal("ERROR: Pipeline ID cannot be Nil")
	}
	return &SequentialPipelineOrchestrator{
		// ID:        uuid.New(),
		ID:        pipelineID,
		Stages:    []Stage{},
		Status:    make(map[uuid.UUID]string),
		DBAdapter: dbAdapter,
	}
}

func (p *SequentialPipelineOrchestrator) AddStage(stage Stage) error {
	if stage == nil {
		return errors.New("stage cannot be nil")
	}
	p.Stages = append(p.Stages, stage)
	return nil
}
 
func (p *SequentialPipelineOrchestrator) Execute(ctx context.Context, userID uuid.UUID, pipelineID uuid.UUID, input interface{}) (uuid.UUID, interface{}, error) {
	// 🚀 **Fix: Ensure stages exist before execution**
	// if len(p.Stages) == 0 {
	// 	log.Println("Error: No stages found for execution")
	// 	return uuid.Nil, nil, errors.New("pipeline has no stages to execute")
	// }

	// Ensure the user exists before proceeding
	user, err := p.DBAdapter.GetUserByID(userID)
	if err != nil {
		return uuid.Nil, nil, errors.New("user not found")
	}

	// pipelineExecution := &models.PipelineExecution{
	// 	// PipelineID: p.ID,
	// 	PipelineID: pipelineID,
	// 	UserID:     user.UserID,
	// 	Status:     "Running",
	// 	CreatedAt:  time.Now(),
	// 	UpdatedAt:  time.Now(),
	// }

	// if err := p.DBAdapter.SavePipelineExecution(pipelineExecution); err != nil {
	// 	return uuid.Nil, nil, err
	// }

	err = p.DBAdapter.UpdatePipelineExecution(&models.PipelineExecution{
        PipelineID: pipelineID,
		UserID:     user.UserID,
        Status:     "Running",
        UpdatedAt:  time.Now(),
    })
    if err != nil {
        return uuid.Nil, nil, err
    }

	if len(p.Stages) == 0 {
		log.Println("Error: No stages found for execution")
		return uuid.Nil, nil, errors.New("pipeline has no stages to execute")
	}

	var result interface{} = input
	var completedStages []Stage

	for _, stage := range p.Stages {
		log.Printf("Executing stage: %v\n", stage.GetID())
		var err error
		result, err = stage.Execute(ctx, result)
		logEntry := &models.ExecutionLog{
			StageID:    stage.GetID(),
			// PipelineID: p.ID,
			PipelineID: pipelineID,
			Status:     "Completed",
			Timestamp:  time.Now(),
		}
		if err != nil {
			logEntry.Status = "Failed"
			logEntry.ErrorMsg = err.Error()
			p.DBAdapter.SaveExecutionLog(logEntry)
			p.rollback(ctx, completedStages, result)
			// Correctly update status with the right ID
			updateErr := p.DBAdapter.UpdatePipelineExecution(&models.PipelineExecution{
				// PipelineID: p.ID,
				PipelineID: pipelineID,
				Status:     "Failed",
				UpdatedAt:  time.Now(),
			})
			if updateErr != nil {
				log.Printf("Failed to update pipeline status: %v", updateErr)
			}

			return stage.GetID(), nil, err
		}
		p.DBAdapter.SaveExecutionLog(logEntry)
		completedStages = append(completedStages, stage)
	}

	// Step 5: Update Pipeline Execution to Completed
	// updateErr := p.DBAdapter.UpdatePipelineExecution(&models.PipelineExecution{
	// 	// PipelineID: p.ID,
	// 	PipelineID: pipelineID,
	// 	Status:     "Completed",
	// 	UpdatedAt:  time.Now(),
	// })
	// if updateErr != nil {
	// 	log.Printf("Failed to update pipeline status: %v", updateErr)
	// }

	err = p.DBAdapter.UpdatePipelineExecution(&models.PipelineExecution{
        PipelineID: pipelineID,
        Status:     "Completed",
        UpdatedAt:  time.Now(),
    })
    if err != nil {
        log.Printf("Failed to update pipeline status: %v", err)
    }

	return uuid.Nil, result, nil
}

func (p *SequentialPipelineOrchestrator) rollback(ctx context.Context, completedStages []Stage, input interface{}) {
	for _, stage := range completedStages {
		stage.Rollback(ctx, input)
	}
}

func (p *SequentialPipelineOrchestrator) Cancel(pipelineID uuid.UUID, userID uuid.UUID) error {
	log.Printf("Checking if pipeline %s exists before cancelling", pipelineID)
	// Step 1: Validate Pipeline Existence
	status, err := p.DBAdapter.GetPipelineStatus(pipelineID.String())
	if err != nil {
		log.Printf("Error fetching pipeline status: %v", err)
		return errors.New("pipeline not found")
	}

	log.Printf("Pipeline status: %s", status)

	// Step 2: Prevent Canceling Completed Pipelines
	if status == "Completed" {
		log.Printf("Pipeline %s is already completed, cannot cancel", pipelineID)
		return errors.New("cannot cancel a completed pipeline")
	}

	// Step 3: Update Status to Cancelled
	log.Printf("Cancelling pipeline %s...", pipelineID)
	// return p.DBAdapter.UpdatePipelineExecution(&models.PipelineExecution{
	// 	PipelineID: pipelineID,
	// 	Status:     "Cancelled",
	// 	UpdatedAt:  time.Now(),
	// })
	err = p.DBAdapter.UpdatePipelineExecution(&models.PipelineExecution{
		PipelineID: pipelineID,
		Status:     "Cancelled",
		UpdatedAt:  time.Now(),
	})

	if err != nil {
		log.Printf("Failed to update pipeline status: %v", err)
		return errors.New("failed to update pipeline status")
	}

	log.Printf("Pipeline %s successfully cancelled", pipelineID)
	return nil
}

func (p *SequentialPipelineOrchestrator) GetStatus(pipelineID uuid.UUID) (string, error) {
	status, err := p.DBAdapter.GetPipelineStatus(pipelineID.String())
	if err != nil {
		return "", errors.New("failed to retrieve pipeline status")
	}
	return status, nil
}
