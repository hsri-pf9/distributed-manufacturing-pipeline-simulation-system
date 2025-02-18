package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/api/rest"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/adapters/secondary"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/domain"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/services"
)

func main() {

	secondary.InitDatabase()

	// Initialize database repository
	dbRepo := secondary.NewDatabaseAdapter()

	// Initialize pipeline orchestrators
	sequentialOrchestrator := domain.NewSequentialPipelineOrchestrator(dbRepo)
	parallelOrchestrator := domain.NewParallelPipelineOrchestrator(dbRepo)

	// Initialize pipeline service
	pipelineService := services.NewPipelineService(sequentialOrchestrator, parallelOrchestrator, dbRepo)

	// Initialize REST API handler
	handler := &rest.PipelineHandler{Service: pipelineService}

	// Setup Gin router
	r := gin.Default()
	r.POST("/pipelines", handler.CreatePipeline)
	r.POST("/pipelines/:id/start", handler.StartPipeline)
	r.GET("/pipelines/:id/status", handler.GetPipelineStatus)
	r.POST("/pipelines/:id/cancel", handler.CancelPipeline)

	// Start server
	log.Println("Starting API server on port 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

