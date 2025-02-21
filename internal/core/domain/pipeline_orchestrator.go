package domain

import (
	"context"

	"github.com/google/uuid"
)

// PipelineOrchestrator defines the core contract for executing pipelines
type PipelineOrchestrator interface {
	AddStage(stage Stage) error
	Execute(ctx context.Context, userID uuid.UUID, input interface{}) (uuid.UUID, interface{}, error)
	GetStatus(pipelineID uuid.UUID) (string, error)
	Cancel(pipelineID uuid.UUID, userID uuid.UUID) error
}
