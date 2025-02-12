package ports

import "github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/models"

type PipelineRepository interface {
	SavePipelineExecution(execution *models.PipelineExecution) error
	UpdatePipelineExecution(execution *models.PipelineExecution) error
	SaveExecutionLog(logEntry *models.ExecutionLog) error
	GetPipelineStatus(pipelineID string) (string, error)
}
