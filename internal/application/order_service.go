// internal/application/order_service.go
package application

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/domain"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/ports"
)

type OrderService struct {
	repo ports.OrderRepositoryPort
}

func NewOrderService(repo ports.OrderRepositoryPort) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *domain.Order, userID int64) (*domain.Order, error) {
	if req.RecipientName == "" || req.RecipientPhone == "" || req.RecipientAddress == "" || req.ItemQuantity == 0 || req.ItemWeight == 0 || req.AmountToCollect == 0 {
		return nil, errors.New("missing required fields")
	}
	phoneRegex := regexp.MustCompile(`^(01)[3-9]{1}[0-9]{8}$`)
	if !phoneRegex.MatchString(req.RecipientPhone) {
		return nil, errors.New("invalid phone number")
	}
	if req.RecipientAddress == "" {
		req.RecipientAddress = "banani, gulshan 2, dhaka, bangladesh"
	}

	baseFee := 60.0
	if req.RecipientCity != 1 {
		baseFee = 100.0
	}
	deliveryFee := baseFee
	if req.ItemWeight > 0.5 && req.ItemWeight <= 1 {
		deliveryFee = 70.0
	} else if req.ItemWeight > 1 {
		extraKg := req.ItemWeight - 1
		deliveryFee = baseFee + 10 + (extraKg * 15)
	}
	req.DeliveryFee = deliveryFee
	req.DeliveryCharge = deliveryFee
	req.CODFee = req.AmountToCollect * 0.01
	req.TotalFee = req.DeliveryFee + req.CODFee
	req.CODAmount = req.AmountToCollect
	req.OrderAmount = req.AmountToCollect

	req.StoreName = "Default Store"
	req.StoreContactPhone = "123456789"
	req.OrderType = "Delivery"
	req.OrderTypeID = 1
	req.PromoDiscount = 0
	req.Discount = 0

	req.ConsignmentID = fmt.Sprintf("DA%vBNWWN%d", time.Now().Format("060102"), time.Now().UnixNano()%1000)
	req.CreatedAt = time.Now()
	req.Status = "Pending"
	req.UserID = userID

	err := s.repo.CreateOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (s *OrderService) ListOrders(ctx context.Context, userID int64, limit, page int64) ([]*domain.Order, int64, error) {
	return s.repo.ListOrders(ctx, userID, limit, page)
}

func (s *OrderService) CancelOrder(ctx context.Context, consignmentID string, userID int64) error {
	return s.repo.CancelOrder(ctx, consignmentID, userID)
}