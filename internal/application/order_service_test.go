// internal/application/order_service_test.go
package application

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/domain"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/ports"
	// "github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/redis"
)

type mockCache struct {
	get    func(ctx context.Context, key string) ([]byte, error)
	set    func(ctx context.Context, key string, value interface{}) error
	delete func(ctx context.Context, prefix string) error
	ping   func(ctx context.Context) error
}

func (m *mockCache) Get(ctx context.Context, key string) ([]byte, error) {
	return m.get(ctx, key)
}

func (m *mockCache) Set(ctx context.Context, key string, value interface{}) error {
	return m.set(ctx, key, value)
}

func (m *mockCache) DeleteByPrefix(ctx context.Context, prefix string) error {
	return m.delete(ctx, prefix)
}

func (m *mockCache) Ping(ctx context.Context) error {
	return m.ping(ctx)
}

func TestOrderService_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := ports.NewMockOrderRepositoryPort(ctrl)
	mockCache := &mockCache{
		delete: func(ctx context.Context, prefix string) error { return nil },
		ping:   func(ctx context.Context) error { return nil },
	}
	svc := NewOrderService(mockRepo, mockCache)

	validOrder := &domain.Order{
		RecipientName:    "John Doe",
		RecipientPhone:   "01712345678",
		RecipientAddress: "123 Main St",
		ItemQuantity:     5,
		ItemWeight:       1.5,
		AmountToCollect:  1000.0,
		RecipientCity:    1,
	}

	tests := []struct {
		name      string
		order     *domain.Order
		userID    int64
		mockSetup func()
		wantErr   bool
		errMsg    string
	}{
		{
			name:   "Successful order creation",
			order:  validOrder,
			userID: 1,
			mockSetup: func() {
				mockRepo.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(nil)
				mockCache.delete = func(ctx context.Context, prefix string) error { return nil }
			},
			wantErr: false,
		},
		{
			name: "Missing required fields",
			order: &domain.Order{
				RecipientName: "",
				RecipientPhone: "01712345678",
				ItemQuantity:   5,
				ItemWeight:     1.5,
				AmountToCollect: 1000.0,
			},
			userID:    1,
			mockSetup: func() {},
			wantErr:   true,
			errMsg:    "missing required fields",
		},
		{
			name: "Invalid phone number",
			order: &domain.Order{
				RecipientName:    "John Doe",
				RecipientPhone:   "123",
				RecipientAddress: "123 Main St",
				ItemQuantity:     5,
				ItemWeight:       1.5,
				AmountToCollect:  1000.0,
			},
			userID:    1,
			mockSetup: func() {},
			wantErr:   true,
			errMsg:    "invalid phone number",
		},
		{
			name:   "Cache deletion error",
			order:  validOrder,
			userID: 1,
			mockSetup: func() {
				mockRepo.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(nil)
				mockCache.delete = func(ctx context.Context, prefix string) error { return errors.New("cache error") }
			},
			wantErr: false, // Cache error doesn't fail the operation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			created, err := svc.CreateOrder(context.Background(), tt.order, tt.userID)
			if tt.wantErr {
				if err == nil || err.Error() != tt.errMsg {
					t.Errorf("CreateOrder() error = %v, wantErr %v, errMsg %v", err, tt.wantErr, tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("CreateOrder() unexpected error: %v", err)
			}
			if created == nil || created.ConsignmentID == "" {
				t.Errorf("CreateOrder() created = %v, want non-nil order with consignment ID", created)
			}
		})
	}
}

func TestOrderService_ListOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := ports.NewMockOrderRepositoryPort(ctrl)
	mockCache := &mockCache{
		ping: func(ctx context.Context) error { return nil },
	}
	svc := NewOrderService(mockRepo, mockCache)

	orders := []*domain.Order{
		{
			ConsignmentID: "DA251021BNWWN123",
			CreatedAt:     time.Now(),
			RecipientName: "John Doe",
			UserID:        1,
		},
	}
	total := int64(1)
	cacheData := struct {
		Orders []*domain.Order
		Total  int64
	}{Orders: orders, Total: total}
	cacheBytes, _ := json.Marshal(cacheData)

	tests := []struct {
		name      string
		userID    int64
		limit     int64
		page      int64
		mockSetup func()
		wantErr   bool
	}{
		{
			name:   "Cache hit",
			userID: 1,
			limit:  10,
			page:   1,
			mockSetup: func() {
				mockCache.get = func(ctx context.Context, key string) ([]byte, error) { return cacheBytes, nil }
			},
			wantErr: false,
		},
		{
			name:   "Cache miss, successful DB query",
			userID: 1,
			limit:  10,
			page:   1,
			mockSetup: func() {
				mockCache.get = func(ctx context.Context, key string) ([]byte, error) { return nil, errors.New("cache miss") }
				mockRepo.EXPECT().ListOrders(gomock.Any(), int64(1), int64(10), int64(1)).Return(orders, total, nil)
				mockCache.set = func(ctx context.Context, key string, value interface{}) error { return nil }
			},
			wantErr: false,
		},
		{
			name:   "Repository error",
			userID: 1,
			limit:  10,
			page:   1,
			mockSetup: func() {
				mockCache.get = func(ctx context.Context, key string) ([]byte, error) { return nil, errors.New("cache miss") }
				mockRepo.EXPECT().ListOrders(gomock.Any(), int64(1), int64(10), int64(1)).Return(nil, int64(0), errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name:   "Cache set error",
			userID: 1,
			limit:  10,
			page:   1,
			mockSetup: func() {
				mockCache.get = func(ctx context.Context, key string) ([]byte, error) { return nil, errors.New("cache miss") }
				mockRepo.EXPECT().ListOrders(gomock.Any(), int64(1), int64(10), int64(1)).Return(orders, total, nil)
				mockCache.set = func(ctx context.Context, key string, value interface{}) error { return errors.New("cache set error") }
			},
			wantErr: false, // Cache error doesn't fail the operation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			result, total, err := svc.ListOrders(context.Background(), tt.userID, tt.limit, tt.page)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ListOrders() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("ListOrders() unexpected error: %v", err)
			}
			if len(result) != len(orders) || total != 1 {
				t.Errorf("ListOrders() result = %v, total = %v, want %v, 1", result, total, orders)
			}
		})
	}
}

func TestOrderService_CancelOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := ports.NewMockOrderRepositoryPort(ctrl)
	mockCache := &mockCache{
		delete: func(ctx context.Context, prefix string) error { return nil },
		ping:   func(ctx context.Context) error { return nil },
	}
	svc := NewOrderService(mockRepo, mockCache)

	tests := []struct {
		name          string
		consignmentID string
		userID        int64
		mockSetup     func()
		wantErr       bool
		errMsg        string
	}{
		{
			name:          "Successful cancel",
			consignmentID: "DA251021BNWWN123",
			userID:        1,
			mockSetup: func() {
				mockRepo.EXPECT().CancelOrder(gomock.Any(), "DA251021BNWWN123", int64(1)).Return(nil)
				mockCache.delete = func(ctx context.Context, prefix string) error { return nil }
			},
			wantErr: false,
		},
		{
			name:          "Order not found",
			consignmentID: "DA251021BNWWN123",
			userID:        1,
			mockSetup: func() {
				mockRepo.EXPECT().CancelOrder(gomock.Any(), "DA251021BNWWN123", int64(1)).Return(errors.New("order not found"))
				mockCache.delete = func(ctx context.Context, prefix string) error { return nil }
			},
			wantErr: true,
			errMsg:  "order not found",
		},
		{
			name:          "Cache deletion error",
			consignmentID: "DA251021BNWWN123",
			userID:        1,
			mockSetup: func() {
				mockRepo.EXPECT().CancelOrder(gomock.Any(), "DA251021BNWWN123", int64(1)).Return(nil)
				mockCache.delete = func(ctx context.Context, prefix string) error { return errors.New("cache error") }
			},
			wantErr: false, // Cache error doesn't fail the operation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := svc.CancelOrder(context.Background(), tt.consignmentID, tt.userID)
			if tt.wantErr {
				if err == nil || err.Error() != tt.errMsg {
					t.Errorf("CancelOrder() error = %v, wantErr %v, errMsg %v", err, tt.wantErr, tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("CancelOrder() unexpected error: %v", err)
			}
		})
	}
}