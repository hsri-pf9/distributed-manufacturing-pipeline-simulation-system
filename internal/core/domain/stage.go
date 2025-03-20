package domain

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/utils"
)

type Stage interface {
	GetID() uuid.UUID
	// Execute(ctx context.Context, input interface{}) (interface{}, error)
	Execute(ctx context.Context, input interface{}, sse *utils.SSEManager, pipelineID uuid.UUID) (interface{}, error)
	HandleError(ctx context.Context, err error) error
	Rollback(ctx context.Context, input interface{}) error
}

type BaseStage struct {
	ID uuid.UUID
}

func NewBaseStage() *BaseStage {
	return &BaseStage{ID: uuid.New()}
}

func (s *BaseStage) GetID() uuid.UUID {
	return s.ID
}

func (s *BaseStage) Execute(ctx context.Context, input interface{}, sse *utils.SSEManager, pipelineID uuid.UUID) (interface{}, error) {
	log.Printf("Executing stage: %s with input: %v\n", s.ID, input)

	// ✅ Broadcast stage execution start as JSON
	sse.BroadcastUpdate(map[string]interface{}{
		"type":        "stage",
		"stage_id":    s.ID.String(),
		"pipeline_id": pipelineID.String(),
		"status":      "Running",
	})

	if input == nil {
		err := errors.New("input is nil, stage execution failed")
		log.Printf("Stage %s execution failed: %v", s.ID, err)

		// ✅ Broadcast stage failure as JSON
		sse.BroadcastUpdate(map[string]interface{}{
			"type":        "stage",
			"stage_id":    s.ID.String(),
			"pipeline_id": pipelineID.String(),
			"status":      "Failed",
		})

		return nil, err
	}

	time.Sleep(5 * time.Second)

	log.Printf("Stage %s executed successfully", s.ID)

	// ✅ Broadcast stage completion as JSON
	sse.BroadcastUpdate(map[string]interface{}{
		"type":        "stage",
		"stage_id":    s.ID.String(),
		"pipeline_id": pipelineID.String(),
		"status":      "Completed",
	})

	return input, nil
}

func (s *BaseStage) HandleError(ctx context.Context, err error) error {
	log.Printf("Error in stage %s execution: %v", s.ID, err)
	return errors.New("stage execution failed: " + err.Error())
}

func (s *BaseStage) Rollback(ctx context.Context, input interface{}) error {
	log.Printf("Rolling back stage %s due to failure. Input: %v", s.ID, input)
	return nil
}