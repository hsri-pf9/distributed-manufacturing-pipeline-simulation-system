package primary

import (
	"context"
	"errors"
	"log"

	proto "github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/api/grpc/proto"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/internal/core/services"
)

type AuthServer struct {
	proto.UnimplementedAuthServiceServer
	AuthService *services.AuthService
}

// Register handles user registration via gRPC
func (s *AuthServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	userID, email, token, err := s.AuthService.RegisterUser(req.Email, req.Password)
	if err != nil {
		log.Println("Registration error:", err)
		return nil, errors.New("registration failed")
	}

	return &proto.RegisterResponse{
		UserId: userID,
		Email:  email,
		Token:  token,
	}, nil
}

// Login handles user authentication via gRPC
func (s *AuthServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	userID, email, token, err := s.AuthService.LoginUser(req.Email, req.Password)
	if err != nil {
		log.Println("Login error:", err)
		return nil, errors.New("login failed")
	}

	return &proto.LoginResponse{
		UserId: userID,
		Email:  email,
		Token:  token,
	}, nil
}
