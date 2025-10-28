package order

import (
	"context"
	"fmt"

	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/ports"
)

type CancelOrderService struct {
	repo  ports.OrderRepositoryPort
	cache ports.CachePort
}

func NewCancelOrderService(repo ports.OrderRepositoryPort, cache ports.CachePort) *CancelOrderService {
	return &CancelOrderService{repo: repo, cache: cache}
}

func (s *CancelOrderService) CancelOrder(ctx context.Context, consignmentID string, userID int64) error {
	err := s.repo.CancelOrder(ctx, consignmentID, userID)
	if err != nil {
		return err
	}

	if s.cache != nil {
		err = s.cache.DeleteByPrefix(ctx, fmt.Sprintf("orders:user:%d", userID))
		if err != nil {
			fmt.Printf("Failed to invalidate cache: %v\n", err)
		}
	}
	return nil
}