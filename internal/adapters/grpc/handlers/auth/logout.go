package auth

import (
	"context"

	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/pkg/auth"
	pb "github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *AuthHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}

	err = h.server.AuthService.Logout(ctx, userID)
	if err != nil {
		return &pb.LogoutResponse{Message: err.Error(), Type: "error", Code: 400}, nil
	}
	return &pb.LogoutResponse{Message: "Successfully logged out", Type: "success", Code: 200}, nil
}
