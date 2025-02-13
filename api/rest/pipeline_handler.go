package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/services"
)

type PipelineHandler struct {
	Service *services.PipelineService
}

type CreatePipelineRequest struct {
	Stages    int  `json:"stages"`
	IsParallel bool `json:"is_parallel"`
}

func (h *PipelineHandler) CreatePipeline(c *gin.Context) {
	var req CreatePipelineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	pipelineID, err := h.Service.CreatePipeline(req.Stages, req.IsParallel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create pipeline"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Pipeline created", "pipeline_id": pipelineID})
}

type StartPipelineRequest struct {
	Input      interface{} `json:"input"`
	IsParallel bool        `json:"is_parallel"`
}

func (h *PipelineHandler) StartPipeline(c *gin.Context) {
	pipelineID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	var req StartPipelineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	go func() {
		h.Service.StartPipeline(context.Background(), pipelineID, req.Input, req.IsParallel)
	}()

	c.JSON(http.StatusAccepted, gin.H{"message": "Pipeline execution started", "pipeline_id": pipelineID})
}

func (h *PipelineHandler) GetPipelineStatus(c *gin.Context) {
	pipelineID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	status, err := h.Service.GetPipelineStatus(pipelineID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pipeline not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pipeline_id": pipelineID, "status": status})
}

func (h *PipelineHandler) CancelPipeline(c *gin.Context) {
	pipelineID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	err = h.Service.CancelPipeline(pipelineID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel pipeline"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pipeline cancelled", "pipeline_id": pipelineID})
}