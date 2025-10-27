package order

import "github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/grpc"

type OrderHandler struct {
	server *grpc.Server
}

func NewAuthHandler(s *grpc.Server) *OrderHandler {
	return &OrderHandler{server: s}
}
