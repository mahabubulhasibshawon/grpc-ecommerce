package grpc

import (
	"context"
	"strings"

	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == "/order.OrderService/Login" || info.FullMethod == "/order.OrderService/Signup" {
		return handler(ctx, req)
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}
	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing authorization")
	}
	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	claims, err := auth.ValidateToken(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	ctx = context.WithValue(ctx, "userID", claims.UserID)
	return handler(ctx, req)
}
