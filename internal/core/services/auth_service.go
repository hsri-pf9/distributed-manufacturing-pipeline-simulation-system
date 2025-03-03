package services

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/adapters/secondary"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/models"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/ports"
	"github.com/nedpals/supabase-go"
)

// AuthService handles user authentication
type AuthService struct {
	SupabaseClient *supabase.Client
	Repo           ports.PipelineRepository
}

// NewAuthService initializes the AuthService
func NewAuthService(repo ports.PipelineRepository) *AuthService {
	return &AuthService{
		SupabaseClient: secondary.InitSupabaseClient(),
		Repo:           repo,
	}
}

// // RegisterUser registers a new user and saves their details in the database
// func (s *AuthService) RegisterUser(email, password string) (string, string, string, error) {
// 	user, err := s.SupabaseClient.Auth.SignUp(context.Background(), supabase.UserCredentials{
// 		Email:    email,
// 		Password: password,
// 	})
// 	if err != nil {
// 		return "", "", "", errors.New("registration failed: " + err.Error())
// 	}

// 	if user == nil {
// 		return "", "", "", errors.New("unexpected response from Supabase: user is nil")
// 	}

// 	// Since SignUp does NOT return a session, manually log in to get a session token
// 	session, err := s.SupabaseClient.Auth.SignIn(context.Background(), supabase.UserCredentials{
// 		Email:    email,
// 		Password: password,
// 	})
// 	if err != nil {
// 		return "", "", "", errors.New("login after registration failed: " + err.Error())
// 	}

// 	// Save user details in the database
// 	newUser := &models.User{
// 		UserID: uuid.MustParse(user.ID),
// 		Email:  user.Email,
// 		Name:   "Harsh Srivastava", // Name should be updated from the UI later
// 		Role:   "admin", // Default role
// 	}

// 	if err := s.Repo.SaveUser(newUser); err != nil {
// 		return "", "", "", errors.New("failed to save user in the database: " + err.Error())
// 	}

// 	return user.ID, user.Email, session.AccessToken, nil
// }

// // LoginUser authenticates a user and returns ID, email, and token
// func (s *AuthService) LoginUser(email, password string) (string, string, string, error) {
// 	session, err := s.SupabaseClient.Auth.SignIn(context.Background(), supabase.UserCredentials{
// 		Email:    email,
// 		Password: password,
// 	})
// 	if err != nil {
// 		return "", "", "", errors.New("login failed: " + err.Error())
// 	}

// 	if session.User.ID == "" || session.User.Email == "" || session.AccessToken == "" {
// 		return "", "", "", errors.New("unexpected response from Supabase: session or user is nil")
// 	}

// 	return session.User.ID, session.User.Email, session.AccessToken, nil
// }

func (s *AuthService) RegisterUser(email, password string) (string, string, string, error) {
	log.Printf("[DEBUG] RegisterUser called with Email: %s", email)

	user, err := s.SupabaseClient.Auth.SignUp(context.Background(), supabase.UserCredentials{
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Printf("[ERROR] Supabase SignUp failed: %v", err) // Log error
		return "", "", "", errors.New("registration failed: " + err.Error())
	}

	if user == nil {
		log.Printf("[ERROR] Unexpected response from Supabase: user is nil") // Log if user is nil
		return "", "", "", errors.New("unexpected response from Supabase: user is nil")
	}

	// ✅ Wait for the email confirmation before proceeding
	log.Println("[INFO] Registration successful. Waiting for email confirmation.")
	// ✅ Wait for the email confirmation before proceeding
	fmt.Println("Please confirm your email before proceeding.")

	return user.ID, user.Email, "", nil // Don't return token yet
}

func (s *AuthService) LoginUser(email, password string) (string, string, string, error) {
	session, err := s.SupabaseClient.Auth.SignIn(context.Background(), supabase.UserCredentials{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", "", "", errors.New("login failed: " + err.Error())
	}

	if session.User.ID == "" || session.User.Email == "" || session.AccessToken == "" {
		return "", "", "", errors.New("unexpected response from Supabase: session or user is nil")
	}

	// ✅ Convert string UserID to uuid.UUID
	userUUID, err := uuid.Parse(session.User.ID)
	if err != nil {
		return "", "", "", errors.New("invalid user UUID: " + err.Error())
	}

	// ✅ Check if user already exists in DB
	existingUser, _ := s.Repo.GetUserByID(userUUID)
	if existingUser == nil {
		// ✅ Save user details after email confirmation
		newUser := &models.User{
			UserID: userUUID,
			Email:  session.User.Email,
			// Name:   "Harsh Srivastava",
			// Role:   "admin",
		}
		if err := s.Repo.SaveUser(newUser); err != nil {
			return "", "", "", errors.New("failed to save user in the database: " + err.Error())
		}
	}

	return session.User.ID, session.User.Email, session.AccessToken, nil
}

