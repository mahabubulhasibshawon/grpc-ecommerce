package auth

import (
	"context"
	pb "github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/grpc/proto"
)

func (h *AuthHandler) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	_, err := h.signupService.Signup(ctx, req.Username, req.Password)
	if err != nil {
		return &pb.SignupResponse{Message: err.Error(), Type: "error", Code: 400}, nil
	}
	return &pb.SignupResponse{
		Message: "User registered successfully",
		Type:    "success",
		Code:    200,
	}, nil
}
