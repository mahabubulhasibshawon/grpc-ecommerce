package order

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/domain"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/ports"
)

type ListOrdersService struct {
	repo  ports.OrderRepositoryPort
	cache ports.CachePort
}

func NewListOrdersService(repo ports.OrderRepositoryPort, cache ports.CachePort) *ListOrdersService {
	return &ListOrdersService{repo: repo, cache: cache}
}

func (s *ListOrdersService) ListOrders(ctx context.Context, userID int64, limit, page int64) ([]*domain.Order, int64, error) {
	if limit < 1 {
		limit = 10
	}
	if page < 1 {
		page = 1
	}

	if s.cache != nil {
		cacheKey := fmt.Sprintf("orders:user:%d:page:%d:limit:%d", userID, page, limit)
		cached, err := s.cache.Get(ctx, cacheKey)
		if err == nil {
			type cachedData struct {
				Orders []*domain.Order
				Total  int64
			}
			var data cachedData
			if err := json.Unmarshal(cached, &data); err == nil {
				log.Printf("[cache hit] key=%s", cacheKey)
				return data.Orders, data.Total, nil
			}
			log.Printf("[cache unmarshal error] key=%s err=%v", cacheKey, err)
		} else {
			log.Printf("[cache miss] key=%s err=%v", cacheKey, err)
		}
	}

	orders, total, err := s.repo.ListOrders(ctx, userID, limit, page)
	if err != nil {
		return nil, 0, err
	}

	if s.cache != nil {
		cacheKey := fmt.Sprintf("orders:user:%d:page:%d:limit:%d", userID, page, limit)
		data := struct {
			Orders []*domain.Order
			Total  int64
		}{Orders: orders, Total: total}
		if err := s.cache.Set(ctx, cacheKey, data); err != nil {
			fmt.Printf("Failed to cache orders: %v\n", err)
		}
	}

	return orders, total, nil
}
