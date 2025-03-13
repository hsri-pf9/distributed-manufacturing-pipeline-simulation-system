// package main

// import (
// 	"log"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/api/rest"
// 	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/cmd/middleware"
// 	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/adapters/secondary"
// 	// "github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/domain"
// 	"github.com/gin-contrib/cors"
// 	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/services"
// )

// func main() {

// 	secondary.InitDatabase()

// 	// Initialize database repository
// 	dbRepo := secondary.NewDatabaseAdapter()

// 	// Initialize pipeline service
// 	pipelineService := services.NewPipelineService(dbRepo)

// 	// email := "harshsrivastava2404@gmail.com"
// 	// password := "harsh123"

// 	authService := services.NewAuthService(dbRepo)
// 	authMiddleware := middleware.AuthMiddleware()


// 	// Initialize REST API handler
// 	handler := &rest.PipelineHandler{Service: pipelineService}
// 	authHandler := &rest.AuthHandler{Service: authService}
// 	userHandler := &rest.UserHandler{Service: authService} // New user handler
// 	// Setup Gin router
// 	r := gin.Default()

// 	r.Use(cors.New(cors.Config{
// 		AllowOrigins:     []string{"http://localhost:3000"}, // Allows requests from any origin (including Postman Web & React)
// 		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
// 		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
// 		AllowCredentials: true,
// 	}))

// 	r.POST("/register", gin.WrapF(authHandler.RegisterHandler))
// 	r.POST("/login", gin.WrapF(authHandler.LoginHandler))

// 	// // User profile routes
// 	// r.GET("/user/:id", userHandler.GetUserProfile)  // Fetch user profile
// 	// r.PUT("/user/:id", userHandler.UpdateUserProfile) // Update user profil
// 	// r.GET("/pipelines", handler.GetUserPipelines)
// 	// r.GET("/pipelines/:id/stages", handler.GetPipelineStages) // ✅ New route



// 	// r.POST("/createpipelines", handler.CreatePipeline)
// 	// r.POST("/pipelines/:id/start", handler.StartPipeline)
// 	// r.GET("/pipelines/:id/status", handler.GetPipelineStatus)
// 	// r.POST("/pipelines/:id/cancel", handler.CancelPipeline)

// 	r.GET("/user/:id", authMiddleware, userHandler.GetUserProfile)
// 	r.PUT("/user/:id", authMiddleware, userHandler.UpdateUserProfile)

// 	r.GET("/pipelines", authMiddleware, handler.GetUserPipelines)
// 	r.GET("/pipelines/:id/stages", authMiddleware, handler.GetPipelineStages)

// 	r.POST("/createpipelines", authMiddleware, handler.CreatePipeline)
// 	r.POST("/pipelines/:id/start", authMiddleware, handler.StartPipeline)
// 	r.GET("/pipelines/:id/status", authMiddleware, handler.GetPipelineStatus)
// 	r.POST("/pipelines/:id/cancel", authMiddleware, handler.CancelPipeline)

// 	// Start server
// 	log.Println("Starting API server on port 8080...")
// 	if err := http.ListenAndServe(":8080", r); err != nil {
// 		log.Fatalf("Failed to start server: %v", err)
// 	}
// }

package main

import (
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/api/rest"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/cmd/middleware"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/adapters/primary"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/adapters/secondary"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/services"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/utils"

	"github.com/gin-contrib/cors"
	proto "github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/api/grpc/proto/auth"
	pipeline_proto "github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/api/grpc/proto/pipeline"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func startRESTServer(authService *services.AuthService, pipelineService *services.PipelineService, sseManager *utils.SSEManager, wg *sync.WaitGroup) {
	defer wg.Done()

	authMiddleware := middleware.AuthMiddleware()

	// Initialize REST API handlers
	// handler := &rest.PipelineHandler{Service: pipelineService}
	handler := &rest.PipelineHandler{Service: pipelineService, SSE: sseManager}
	authHandler := &rest.AuthHandler{Service: authService}
	userHandler := &rest.UserHandler{Service: authService}

	// Setup Gin router
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
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
	log.Println("Starting REST API & Frontend on port 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start REST server: %v", err)
	}
}

func startGRPCServer(authService *services.AuthService, pipelineService *services.PipelineService, wg *sync.WaitGroup) {
	defer wg.Done()

	// Create gRPC server
	grpcServer := grpc.NewServer()
	authServer := &primary.AuthServer{AuthService: authService}
	pipelineServer := &primary.PipelineServer{Service: pipelineService}

	// Register gRPC services
	proto.RegisterAuthServiceServer(grpcServer, authServer)
	pipeline_proto.RegisterPipelineServiceServer(grpcServer, pipelineServer)
	reflection.Register(grpcServer)

	// Start gRPC server
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	log.Println("Starting gRPC server on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}

// ✅ Correctly serve React frontend
func startFrontendServer(wg *sync.WaitGroup) {
	defer wg.Done()

	fs := http.FileServer(http.Dir("cmd/api-server/build")) // Ensure correct path
	http.Handle("/", fs)

	log.Println("Starting Frontend server on port 3000...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalf("Failed to start Frontend server: %v", err)
	}
}

func main() {
	// Initialize database
	secondary.InitDatabase()
	dbRepo := secondary.NewDatabaseAdapter()

	sseManager := utils.NewSSEManager()

	// Initialize services
	authService := services.NewAuthService(dbRepo)
	pipelineService := services.NewPipelineService(dbRepo,sseManager)
	


	var wg sync.WaitGroup
	wg.Add(3)

	// Start REST API server
	go startRESTServer(authService, pipelineService, sseManager, &wg)

	// Start gRPC server
	go startGRPCServer(authService, pipelineService, &wg)

	// go startFrontendServer(&wg) 

	// Wait for both servers to run
	wg.Wait()
}
