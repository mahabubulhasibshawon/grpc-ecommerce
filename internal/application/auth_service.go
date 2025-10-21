// internal/application/auth_service.go
package application

import (
	"context"
	"errors"

	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/domain"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/ports"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo ports.OrderRepositoryPort
}

func NewAuthService(repo ports.OrderRepositoryPort) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Signup(ctx context.Context, username, password string) (*domain.User, error) {
	if username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	user, err := s.repo.CreateUser(ctx, username, string(hashedPassword))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, *domain.User, error) {
	user, err := s.repo.FindUserByUsername(ctx, username)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
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
