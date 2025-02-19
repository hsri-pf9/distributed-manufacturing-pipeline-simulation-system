package services

import (
	"context"
	"errors"

	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/adapters/secondary"
	"github.com/nedpals/supabase-go"
)

// AuthService handles user authentication
type AuthService struct {
	SupabaseClient *supabase.Client
}

// NewAuthService initializes the AuthService
func NewAuthService() *AuthService {
	return &AuthService{
		SupabaseClient: secondary.InitSupabaseClient(),
	}
}

// RegisterUser registers a new user and logs them in to get a session token
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


