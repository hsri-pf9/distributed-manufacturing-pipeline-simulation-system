package services

import (
	"context"
	"errors"

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

// RegisterUser registers a new user and saves their details in the database
func (s *AuthService) RegisterUser(email, password string) (string, string, string, error) {
	user, err := s.SupabaseClient.Auth.SignUp(context.Background(), supabase.UserCredentials{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", "", "", errors.New("registration failed: " + err.Error())
	}

	if user == nil {
		return "", "", "", errors.New("unexpected response from Supabase: user is nil")
	}

	// Since SignUp does NOT return a session, manually log in to get a session token
	session, err := s.SupabaseClient.Auth.SignIn(context.Background(), supabase.UserCredentials{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", "", "", errors.New("login after registration failed: " + err.Error())
	}

	// Save user details in the database
	newUser := &models.User{
		UserID: uuid.MustParse(user.ID),
		Email:  user.Email,
		Name:   "", // Name should be updated from the UI later
		Role:   "worker", // Default role
	}

	if err := s.Repo.SaveUser(newUser); err != nil {
		return "", "", "", errors.New("failed to save user in the database: " + err.Error())
	}

	return user.ID, user.Email, session.AccessToken, nil
}

// LoginUser authenticates a user and returns ID, email, and token
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

	return session.User.ID, session.User.Email, session.AccessToken, nil
}
