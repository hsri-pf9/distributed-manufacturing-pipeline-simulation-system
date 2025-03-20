package secondary

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "github.com/joho/godotenv"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/models"
)

var DB *gorm.DB

// InitDatabase initializes the database connection
func InitDatabase() {
	// // Load environment variables
	// if err := godotenv.Load(); err != nil {
	// 	log.Println("⚠️ Warning: No .env file found. Using system environment variables.")
	// }

	dsn := os.Getenv("SUPABASE_DB")
	if dsn == "" {
		log.Fatal("❌ SUPABASE_DB environment variable is not set")
	}

	log.Printf("🔗 Connecting to database: %s", dsn)

	var err error
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // ✅ Disable prepared statement caching
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	log.Println("✅ Database connection established.")

	// Run database migrations
	if err := DB.AutoMigrate(&models.User{}, &models.PipelineExecution{}, &models.ExecutionLog{}); err != nil {
		log.Fatalf("❌ Database migration failed: %v", err)
	}

	log.Println("✅ Database migration completed.")
}

// CloseDatabase closes the database connection
func CloseDatabase() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Println("⚠️ Warning: Unable to close database connection properly.")
		return
	}
	sqlDB.Close()
	log.Println("✅ Database connection closed.")
}






