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
	ID     uuid.UUID
	Stages []Stage
	Status map[uuid.UUID]string 
	DBAdapter  ports.PipelineRepository           
}

func NewSequentialPipelineOrchestrator(dbAdapter ports.PipelineRepository) *SequentialPipelineOrchestrator {
	return &SequentialPipelineOrchestrator{
		ID:     uuid.New(),
		Stages: []Stage{},
		Status: make(map[uuid.UUID]string),
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

func (p *SequentialPipelineOrchestrator) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	pipelineID := p.ID

	pipelineExecution := &models.PipelineExecution{
		ID:         pipelineID,
		Status:     "Running",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := p.DBAdapter.SavePipelineExecution(pipelineExecution); err != nil {
		return nil, err
	}

	var result interface{} = input
	var completedStages []Stage

	for _, stage := range p.Stages {
		log.Printf("Executing stage: %v\n", stage.GetID())
		var err error
		result, err = stage.Execute(ctx, result)
		logEntry := &models.ExecutionLog{
			ID:         uuid.New(),
			StageID:    stage.GetID(),
			PipelineID: pipelineID,
			Status:     "Completed",
			Timestamp:  time.Now(),
		}
		if err != nil {
			logEntry.Status = "Failed"
			logEntry.ErrorMsg = err.Error()
			p.DBAdapter.SaveExecutionLog(logEntry)
			p.rollback(ctx, completedStages, result)
			p.DBAdapter.UpdatePipelineExecution(&models.PipelineExecution{ID: pipelineID, Status: "Failed"})
			return nil, err
		}
		p.DBAdapter.SaveExecutionLog(logEntry)
		completedStages = append(completedStages, stage)
	}

	p.DBAdapter.UpdatePipelineExecution(&models.PipelineExecution{ID: pipelineID, Status: "Completed"})
	return result, nil
}

func (p *SequentialPipelineOrchestrator) rollback(ctx context.Context, completedStages []Stage, input interface{}) {
	for _, stage := range completedStages {
		stage.Rollback(ctx, input)
	}
}

func (p *SequentialPipelineOrchestrator) Cancel(pipelineID uuid.UUID) error {
	return p.DBAdapter.UpdatePipelineExecution(&models.PipelineExecution{ID: pipelineID, Status: "Cancelled"})
}

func (p *SequentialPipelineOrchestrator) GetStatus(pipelineID uuid.UUID) (string, error) {
	return p.DBAdapter.GetPipelineStatus(pipelineID.String())
}
