package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/api/rest"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/cmd/middleware"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/adapters/secondary"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/services"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/utils"

	"github.com/gin-contrib/cors"
)

func startRESTServer(authService *services.AuthService, pipelineService *services.PipelineService, sseManager *utils.SSEManager, wg *sync.WaitGroup) {
	defer wg.Done()

	authMiddleware := middleware.AuthMiddleware()

	// Initialize REST API handlers
	handler := &rest.PipelineHandler{Service: pipelineService, SSE: sseManager}
	authHandler := &rest.AuthHandler{Service: authService}
	userHandler := &rest.UserHandler{Service: authService}

	// Setup Gin router
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{"http://localhost:3000"},
		AllowOrigins:     []string{"http://localhost:30080", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Public routes
	r.POST("/register", gin.WrapF(authHandler.RegisterHandler))
	r.POST("/login", gin.WrapF(authHandler.LoginHandler))

	// Protected routes
	r.GET("/user/:id", authMiddleware, userHandler.GetUserProfile)
	r.PUT("/user/:id", authMiddleware, userHandler.UpdateUserProfile)

	r.GET("/pipelines", authMiddleware, handler.GetUserPipelines)
	r.GET("/pipelines/:id/stages", authMiddleware, handler.GetPipelineStages)

	r.POST("/createpipelines", authMiddleware, handler.CreatePipeline)
	r.POST("/pipelines/:id/start", authMiddleware, handler.StartPipeline)
	r.GET("/pipelines/:id/status", authMiddleware, handler.GetPipelineStatus)
	r.POST("/pipelines/:id/cancel", authMiddleware, handler.CancelPipeline)

	// SSE Route
	r.GET("/pipelines/:id/stream", authMiddleware, sseManager.RegisterClient)

	// Start REST API server
	log.Println("🚀 Starting REST API & Frontend on port 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("❌ Failed to start REST server: %v", err)
	}
}

func main() {
	// Initialize database
	secondary.InitDatabase()
	defer secondary.CloseDatabase()

	dbRepo := secondary.NewDatabaseAdapter()
	sseManager := utils.NewSSEManager()

	// Initialize services
	authService := services.NewAuthService(dbRepo)
	pipelineService := services.NewPipelineService(dbRepo, sseManager)

	var wg sync.WaitGroup
	wg.Add(1) // Only 1 (REST)

	// Start REST API server
	go startRESTServer(authService, pipelineService, sseManager, &wg)

	// Wait for server
	wg.Wait()
}
