package secondary

import (
	// "log"
	"os"

	// "github.com/joho/godotenv"
	"github.com/nedpals/supabase-go"
)

// InitSupabaseClient initializes and returns a Supabase client
func InitSupabaseClient() *supabase.Client {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }

	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_API_KEY")

	return supabase.CreateClient(url, key)
}
