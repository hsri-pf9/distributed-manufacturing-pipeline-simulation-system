// package models

// import (
// 	"time"
// 	"github.com/google/uuid"
// )

// // PipelineExecution stores pipeline execution details
// type PipelineExecution struct {
// 	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
// 	PipelineID uuid.UUID `gorm:"type:uuid;not null"`
// 	Status     string    `gorm:"type:varchar(50);not null"`
// 	CreatedAt  time.Time `gorm:"autoCreateTime"`
// 	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
// }

// // ExecutionLog stores logs related to pipeline execution
// type ExecutionLog struct {
// 	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
// 	StageID    uuid.UUID `gorm:"type:uuid;not null"`
// 	PipelineID uuid.UUID `gorm:"type:uuid;not null"`
// 	Status     string    `gorm:"type:varchar(50);not null"`
// 	ErrorMsg   string    `gorm:"type:text"`
// 	Timestamp  time.Time `gorm:"autoCreateTime"`
// }

package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	UserID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name   string    `gorm:"type:varchar(100);not null"`
	Email  string    `gorm:"type:varchar(100);unique;not null"`
	Role   string    `gorm:"type:varchar(20);not null;default:'worker';check:role IN ('super_admin', 'admin', 'manager', 'worker')"`
}


// PipelineExecution stores pipeline execution details for a user
type PipelineExecution struct {
	PipelineID uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	Status     string    `gorm:"type:varchar(50);not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

// ExecutionLog stores logs related to pipeline execution stages
type ExecutionLog struct {
	StageID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	PipelineID uuid.UUID `gorm:"type:uuid;not null;index"`
	Status     string    `gorm:"type:varchar(50);not null"`
	ErrorMsg   string    `gorm:"type:text"`
	Timestamp  time.Time `gorm:"autoCreateTime"`

	PipelineExecution PipelineExecution `gorm:"foreignKey:PipelineID;constraint:OnDelete:CASCADE;"`
}
