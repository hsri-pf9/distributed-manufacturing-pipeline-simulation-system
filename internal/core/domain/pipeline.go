package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"log"
)

type PipelineOrchestrator struct {
	ID uuid.UUID
	Stages []Stage
}

func NewPipelineOrchestrator() *PipelineOrchestrator {
	return &PipelineOrchestrator{
		ID: uuid.New(),
		Stages: []Stage{},
	}
}

//Addition of the stages
func (p *PipelineOrchestrator) AddStage(stage Stage) error {
	if(stage == nil){
		return errors.New("errors cannot be nil")
	}
	p.Stages = append(p.Stages, stage)
	return nil
}

//Execution of all the stages SEQUENTIALLY
func (p *PipelineOrchestrator) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	var result interface{} = input
	var err error

	for _, stage := range p.Stages {
		log.Println("Executing stage:", stage.GetID())

		result, err = stage.Execute(ctx, result)
		if err != nil {
			stage.HandleError(ctx, err)
			p.Rollback(ctx, result)
			return nil, err
		}
	}
	return result,nil
}

func (p *PipelineOrchestrator) Rollback(ctx context.Context, input interface{}) {
	for _, stage := range p.Stages {
		stage.Rollback(ctx, input)
	}
}