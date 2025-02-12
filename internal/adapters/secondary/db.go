package secondary

import (
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/models"
)

var DB *gorm.DB

func InitDatabase() {
	dsn := os.Getenv("SUPABASE_DB_URL")
	if dsn == "" {
		log.Fatal("SUPABASE_DB_URL environment variable is not set")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to Supabase database: %v", err)
	}

	if err := DB.AutoMigrate(&models.PipelineExecution{}, &models.ExecutionLog{}); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
}