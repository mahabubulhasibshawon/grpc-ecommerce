package order

import (
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/application/order"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/ports"
)

type OrderHandler struct {
	createOrderService *order.CreateOrderService
	listOrdersService  *order.ListOrdersService
	cancelOrderService *order.CancelOrderService
}

func NewOrderHandler(repo ports.OrderRepositoryPort, cache ports.CachePort) *OrderHandler {
	return &OrderHandler{
		createOrderService: order.NewCreateOrderService(repo, cache),
		listOrdersService:  order.NewListOrdersService(repo, cache),
		cancelOrderService: order.NewCancelOrderService(repo, cache),
	}
}
