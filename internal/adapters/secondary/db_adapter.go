// package secondary

// import (
// 	"github.com/google/uuid"
// 	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/models"
// 	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/ports"
// )

// type DatabaseAdapter struct{}

// var _ ports.PipelineRepository = (*DatabaseAdapter)(nil)

// func NewDatabaseAdapter() *DatabaseAdapter {
// 	return &DatabaseAdapter{}
// }

// // SaveUser inserts a new user into the database
// func (d *DatabaseAdapter) SaveUser(user *models.User) error {
// 	return DB.Create(user).Error
// }

// // GetUserByID retrieves a user by their ID
// func (d *DatabaseAdapter) GetUserByID(userID uuid.UUID) (*models.User, error) {
// 	var user models.User
// 	if err := DB.First(&user, "user_id = ?", userID).Error; err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }

// // UpdateUser updates user details
// func (d *DatabaseAdapter) UpdateUser(userID uuid.UUID, updates map[string]interface{}) error {
// 	return DB.Model(&models.User{}).Where("user_id = ?", userID).Updates(updates).Error
// }

// func (d *DatabaseAdapter) SavePipelineExecution(execution *models.PipelineExecution) error {
// 	return DB.Create(execution).Error
// }

// func (d *DatabaseAdapter) UpdatePipelineExecution(execution *models.PipelineExecution) error {
// 	return DB.Model(execution).Where("pipeline_id = ?", execution.PipelineID).Update("status", execution.Status).Error
// }

// func (d *DatabaseAdapter) SaveExecutionLog(logEntry *models.ExecutionLog) error {
// 	return DB.Create(logEntry).Error
// }

// func (d *DatabaseAdapter) GetPipelineStatus(pipelineID string) (string, error) {
// 	var execution models.PipelineExecution
// 	if err := DB.Where("pipeline_id = ?", pipelineID).First(&execution).Error; err != nil {
// 		return "", err
// 	}
// 	return execution.Status, nil
// }

// // GetPipelinesByUser retrieves all pipelines for a specific user
// func (d *DatabaseAdapter) GetPipelinesByUser(userID string) ([]models.PipelineExecution, error) {
// 	var pipelines []models.PipelineExecution
// 	err := DB.Where("user_id = ?", userID).Find(&pipelines).Error
// 	return pipelines, err
// }
// // GetPipelineStages fetches all stages associated with a pipeline
// func (d *DatabaseAdapter) GetPipelineStages(pipelineID uuid.UUID) ([]models.ExecutionLog, error) {
// 	var stages []models.ExecutionLog
// 	if err := DB.Where("pipeline_id = ?", pipelineID).Find(&stages).Error; err != nil {
// 		return nil, err
// 	}
// 	return stages, nil
// }

package secondary

import (
	"log"

	"github.com/google/uuid"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/models"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/ports"
	"gorm.io/gorm"
)

type DatabaseAdapter struct {
	DB *gorm.DB
}

var _ ports.PipelineRepository = (*DatabaseAdapter)(nil)

// NewDatabaseAdapter initializes a new database adapter
func NewDatabaseAdapter() *DatabaseAdapter {
	return &DatabaseAdapter{DB: DB}
}

// SaveUser inserts a new user into the database, ensuring uniqueness
func (d *DatabaseAdapter) SaveUser(user *models.User) error {
	sqlDB, _ := d.DB.DB()
	sqlDB.Exec("DEALLOCATE ALL") // ✅ Reset prepared statements

	var existingUser models.User
	err := d.DB.Where("email = ?", user.Email).First(&existingUser).Error

	if err == nil {
		log.Printf("⚠️ User already exists with email: %s", user.Email)
		return nil
	} else if err != gorm.ErrRecordNotFound {
		return err
	}

	if user.UserID == uuid.Nil {
		user.UserID = uuid.New()
	}

	if err := d.DB.Create(user).Error; err != nil {
		log.Printf("❌ Failed to save user in DB: %v", err)
		return err
	}

	log.Println("✅ User saved successfully:", user.Email)
	return nil
}

// GetUserByID retrieves a user by their ID
func (d *DatabaseAdapter) GetUserByID(userID uuid.UUID) (*models.User, error) {
	sqlDB, _ := d.DB.DB()
	sqlDB.Exec("DEALLOCATE ALL")

	var user models.User
	if err := d.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates user details
func (d *DatabaseAdapter) UpdateUser(userID uuid.UUID, updates map[string]interface{}) error {
	sqlDB, _ := d.DB.DB()
	sqlDB.Exec("DEALLOCATE ALL")
	return d.DB.Model(&models.User{}).Where("user_id = ?", userID).Updates(updates).Error
}

// SavePipelineExecution saves pipeline execution details
func (d *DatabaseAdapter) SavePipelineExecution(execution *models.PipelineExecution) error {
	sqlDB, _ := d.DB.DB()
	sqlDB.Exec("DEALLOCATE ALL")
	return d.DB.Create(execution).Error
}

// UpdatePipelineExecution updates the pipeline execution status
func (d *DatabaseAdapter) UpdatePipelineExecution(execution *models.PipelineExecution) error {
	sqlDB, _ := d.DB.DB()
	sqlDB.Exec("DEALLOCATE ALL")
	return d.DB.Model(execution).Where("pipeline_id = ?", execution.PipelineID).Update("status", execution.Status).Error
}

// SaveExecutionLog saves execution logs
func (d *DatabaseAdapter) SaveExecutionLog(logEntry *models.ExecutionLog) error {
	sqlDB, _ := d.DB.DB()
	sqlDB.Exec("DEALLOCATE ALL")
	return d.DB.Create(logEntry).Error
}

// GetPipelineStatus retrieves pipeline execution status
func (d *DatabaseAdapter) GetPipelineStatus(pipelineID string) (string, error) {
	sqlDB, _ := d.DB.DB()
	sqlDB.Exec("DEALLOCATE ALL")

	var execution models.PipelineExecution
	if err := d.DB.Where("pipeline_id = ?", pipelineID).First(&execution).Error; err != nil {
		return "", err
	}
	return execution.Status, nil
}

// GetPipelinesByUser retrieves all pipelines for a user
func (d *DatabaseAdapter) GetPipelinesByUser(userID string) ([]models.PipelineExecution, error) {
	sqlDB, _ := d.DB.DB()
	sqlDB.Exec("DEALLOCATE ALL")

	var pipelines []models.PipelineExecution
	err := d.DB.Where("user_id = ?", userID).Find(&pipelines).Error
	return pipelines, err
}

// GetPipelineStages fetches all stages associated with a pipeline
func (d *DatabaseAdapter) GetPipelineStages(pipelineID uuid.UUID) ([]models.ExecutionLog, error) {
	sqlDB, _ := d.DB.DB()
	sqlDB.Exec("DEALLOCATE ALL")

	var stages []models.ExecutionLog
	if err := d.DB.Where("pipeline_id = ?", pipelineID).Find(&stages).Error; err != nil {
		return nil, err
	}
	return stages, nil
}




