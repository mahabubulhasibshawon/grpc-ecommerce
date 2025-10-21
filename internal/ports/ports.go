// internal/ports/ports.go
package ports

import (
	"context"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/domain"
)

type AuthPort interface {
	Signup(ctx context.Context, username, password string) (*domain.User, error)
	Login(ctx context.Context, username, password string) (string, *domain.User, error)
	Logout(ctx context.Context, userID int64) error
}

type OrderRepositoryPort interface {
	CreateUser(ctx context.Context, username, password string) (*domain.User, error)
	FindUserByUsername(ctx context.Context, username string) (*domain.User, error)
	CreateOrder(ctx context.Context, order *domain.Order) error
	ListOrders(ctx context.Context, userID int64, limit, page int64) ([]*domain.Order, int64, error)
	CancelOrder(ctx context.Context, consignmentID string, userID int64) error
}