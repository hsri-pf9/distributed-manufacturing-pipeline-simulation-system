package domain

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/google/uuid"
)

type SequentialPipelineOrchestrator struct {
	ID     uuid.UUID
	Stages []Stage
	Status map[uuid.UUID]string 
	mu     sync.Mutex           
}

func NewSequentialPipelineOrchestrator() *SequentialPipelineOrchestrator {
	return &SequentialPipelineOrchestrator{
		ID:     uuid.New(),
		Stages: []Stage{},
		Status: make(map[uuid.UUID]string),
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

	// Mark pipeline as running
	p.updateStatus(pipelineID, "Running")

	var result interface{} = input
	var completedStages []Stage

	for _, stage := range p.Stages {
		log.Printf("Executing stage: %v\n", stage.GetID())

		var err error
		result, err = stage.Execute(ctx, result)
		if err != nil {
			log.Printf("Error in stage %v: %v. Rolling back...\n", stage.GetID(), err)
			stage.HandleError(ctx, err)
			p.rollback(ctx, completedStages, result)

			p.updateStatus(pipelineID, "Failed")
			return nil, err
		}

		completedStages = append(completedStages, stage)
	}

	p.updateStatus(pipelineID, "Completed")
	return result, nil
}

func (p *SequentialPipelineOrchestrator) rollback(ctx context.Context, completedStages []Stage, input interface{}) {
	for _, stage := range completedStages {
		stage.Rollback(ctx, input)
	}
}

func (p *SequentialPipelineOrchestrator) Cancel(pipelineID uuid.UUID) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.Status[pipelineID]; !exists {
		return errors.New("pipeline ID not found")
	}

	log.Printf("Cancelling pipeline: %v\n", pipelineID)
	p.Status[pipelineID] = "Cancelled"
	return nil
}

func (p *SequentialPipelineOrchestrator) GetStatus(pipelineID uuid.UUID) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	status, exists := p.Status[pipelineID]
	if !exists {
		return "", errors.New("pipeline ID not found")
	}
	return status, nil
}

func (p *SequentialPipelineOrchestrator) updateStatus(pipelineID uuid.UUID, status string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Status[pipelineID] = status
}
