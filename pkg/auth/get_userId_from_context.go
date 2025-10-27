package auth

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/metadata"
)

func GetUserIDFromContext(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, errors.New("missing metadata")
	}
	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return 0, errors.New("missing authorization")
	}
	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	claims, err := ValidateToken(token)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
