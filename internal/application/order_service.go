// internal/application/order_service.go
package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	// "github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/redis"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/domain"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/ports"
)

type OrderService struct {
	repo  ports.OrderRepositoryPort
	cache ports.CachePort
}

func NewOrderService(repo ports.OrderRepositoryPort, cache ports.CachePort) *OrderService {
	return &OrderService{repo: repo, cache: cache}
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

	// Invalidate cache for this user
	if s.cache != nil {
		err = s.cache.DeleteByPrefix(ctx, fmt.Sprintf("orders:user:%d", userID))
		if err != nil {
			fmt.Printf("Failed to invalidate cache: %v\n", err)
		}
	}

	return req, nil
}

func (s *OrderService) ListOrders(ctx context.Context, userID int64, limit, page int64) ([]*domain.Order, int64, error) {
	if limit < 1 {
		limit = 10
	}
	if page < 1 {
		page = 1
	}

	// check cache
	if s.cache != nil {
		cacheKey := fmt.Sprintf("orders:user:%d:page:%d:limit:%d", userID, page, limit)
		cached, err := s.cache.Get(ctx, cacheKey)
		if err != nil {
			// var orders []*domain.Order
			// var total int64
			type cachedData struct {
				Orders []*domain.Order
				Total  int64
			}
			var data cachedData
			if err := json.Unmarshal(cached, &data); err == nil {
				return data.Orders, data.Total, nil
			}
		}
	}
	// Cache miss, query database
	orders, total, err := s.repo.ListOrders(ctx, userID, limit, page)
	if err != nil {
		return nil, 0, err
	}

	// Cache result
	if s.cache != nil {
		cacheKey := fmt.Sprintf("orders:user:%d:page:%d:limit:%d", userID, page, limit)
		data := struct {
			Orders []*domain.Order
			Total  int64
		}{Orders: orders, Total: total}
		if err := s.cache.Set(ctx, cacheKey, data); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Failed to cache orders: %v\n", err)
		}
	}

	// return s.repo.ListOrders(ctx, userID, limit, page)
	return orders, total, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, consignmentID string, userID int64) error {
	err := s.repo.CancelOrder(ctx, consignmentID, userID)
	if err != nil {
		return err
	}

	// invalidate cache for this user
	if s.cache != nil {
		err = s.cache.DeleteByPrefix(ctx, fmt.Sprint("orders:user:%d", userID))
		if err != nil {
			fmt.Printf("failed to invalidate cache %v", err)
		}
	}
	return nil
}
