package auth

import "github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/grpc"

type AuthHandler struct {
	server *grpc.Server
}

func NewAuthHandler(s *grpc.Server) *AuthHandler {
	return &AuthHandler{server: s}
}
