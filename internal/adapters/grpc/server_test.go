// internal/adapters/grpc/server_test.go
package grpc

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	_ "github.com/lib/pq"
	pb "github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/grpc/proto"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/redis"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/repository"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func setupTestServer(t *testing.T) (*grpc.ClientConn, pb.OrderServiceClient) {
	lis = bufconn.Listen(bufSize)
	dsn := "host=localhost port=5432 user=postgres password=pass dbname=orderdb sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping DB: %v", err)
	}
	cache := redis.NewCache("localhost:6379", "", "",0, 5*time.Minute)
	if err := cache.Ping(context.Background()); err != nil {
		t.Fatalf("failed to connect to Redis: %v", err)
	}
	repo := repository.NewPostgresRepository(db)
	srv := NewServer(repo, cache)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(AuthInterceptor))
	pb.RegisterOrderServiceServer(grpcServer, srv)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("Server failed: %v", err)
		}
	}()

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
		return lis.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	client := pb.NewOrderServiceClient(conn)
	return conn, client
}

func TestGRPCServer(t *testing.T) {
	conn, client := setupTestServer(t)
	defer conn.Close()

	ctx := context.Background()

	// Test Signup
	t.Run("Signup_Success", func(t *testing.T) {
		resp, err := client.Signup(ctx, &pb.SignupRequest{
			Username: fmt.Sprintf("testuser%d@example.com", time.Now().UnixNano()),
			Password: "securepass",
		})
		if err != nil {
			t.Errorf("Signup failed: %v", err)
		}
		if resp.Code != 200 || resp.Type != "success" {
			t.Errorf("Signup response = %v, want code 200, type success", resp)
		}
	})

	t.Run("Signup_DuplicateUsername", func(t *testing.T) {
		username := fmt.Sprintf("testuser%d@example.com", time.Now().UnixNano())
		_, err := client.Signup(ctx, &pb.SignupRequest{Username: username, Password: "securepass"})
		if err != nil {
			t.Errorf("First signup failed: %v", err)
		}
		resp, err := client.Signup(ctx, &pb.SignupRequest{Username: username, Password: "securepass"})
		if err != nil {
			t.Errorf("Signup failed: %v", err)
		}
		if resp.Code != 400 || resp.Message != "username already exists" {
			t.Errorf("Signup response = %v, want code 400, message 'username already exists'", resp)
		}
	})

	// Test Login
	var token string
	t.Run("Login_Success", func(t *testing.T) {
		username := fmt.Sprintf("testuser%d@example.com", time.Now().UnixNano())
		_, err := client.Signup(ctx, &pb.SignupRequest{Username: username, Password: "securepass"})
		if err != nil {
			t.Errorf("Signup failed: %v", err)
		}
		resp, err := client.Login(ctx, &pb.LoginRequest{Username: username, Password: "securepass"})
		if err != nil {
			t.Errorf("Login failed: %v", err)
		}
		if resp.Code != 200 || resp.Type != "success" || resp.AccessToken == "" {
			t.Errorf("Login response = %v, want code 200, type success, non-empty token", resp)
		}
		token = resp.AccessToken
	})

	t.Run("Login_InvalidCredentials", func(t *testing.T) {
		resp, err := client.Login(ctx, &pb.LoginRequest{Username: "nonexistent@example.com", Password: "wrongpass"})
		if err != nil {
			t.Errorf("Login failed: %v", err)
		}
		if resp.Code != 400 || resp.Message != "invalid credentials" {
			t.Errorf("Login response = %v, want code 400, message 'invalid credentials'", resp)
		}
	})

	// Test CreateOrder
	t.Run("CreateOrder_Success", func(t *testing.T) {
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))
		resp, err := client.CreateOrder(ctx, &pb.CreateOrderRequest{
			StoreId:          1,
			RecipientName:    "John Doe",
			RecipientPhone:   "01712345678",
			RecipientAddress: "123 Main St",
			RecipientCity:    1,
			RecipientZone:    1,
			RecipientArea:    1,
			DeliveryType:     48,
			ItemType:         2,
			ItemQuantity:     5,
			ItemWeight:       1.5,
			AmountToCollect:  1000.0,
		})
		if err != nil {
			t.Errorf("CreateOrder failed: %v", err)
		}
		if resp.Code != 200 || resp.Type != "success" || resp.Data.ConsignmentId == "" {
			t.Errorf("CreateOrder response = %v, want code 200, type success, non-empty consignmentId", resp)
		}
	})

	t.Run("CreateOrder_Unauthorized", func(t *testing.T) {
		resp, err := client.CreateOrder(ctx, &pb.CreateOrderRequest{
			RecipientName:    "John Doe",
			RecipientPhone:   "01712345678",
			RecipientAddress: "123 Main St",
			ItemQuantity:     5,
			ItemWeight:       1.5,
			AmountToCollect:  1000.0,
		})
		if err == nil {
			t.Errorf("CreateOrder expected error, got response: %v", resp)
		}
		if s, ok := status.FromError(err); !ok || s.Code() != codes.Unauthenticated {
			t.Errorf("CreateOrder error = %v, want Unauthenticated", err)
		}
	})

	// Test ListOrders
	t.Run("ListOrders_Success", func(t *testing.T) {
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))
		resp, err := client.ListOrders(ctx, &pb.ListOrdersRequest{Limit: 10, Page: 1})
		if err != nil {
			t.Errorf("ListOrders failed: %v", err)
		}
		if resp.Code != 200 || resp.Type != "success" {
			t.Errorf("ListOrders response = %v, want code 200, type success", resp)
		}
	})

	t.Run("ListOrders_Unauthorized", func(t *testing.T) {
		resp, err := client.ListOrders(ctx, &pb.ListOrdersRequest{Limit: 10, Page: 1})
		if err == nil {
			t.Errorf("ListOrders expected error, got response: %v", resp)
		}
		if s, ok := status.FromError(err); !ok || s.Code() != codes.Unauthenticated {
			t.Errorf("ListOrders error = %v, want Unauthenticated", err)
		}
	})

	// Test ListOrders with cache
	t.Run("ListOrders_CacheHit", func(t *testing.T) {
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))
		// First call to populate cache
		_, err := client.ListOrders(ctx, &pb.ListOrdersRequest{Limit: 10, Page: 1})
		if err != nil {
			t.Errorf("ListOrders failed: %v", err)
		}
		// Second call should hit cache
		resp, err := client.ListOrders(ctx, &pb.ListOrdersRequest{Limit: 10, Page: 1})
		if err != nil {
			t.Errorf("ListOrders failed: %v", err)
		}
		if resp.Code != 200 || resp.Type != "success" {
			t.Errorf("ListOrders response = %v, want code 200, type success", resp)
		}
	})

	// Test CancelOrder
	var consignmentID string
	t.Run("CreateOrder_ForCancel", func(t *testing.T) {
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))
		resp, err := client.CreateOrder(ctx, &pb.CreateOrderRequest{
			StoreId:          1,
			RecipientName:    "John Doe",
			RecipientPhone:   "01712345678",
			RecipientAddress: "123 Main St",
			RecipientCity:    1,
			RecipientZone:    1,
			RecipientArea:    1,
			DeliveryType:     48,
			ItemType:         2,
			ItemQuantity:     5,
			ItemWeight:       1.5,
			AmountToCollect:  1000.0,
		})
		if err != nil {
			t.Errorf("CreateOrder failed: %v", err)
		}
		consignmentID = resp.Data.ConsignmentId
	})

	t.Run("CancelOrder_Success", func(t *testing.T) {
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))
		resp, err := client.CancelOrder(ctx, &pb.CancelOrderRequest{ConsignmentId: consignmentID})
		if err != nil {
			t.Errorf("CancelOrder failed: %v", err)
		}
		if resp.Code != 200 || resp.Type != "success" {
			t.Errorf("CancelOrder response = %v, want code 200, type success", resp)
		}
	})

	t.Run("CancelOrder_Unauthorized", func(t *testing.T) {
		resp, err := client.CancelOrder(ctx, &pb.CancelOrderRequest{ConsignmentId: consignmentID})
		if err == nil {
			t.Errorf("CancelOrder expected error, got response: %v", resp)
		}
		if s, ok := status.FromError(err); !ok || s.Code() != codes.Unauthenticated {
			t.Errorf("CancelOrder error = %v, want Unauthenticated", err)
		}
	})

	// Test Logout
	t.Run("Logout_Success", func(t *testing.T) {
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))
		resp, err := client.Logout(ctx, &pb.LogoutRequest{})
		if err != nil {
			t.Errorf("Logout failed: %v", err)
		}
		if resp.Code != 200 || resp.Type != "success" {
			t.Errorf("Logout response = %v, want code 200, type success", resp)
		}

		// Verify token is blacklisted
		_, err = client.ListOrders(ctx, &pb.ListOrdersRequest{Limit: 10, Page: 1})
		if err == nil {
			t.Errorf("ListOrders expected error after logout")
		}
		if s, ok := status.FromError(err); !ok || s.Code() != codes.Unauthenticated {
			t.Errorf("ListOrders error = %v, want Unauthenticated", err)
		}
	})

	t.Run("Logout_Unauthorized", func(t *testing.T) {
		resp, err := client.Logout(ctx, &pb.LogoutRequest{})
		if err == nil {
			t.Errorf("Logout expected error, got response: %v", resp)
		}
		if s, ok := status.FromError(err); !ok || s.Code() != codes.Unauthenticated {
			t.Errorf("Logout error = %v, want Unauthenticated", err)
		}
	})
}
