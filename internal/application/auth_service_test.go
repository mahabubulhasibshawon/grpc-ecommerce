// internal/application/auth_service_test.go
package application

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/domain"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/ports"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Signup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := ports.NewMockOrderRepositoryPort(ctrl)
	svc := NewAuthService(mockRepo)

	tests := []struct {
		name      string
		username  string
		password  string
		mockSetup func()
		wantErr   bool
		errMsg    string
	}{
		{
			name:     "Successful signup",
			username: "testuser@example.com",
			password: "securepass",
			mockSetup: func() {
				hashed, _ := bcrypt.GenerateFromPassword([]byte("securepass"), bcrypt.DefaultCost)
				mockRepo.EXPECT().CreateUser(gomock.Any(), "testuser@example.com", gomock.Any()).Return(&domain.User{ID: 1, Username: "testuser@example.com", Password: string(hashed)}, nil)
			},
			wantErr: false,
		},
		{
			name:      "Missing username",
			username:  "",
			password:  "securepass",
			mockSetup: func() {},
			wantErr:   true,
			errMsg:    "username and password are required",
		},
		{
			name:     "Repository error",
			username: "testuser@example.com",
			password: "securepass",
			mockSetup: func() {
				mockRepo.EXPECT().CreateUser(gomock.Any(), "testuser@example.com", gomock.Any()).Return(nil, errors.New("username already exists"))
			},
			wantErr: true,
			errMsg:  "username already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			user, err := svc.Signup(context.Background(), tt.username, tt.password)
			if tt.wantErr {
				if err == nil || err.Error() != tt.errMsg {
					t.Errorf("Signup() error = %v, wantErr %v, errMsg %v", err, tt.wantErr, tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("Signup() unexpected error: %v", err)
			}
			if user == nil || user.Username != tt.username {
				t.Errorf("Signup() user = %v, want username %v", user, tt.username)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := ports.NewMockOrderRepositoryPort(ctrl)
	svc := NewAuthService(mockRepo)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("securepass"), bcrypt.DefaultCost)

	tests := []struct {
		name      string
		username  string
		password  string
		mockSetup func()
		wantErr   bool
		errMsg    string
	}{
		{
			name:     "Successful login",
			username: "testuser@example.com",
			password: "securepass",
			mockSetup: func() {
				mockRepo.EXPECT().FindUserByUsername(gomock.Any(), "testuser@example.com").Return(&domain.User{ID: 1, Username: "testuser@example.com", Password: string(hashed)}, nil)
			},
			wantErr: false,
		},
		{
			name:     "Invalid credentials",
			username: "testuser@example.com",
			password: "wrongpass",
			mockSetup: func() {
				mockRepo.EXPECT().FindUserByUsername(gomock.Any(), "testuser@example.com").Return(&domain.User{ID: 1, Username: "testuser@example.com", Password: string(hashed)}, nil)
			},
			wantErr: true,
			errMsg:  "invalid credentials",
		},
		{
			name:     "User not found",
			username: "testuser@example.com",
			password: "securepass",
			mockSetup: func() {
				mockRepo.EXPECT().FindUserByUsername(gomock.Any(), "testuser@example.com").Return(nil, nil)
			},
			wantErr: true,
			errMsg:  "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			token, user, err := svc.Login(context.Background(), tt.username, tt.password)
			if tt.wantErr {
				if err == nil || err.Error() != tt.errMsg {
					t.Errorf("Login() error = %v, wantErr %v, errMsg %v", err, tt.wantErr, tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("Login() unexpected error: %v", err)
			}
			if user == nil || user.Username != tt.username || token == "" {
				t.Errorf("Login() user = %v, token = %v, want username %v, non-empty token", user, token, tt.username)
			}
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := ports.NewMockOrderRepositoryPort(ctrl)
	svc := NewAuthService(mockRepo)

	token, _ := auth.GenerateToken("testuser@example.com", 1)

	tests := []struct {
		name    string
		ctx     context.Context
		userID  int64
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Successful logout",
			ctx:     context.WithValue(context.Background(), "token", token),
			userID:  1,
			wantErr: false,
		},
		{
			name:    "Missing token",
			ctx:     context.Background(),
			userID:  1,
			wantErr: true,
			errMsg:  "token not found in context",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Logout(tt.ctx, tt.userID)
			if tt.wantErr {
				if err == nil || err.Error() != tt.errMsg {
					t.Errorf("Logout() error = %v, wantErr %v, errMsg %v", err, tt.wantErr, tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("Logout() unexpected error: %v", err)
			}
			_, err = auth.ValidateToken(token)
			if err == nil || !strings.Contains(err.Error(), "token is blacklisted") {
				t.Errorf("Logout() token not blacklisted, err = %v", err)
			}
		})
	}
}
