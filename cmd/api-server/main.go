package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/api/rest"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/adapters/secondary"
	// "github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/domain"
	"github.com/gin-contrib/cors"
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
	userHandler := &rest.UserHandler{Service: authService} // New user handler
	// Setup Gin router
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Allows requests from any origin (including Postman Web & React)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.POST("/register", gin.WrapF(authHandler.RegisterHandler))
	r.POST("/login", gin.WrapF(authHandler.LoginHandler))

	// User profile routes
	r.GET("/user/:id", userHandler.GetUserProfile)  // Fetch user profile
	r.PUT("/user/:id", userHandler.UpdateUserProfile) // Update user profil
	r.GET("/pipelines", handler.GetUserPipelines)


	r.POST("/createpipelines", handler.CreatePipeline)
	r.POST("/pipelines/:id/start", handler.StartPipeline)
	r.GET("/pipelines/:id/status", handler.GetPipelineStatus)
	r.POST("/pipelines/:id/cancel", handler.CancelPipeline)

	// Start server
	log.Println("Starting API server on port 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

