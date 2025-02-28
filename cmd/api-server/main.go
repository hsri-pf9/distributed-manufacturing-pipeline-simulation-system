package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/api/rest"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/adapters/secondary"
	// "github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/domain"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/services"
)

func main() {

	secondary.InitDatabase()

	// Initialize database repository
	dbRepo := secondary.NewDatabaseAdapter()

	// Initialize pipeline service
	pipelineService := services.NewPipelineService(dbRepo)

	// email := "harshsrivastava2404@gmail.com"
	// password := "harsh123"

	authService := services.NewAuthService(dbRepo)

	// Initialize REST API handler
	handler := &rest.PipelineHandler{Service: pipelineService}
	authHandler := &rest.AuthHandler{Service: authService}
	// Setup Gin router
	r := gin.Default()
	r.POST("/register", gin.WrapF(authHandler.RegisterHandler))
	r.POST("/login", gin.WrapF(authHandler.LoginHandler))

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

