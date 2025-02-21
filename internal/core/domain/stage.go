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
	log.Printf("Executing stage: %s\n", s.ID)
	return input, nil
}

func (s *BaseStage) HandleError(ctx context.Context, err error) error {
	log.Println("Error in stage execution:", err)
	return errors.New("stage execution failed: " + err.Error())
}

func (s *BaseStage) Rollback(ctx context.Context, input interface{}) error {
	log.Println("Rolling back stage:", s.ID)
	return nil
}