package auth

import (
	"context"
	"errors"

	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/domain"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/ports"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	repo ports.OrderRepositoryPort
}

func NewLoginService(repo ports.OrderRepositoryPort) *LoginService {
	return &LoginService{repo: repo}
}

func (s *LoginService) Login(ctx context.Context, username, password string) (string, *domain.User, error) {
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
