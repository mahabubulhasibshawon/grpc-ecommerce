// internal/application/auth_service.go
package application

import (
	"context"
	"errors"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/domain"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/ports"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/pkg/auth"
)

type AuthService struct {
	repo ports.OrderRepositoryPort
}

func NewAuthService(repo ports.OrderRepositoryPort) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, *domain.User, error) {
	user, err := s.repo.FindUserByUsername(ctx, username)
	if err != nil {
		return "", nil, err
	}
	if user == nil || user.Password != password {
		return "", nil, errors.New("invalid credentials")
	}
	token, err := auth.GenerateToken(username, user.ID)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (s *AuthService) Logout(ctx context.Context, userID int64) error {
	return nil
}