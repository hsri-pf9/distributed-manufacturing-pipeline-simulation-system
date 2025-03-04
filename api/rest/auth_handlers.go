package rest

import (
	"encoding/json"
	"net/http"
	"log"

	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/services"
)

// AuthHandler handles authentication routes
type AuthHandler struct {
	Service *services.AuthService
}

// Credentials struct for user input
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterHandler handles user registration
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Debug: Print received data
	log.Printf("Received Register Request: Email=%s, Password=%s", creds.Email, creds.Password)

	userID, email, token, err := h.Service.RegisterUser(creds.Email, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Debug: Print success
	log.Printf("User Registered Successfully: ID=%s, Email=%s", userID, email)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"user_id": userID,
		"email":   email,
		"token":   token,
	})
}

// LoginHandler handles user login
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, email, token, err := h.Service.LoginUser(creds.Email, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"user_id": userID,
		"email":   email,
		"token":   token,
	})
}
