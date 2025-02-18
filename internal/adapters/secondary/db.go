package secondary

import (
    "log"
    "os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/joho/godotenv"
    "github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/models"
)

var DB *gorm.DB

func InitDatabase() {
    // Load environment variables from .env file
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: No .env file found. Using system environment variables.")
    }

    dsn := os.Getenv("SUPABASE_DB")
    if dsn == "" {
        log.Fatal("SUPABASE_DB environment variable is not set")
    }

	log.Printf("Connecting to database: %s", dsn)

    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to Supabase database: %v", err)
    }
    log.Println("Database connection established.")

    if err := DB.AutoMigrate(&models.PipelineExecution{}, &models.ExecutionLog{}); err != nil {
        log.Fatalf("Database migration failed: %v", err)
    }
    log.Println("Database migration completed.")
}
