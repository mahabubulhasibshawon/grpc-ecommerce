package grpc

import (
	"context"

	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/grpc/handlers/auth"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/grpc/handlers/order"
	pb "github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/grpc/proto"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/ports"
)

type Server struct {
	pb.UnimplementedOrderServiceServer
	authHandler  *auth.AuthHandler
	orderHandler *order.OrderHandler
}

func NewServer(repo ports.OrderRepositoryPort, cache ports.CachePort) *Server {
	return &Server{
		authHandler:  auth.NewAuthHandler(repo),
		orderHandler: order.NewOrderHandler(repo, cache),
	}
}

func (s *Server) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	return s.authHandler.Signup(ctx, req)
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return s.authHandler.Login(ctx, req)
}

func (s *Server) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	return s.authHandler.Logout(ctx, req)
}

func (s *Server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	return s.orderHandler.CreateOrder(ctx, req)
}

func (s *Server) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	return s.orderHandler.ListOrders(ctx, req)
}

func (s *Server) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.CancelOrderResponse, error) {
	return s.orderHandler.CancelOrder(ctx, req)
}