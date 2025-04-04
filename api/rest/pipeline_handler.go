package rest

import (
	"context"
	"net/http"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/services"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/utils"
)

type PipelineHandler struct {
	Service *services.PipelineService
	SSE     *utils.SSEManager // 🔹 Added SSE Manager
}

type CreatePipelineRequest struct {
	Stages     int  `json:"stages"`
	IsParallel bool `json:"is_parallel"`
	UserID     uuid.UUID `json:"user_id"` // Extracted from the request
}

// CreatePipeline handles pipeline creation
func (h *PipelineHandler) CreatePipeline(c *gin.Context) {
	var req CreatePipelineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.UserID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	pipelineID, err := h.Service.CreatePipeline(req.UserID, req.Stages, req.IsParallel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create pipeline"})
		return
	}

	h.SSE.BroadcastUpdate("Pipeline " + pipelineID.String() + " has been created.")

	c.JSON(http.StatusAccepted, gin.H{"message": "Pipeline created", "pipeline_id": pipelineID})
}

type StartPipelineRequest struct {
	Input      interface{} `json:"input"`
	IsParallel bool        `json:"is_parallel"`
	UserID     uuid.UUID   `json:"user_id"`
}

// StartPipeline handles pipeline execution
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

	if req.UserID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	go func() {
		h.Service.StartPipeline(context.Background(), req.UserID, pipelineID, req.Input, req.IsParallel)

		// 🔹 Notify clients about execution start
		h.SSE.BroadcastUpdate("Pipeline " + pipelineID.String() + " has started execution.")
	}()

	c.JSON(http.StatusAccepted, gin.H{"message": "Pipeline execution started", "pipeline_id": pipelineID})
}

type GetPipelineStatusRequest struct {
	IsParallel bool `json:"is_parallel"`
}

// GetPipelineStatus retrieves the current status of a pipeline
func (h *PipelineHandler) GetPipelineStatus(c *gin.Context) {
	pipelineID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	// var req GetPipelineStatusRequest
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
	// 	return
	// }

	// Get the "is_parallel" query parameter (defaults to false if not provided)
	isParallel := c.DefaultQuery("is_parallel", "false") == "true"

	// status, err := h.Service.GetPipelineStatus(pipelineID, req.IsParallel)
	// if err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Pipeline not found"})
	// 	return
	// }
	status, err := h.Service.GetPipelineStatus(pipelineID, isParallel)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pipeline not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pipeline_id": pipelineID, "status": status})
}

type CancelPipelineRequest struct {
	IsParallel bool      `json:"is_parallel"`
	UserID     uuid.UUID `json:"user_id"`
}

// CancelPipeline cancels an ongoing pipeline execution
func (h *PipelineHandler) CancelPipeline(c *gin.Context) {
	pipelineID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	var req CancelPipelineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.UserID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	err = h.Service.CancelPipeline(pipelineID, req.UserID, req.IsParallel)
	if err != nil {
		log.Printf("Error cancelling pipeline: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel pipeline"})
		return
	}

	// 🔹 Notify clients about pipeline cancellation
	h.SSE.BroadcastUpdate("Pipeline " + pipelineID.String() + " has been cancelled.")

	c.JSON(http.StatusOK, gin.H{"message": "Pipeline cancelled", "pipeline_id": pipelineID})
}

func (h *PipelineHandler) GetUserPipelines(c *gin.Context) {
	userID := c.Query("user_id") // Fetch user_id from query parameters
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	pipelines, err := h.Service.GetPipelinesByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pipelines"})
		return
	}

	c.JSON(http.StatusOK, pipelines)
}

// GetPipelineStages fetches the stages of a pipeline
func (h *PipelineHandler) GetPipelineStages(c *gin.Context) {
	pipelineID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	stages, err := h.Service.GetPipelineStages(pipelineID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pipeline stages"})
		return
	}

	c.JSON(http.StatusOK, stages)
}

