package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"log"
)

type Stage interface {
	GetID() uuid.UUID
	Execute(ctx context.Context, input interface{}) (interface{}, error)
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

func (s *BaseStage) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	log.Printf("Executing stage: %s with input: %v\n", s.ID, input)

	// Simulate execution failure for testing error handling
	if input == nil {
		err := errors.New("input is nil, stage execution failed")
		log.Printf("Stage %s execution failed: %v", s.ID, err)
		return nil, err
	}

	log.Printf("Stage %s executed successfully", s.ID)
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