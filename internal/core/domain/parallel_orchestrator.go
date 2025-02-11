package domain

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
)

type ParallelPipelineOrchestrator struct {
	ID     uuid.UUID
	Stages []Stage
	mu     sync.Mutex
	status map[uuid.UUID]string
}

// NewParallelPipelineOrchestrator initializes a new parallel orchestrator
func NewParallelPipelineOrchestrator() *ParallelPipelineOrchestrator {
	return &ParallelPipelineOrchestrator{
		ID:     uuid.New(),
		status: make(map[uuid.UUID]string),
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

// Execute runs all stages concurrently
func (p *ParallelPipelineOrchestrator) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	pipelineID := uuid.New()
	p.UpdateStatus(pipelineID, "Running")

	var wg sync.WaitGroup
	results := make(chan interface{}, len(p.Stages))
	errorsChan := make(chan error, len(p.Stages))

	for _, stage := range p.Stages {
		wg.Add(1)
		go func(stage Stage) {
			defer wg.Done()
			result, err := stage.Execute(ctx, input)
			if err != nil {
				errorsChan <- err
			} else {
				results <- result
			}
		}(stage)
	}

	wg.Wait()
	close(results)
	close(errorsChan)

	if len(errorsChan) > 0 {
		p.UpdateStatus(pipelineID, "Failed")
		return nil, errors.New("pipeline execution failed")
	}

	p.UpdateStatus(pipelineID, "Completed")
	return <-results, nil
}

// GetStatus retrieves the status of a pipeline
func (p *ParallelPipelineOrchestrator) GetStatus(pipelineID uuid.UUID) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	status, exists := p.status[pipelineID]
	if !exists {
		return "", errors.New("pipeline not found")
	}
	return status, nil
}

// Cancel updates the status to "Canceled"
func (p *ParallelPipelineOrchestrator) Cancel(pipelineID uuid.UUID) error {
	p.UpdateStatus(pipelineID, "Canceled")
	return nil
}

// UpdateStatus modifies the pipeline status
func (p *ParallelPipelineOrchestrator) UpdateStatus(pipelineID uuid.UUID, status string) {
	p.mu.Lock()
	p.status[pipelineID] = status
	p.mu.Unlock()
}

