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

// RegisterUser registers a new user in Supabase and saves them in the database after email confirmation.
func (s *AuthService) RegisterUser(email, password string) (string, string, string, error) {
	log.Printf("[DEBUG] RegisterUser called with Email: %s", email)

	user, err := s.SupabaseClient.Auth.SignUp(context.Background(), supabase.UserCredentials{
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Printf("[ERROR] Supabase SignUp failed: %v", err)
		return "", "", "", errors.New("registration failed: " + err.Error())
	}

	if user == nil {
		log.Printf("[ERROR] Unexpected response from Supabase: user is nil")
		return "", "", "", errors.New("unexpected response from Supabase: user is nil")
	}

	log.Println("[INFO] Registration successful. Waiting for email confirmation.")
	fmt.Println("Please confirm your email before proceeding.")

	// ✅ After email verification, fetch user from Supabase
	session, err := s.SupabaseClient.Auth.SignIn(context.Background(), supabase.UserCredentials{
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Printf("[ERROR] Login after registration failed: %v", err)
		return "", "", "", errors.New("login after registration failed: " + err.Error())
	}

	// ✅ Ensure we get a valid user ID from Supabase
	if session.User.ID == "" {
		return "", "", "", errors.New("failed to retrieve user ID from Supabase")
	}

	// ✅ Convert user ID to UUID
	userUUID, err := uuid.Parse(session.User.ID)
	if err != nil {
		return "", "", "", errors.New("invalid user UUID: " + err.Error())
	}

	// ✅ Check if the user already exists in DB
	existingUser, _ := s.Repo.GetUserByID(userUUID)
	if existingUser == nil {
		newUser := &models.User{
			UserID: userUUID,
			Email:  session.User.Email,
			Role:   "worker", // Default role
		}

		// ✅ Ensure user is stored in DB
		if err := s.Repo.SaveUser(newUser); err != nil {
			log.Printf("[ERROR] Failed to save user in DB: %v", err)
			return "", "", "", errors.New("failed to save user in database")
		}
		log.Printf("[INFO] User saved successfully: ID=%s", userUUID.String())
	}

	return session.User.ID, session.User.Email, session.AccessToken, nil
}

// LoginUser authenticates a user and returns their details and token.
func (s *AuthService) LoginUser(email, password string) (string, string, string, error) {
	log.Printf("[DEBUG] LoginUser called with Email: %s", email)

	session, err := s.SupabaseClient.Auth.SignIn(context.Background(), supabase.UserCredentials{
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Printf("[ERROR] Supabase Login failed: %v", err)
		return "", "", "", errors.New("login failed: " + err.Error())
	}

	// ✅ Ensure we get valid user details
	if session.User.ID == "" || session.User.Email == "" || session.AccessToken == "" {
		return "", "", "", errors.New("unexpected response from Supabase: missing user details")
	}

	// ✅ Convert string UserID to uuid.UUID
	userUUID, err := uuid.Parse(session.User.ID)
	if err != nil {
		log.Printf("[ERROR] Invalid user UUID: %v", err)
		return "", "", "", errors.New("invalid user UUID: " + err.Error())
	}

	// ✅ Check if user exists in DB
	existingUser, _ := s.Repo.GetUserByID(userUUID)
	if existingUser == nil {
		// ✅ Save user details after email confirmation
		newUser := &models.User{
			UserID: userUUID,
			Email:  session.User.Email,
			Role:   "worker", // Default role
		}
		if err := s.Repo.SaveUser(newUser); err != nil {
			log.Printf("[ERROR] Failed to save user in DB: %v", err)
			return "", "", "", errors.New("failed to save user in the database")
		}
		log.Printf("[INFO] New user added to database: %s", session.User.Email)
	}

	return session.User.ID, session.User.Email, session.AccessToken, nil
}

// GetUserByID fetches user details
func (s *AuthService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	return s.Repo.GetUserByID(userID)
}

// UpdateUser updates user details
func (s *AuthService) UpdateUser(userID uuid.UUID, updates map[string]interface{}) error {
	return s.Repo.UpdateUser(userID, updates)
}
