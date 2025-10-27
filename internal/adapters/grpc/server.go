// internal/adapters/grpc/server.go
package grpc

import (
	pb "github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/grpc/proto"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/redis"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/application"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/ports"
)

type Server struct {
	pb.UnimplementedOrderServiceServer
	AuthService  *application.AuthService
	OrderService *application.OrderService
}

func NewServer(repo ports.OrderRepositoryPort, cache *redis.Cache) *Server {
	return &Server{
		AuthService:  application.NewAuthService(repo),
		OrderService: application.NewOrderService(repo, cache),
	}
}

