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

    // Enable UUID extension
    err = DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
    if err != nil {
        log.Fatalf("Failed to enable uuid-ossp extension: %v", err)
    }
    log.Println("UUID-OSSP extension enabled.")

    // // **Migrate User table first**
    // if err := DB.AutoMigrate(&models.Customer{}); err != nil {
    //     log.Fatalf("Failed to migrate Customer table: %v", err)
    // }

    // **Migrate User table first**
    if err := DB.AutoMigrate(&models.User{}); err != nil {
        log.Fatalf("Failed to migrate User table: %v", err)
    }

    // **Migrate PipelineExecution next**
    if err := DB.AutoMigrate(&models.PipelineExecution{}); err != nil {
        log.Fatalf("Failed to migrate PipelineExecution table: %v", err)
    }

    // **Migrate ExecutionLog last (depends on PipelineExecution)**
    if err := DB.AutoMigrate(&models.ExecutionLog{}); err != nil {
        log.Fatalf("Failed to migrate ExecutionLog table: %v", err)
    }

    log.Println("Database migration completed successfully.")
}
