package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid;primaryKey"`
    Name   string    `gorm:"type:varchar(100)"`
    Email  string    `gorm:"type:varchar(100);unique;not null"`
    Role   string    `gorm:"type:varchar(20);not null;default:'worker';check:role IN ('super_admin', 'admin', 'manager', 'worker')"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

    PipelineExecutions []PipelineExecution `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

type PipelineExecution struct {
    PipelineID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
    Status     string    `gorm:"type:varchar(50);not null"`
    CreatedAt  time.Time `gorm:"autoCreateTime"`
    UpdatedAt  time.Time `gorm:"autoUpdateTime"`

	// One PipelineExecution can have multiple ExecutionLogs
	ExecutionLogs []ExecutionLog `gorm:"foreignKey:PipelineID;constraint:OnDelete:CASCADE;"`
}

// ExecutionLog stores logs related to pipeline execution stages
type ExecutionLog struct {
    StageID    uuid.UUID `gorm:"type:uuid;primaryKey"`
    PipelineID uuid.UUID `gorm:"type:uuid;not null;index"`
    Status     string    `gorm:"type:varchar(50);not null"`
    ErrorMsg   string    `gorm:"type:text"`
    Timestamp  time.Time `gorm:"autoCreateTime"`
}