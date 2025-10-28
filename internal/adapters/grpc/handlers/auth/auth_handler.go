package auth

import (
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/application/auth"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/ports"
)

type AuthHandler struct {
	signupService *auth.SignupService
	loginService  *auth.LoginService
	logoutService *auth.LogoutService
}

func NewAuthHandler(repo ports.OrderRepositoryPort) *AuthHandler {
	return &AuthHandler{
		signupService: auth.NewSignupService(repo),
		loginService:  auth.NewLoginService(repo),
		logoutService: auth.NewLogoutService(),
	}
}
