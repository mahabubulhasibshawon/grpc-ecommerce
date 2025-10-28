package auth

import (
	"context"
	"errors"

	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/pkg/auth"
)

type LogoutService struct{}

func NewLogoutService() *LogoutService {
	return &LogoutService{}
}

func (s *LogoutService) Logout(ctx context.Context, userID int64) error {
	token, ok := ctx.Value("token").(string)
	if !ok {
		return errors.New("token not found in context")
	}
	auth.BlacklistToken(token)
	return nil
}