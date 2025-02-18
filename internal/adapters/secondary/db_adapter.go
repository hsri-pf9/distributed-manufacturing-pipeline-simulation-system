package secondary

import (
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/models"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/ports"
)

type DatabaseAdapter struct{}

var _ ports.PipelineRepository = (*DatabaseAdapter)(nil)

func NewDatabaseAdapter() *DatabaseAdapter {
	return &DatabaseAdapter{}
}

func (d *DatabaseAdapter) SavePipelineExecution(execution *models.PipelineExecution) error {
	return DB.Create(execution).Error
}

func (d *DatabaseAdapter) UpdatePipelineExecution(execution *models.PipelineExecution) error {
	return DB.Model(execution).Where("id = ?", execution.ID).Update("status", execution.Status).Error
}

func (d *DatabaseAdapter) SaveExecutionLog(logEntry *models.ExecutionLog) error {
	return DB.Create(logEntry).Error
}

func (d *DatabaseAdapter) GetPipelineStatus(pipelineID string) (string, error) {
	var execution models.PipelineExecution
	if err := DB.Where("pipeline_id = ?", pipelineID).First(&execution).Error; err != nil {
		return "", err
	}
	return execution.Status, nil
}