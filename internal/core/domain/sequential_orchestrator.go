package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"log"
)

// SequentialPipelineOrchestrator executes stages sequentially
type SequentialPipelineOrchestrator struct {
	ID     uuid.UUID
	Stages []Stage
}

// NewSequentialPipelineOrchestrator initializes a new instance
func NewSequentialPipelineOrchestrator() *SequentialPipelineOrchestrator {
	return &SequentialPipelineOrchestrator{
		ID:     uuid.New(),
		Stages: []Stage{},
	}
}

// AddStage adds a stage to the sequential pipeline
func (p *SequentialPipelineOrchestrator) AddStage(stage Stage) error {
	if stage == nil {
		return errors.New("stage cannot be nil")
	}
	p.Stages = append(p.Stages, stage)
	return nil
}

// Execute runs all stages sequentially, rolling back in case of errors
func (p *SequentialPipelineOrchestrator) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	var result interface{} = input
	completedStages := []Stage{}

	for _, stage := range p.Stages {
		log.Println("Executing stage:", stage.GetID())
		var err error
		result, err = stage.Execute(ctx, result)
		if err != nil {
			stage.HandleError(ctx, err)
			log.Println("Rolling back completed stages...")
			p.rollback(ctx, completedStages, result)
			return nil, err
		}
		completedStages = append(completedStages, stage)
	}
	return result, nil
}

// rollback reverts completed stages upon failure
func (p *SequentialPipelineOrchestrator) rollback(ctx context.Context, completedStages []Stage, input interface{}) {
	for _, stage := range completedStages {
		stage.Rollback(ctx, input)
	}
}

func (p *SequentialPipelineOrchestrator) Cancel(pipelineID uuid.UUID) error {
	log.Println("Cancel operation is not fully implemented for SequentialPipelineOrchestrator")
	return nil
}

// GetStatus returns the status of the sequential pipeline execution
func (p *SequentialPipelineOrchestrator) GetStatus(pipelineID uuid.UUID) (string, error) {
	// Currently, sequential execution does not maintain a pipeline status map
	// You can extend this later with actual tracking.
	log.Println("GetStatus is not fully implemented for SequentialPipelineOrchestrator")
	return "Unknown", nil
}