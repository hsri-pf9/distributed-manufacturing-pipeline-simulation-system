package ports

import ( "github.com/google/uuid"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/models"
)

type PipelineRepository interface {
	SavePipelineExecution(execution *models.PipelineExecution) error
	UpdatePipelineExecution(execution *models.PipelineExecution) error
	SaveExecutionLog(logEntry *models.ExecutionLog) error
	GetPipelineStatus(pipelineID string) (string, error)

	GetUserByID(userID uuid.UUID) (*models.User, error)
	SaveUser(user *models.User) error
	UpdateUser(userID uuid.UUID, updates map[string]interface{}) error 
	GetPipelinesByUser(userID string) ([]models.PipelineExecution, error)
}
