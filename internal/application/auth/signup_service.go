package auth

import (
	"context"
	"errors"

	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/domain"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

type SignupService struct {
	repo ports.OrderRepositoryPort
}

func NewSignupService(repo ports.OrderRepositoryPort) *SignupService {
	return &SignupService{repo: repo}
}

func (s *SignupService) Signup(ctx context.Context, username, password string) (*domain.User, error) {
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
