package middleware

import (
	"context"
	// "errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/adapters/secondary"
	// "github.com/nedpals/supabase-go"
)

// // Middleware to authenticate requests using Supabase JWT token
// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
		
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
// 			c.Abort()
// 			return
// 		}

// 		// Extract token from "Bearer <token>"
// 		tokenParts := strings.Split(authHeader, " ")
// 		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
// 			c.Abort()
// 			return
// 		}
// 		token := tokenParts[1]

// 		// Validate token with Supabase
// 		client := secondary.InitSupabaseClient()
// 		user, err := client.Auth.User(context.Background(), token)
// 		if err != nil || user == nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
// 			c.Abort()
// 			return
// 		}

// 		// Store user ID in context for access in handlers
// 		c.Set("user_id", user.ID)
// 		c.Next()
// 	}
// }

// AuthMiddleware authenticates API requests using Supabase JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		// 🔹 Check if token is in the Authorization header (for normal API requests)
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
				token = tokenParts[1]
			}
		}

		// 🔹 If no token in header, check query params (for SSE)
		if token == "" {
			token = c.Query("token")
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			c.Abort()
			return
		}

		// 🔹 Validate token with Supabase
		client := secondary.InitSupabaseClient()
		user, err := client.Auth.User(context.Background(), token)
		if err != nil || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 🔹 Store user ID in context for handlers
		c.Set("user_id", user.ID)
		c.Next()
	}
}
