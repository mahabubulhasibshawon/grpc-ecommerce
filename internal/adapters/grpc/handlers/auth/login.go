package auth

import (
	"context"

	pb "github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/grpc/proto"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/pkg/auth"
)

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, user, err := h.server.AuthService.Login(ctx, req.Username, req.Password)
	if err != nil {
		return &pb.LoginResponse{Message: "Invalid credentials", Type: "error", Code: 400}, nil
	}
	token, err = auth.GenerateToken(req.Username, user.ID)
	if err != nil {
		return &pb.LoginResponse{Message: err.Error(), Type: "error", Code: 500}, nil
	}
	return &pb.LoginResponse{
		TokenType:    "Bearer",
		ExpiresIn:    432000,
		AccessToken:  token,
		RefreshToken: "dummy-refresh",
		Message:      "Logged in",
		Type:         "success",
		Code:         200,
	}, nil
}
