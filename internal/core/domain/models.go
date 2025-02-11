package domain

import (
	"time"
	"github.com/google/uuid"
)

// PipelineExecution stores pipeline execution details
type PipelineExecution struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	PipelineID uuid.UUID `gorm:"type:uuid;not null"`
	Status     string    `gorm:"type:varchar(50);not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

// ExecutionLog stores logs related to pipeline execution
type ExecutionLog struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	StageID    uuid.UUID `gorm:"type:uuid;not null"`
	PipelineID uuid.UUID `gorm:"type:uuid;not null"`
	Status     string    `gorm:"type:varchar(50);not null"`
	ErrorMsg   string    `gorm:"type:text"`
	Timestamp  time.Time `gorm:"autoCreateTime"`
}

