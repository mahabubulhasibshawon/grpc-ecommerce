package order

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/domain"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/ports"
)

type CreateOrderService struct {
	repo  ports.OrderRepositoryPort
	cache ports.CachePort
}

func NewCreateOrderService(repo ports.OrderRepositoryPort, cache ports.CachePort) *CreateOrderService {
	return &CreateOrderService{repo: repo, cache: cache}
}

func (s *CreateOrderService) CreateOrder(ctx context.Context, req *domain.Order, userID int64) (*domain.Order, error) {
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

	if s.cache != nil {
		err = s.cache.DeleteByPrefix(ctx, fmt.Sprintf("orders:user:%d", userID))
		if err != nil {
			fmt.Printf("Failed to invalidate cache: %v\n", err)
		}
	}

	return req, nil
}