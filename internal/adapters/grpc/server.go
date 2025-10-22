// internal/adapters/grpc/server.go
package grpc

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/grpc/proto"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/redis"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/application"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/domain"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/ports"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/pkg/auth"
)

type Server struct {
	pb.UnimplementedOrderServiceServer
	authService  *application.AuthService
	orderService *application.OrderService
}

func NewServer(repo ports.OrderRepositoryPort, cache *redis.Cache) *Server {
	return &Server{
		authService:  application.NewAuthService(repo),
		orderService: application.NewOrderService(repo, cache),
	}
}

func (s *Server) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	_ , err := s.authService.Signup(ctx, req.Username, req.Password)
	if err != nil {
		return &pb.SignupResponse{Message: err.Error(), Type: "error", Code: 400}, nil
	}
	return &pb.SignupResponse{
		Message: "User registered successfully",
		Type:    "success",
		Code:    200,
	}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, user, err := s.authService.Login(ctx, req.Username, req.Password)
	if err != nil {
		return &pb.LoginResponse{Message: "Invalid credentials", Type: "error", Code: 400}, nil
	}
	token, err = auth.GenerateToken(req.Username, user.ID)
	if err != nil {
		return &pb.LoginResponse{Message: err.Error(), Type: "error", Code: 500}, nil
	}
	return &pb.LoginResponse{
		TokenType:    "Bearer",
		ExpiresIn:    432000,
		AccessToken:  token,
		RefreshToken: "dummy-refresh",
		Message:      "Logged in",
		Type:         "success",
		Code:         200,
	}, nil
}

func (s *Server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}

	order := &domain.Order{
		StoreID:          req.StoreId,
		MerchantOrderID:  req.MerchantOrderId,
		RecipientName:    req.RecipientName,
		RecipientPhone:   req.RecipientPhone,
		RecipientAddress: req.RecipientAddress,
		RecipientCity:    req.RecipientCity,
		RecipientZone:    req.RecipientZone,
		RecipientArea:    req.RecipientArea,
		DeliveryType:     req.DeliveryType,
		ItemType:         req.ItemType,
		Instruction:      req.SpecialInstruction,
		ItemQuantity:     req.ItemQuantity,
		ItemWeight:       req.ItemWeight,
		AmountToCollect:  req.AmountToCollect,
		Description:      req.ItemDescription,
	}

	created, err := s.orderService.CreateOrder(ctx, order, userID)
	if err != nil {
		return &pb.CreateOrderResponse{Message: err.Error(), Type: "error", Code: 422}, nil
	}
	return &pb.CreateOrderResponse{
		Message: "Order Created Successfully",
		Type:    "success",
		Code:    200,
		Data: &pb.OrderData{
			ConsignmentId:   created.ConsignmentID,
			MerchantOrderId: created.MerchantOrderID,
			OrderStatus:     created.Status,
			DeliveryFee:     created.DeliveryFee,
		},
	}, nil
}

func (s *Server) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}

	orders, total, err := s.orderService.ListOrders(ctx, userID, req.Limit, req.Page)
	if err != nil {
		return &pb.ListOrdersResponse{Message: err.Error(), Type: "error", Code: 400}, nil
	}

	var pbOrders []*pb.Order
	for _, o := range orders {
		pbOrders = append(pbOrders, &pb.Order{
			OrderConsignmentId: o.ConsignmentID,
			OrderCreatedAt:     o.CreatedAt.Format(time.RFC3339),
			OrderDescription:   o.Description,
			MerchantOrderId:    o.MerchantOrderID,
			RecipientName:      o.RecipientName,
			RecipientAddress:   o.RecipientAddress,
			RecipientPhone:     o.RecipientPhone,
			OrderAmount:        o.OrderAmount,
			TotalFee:           o.TotalFee,
			Instruction:        o.Instruction,
			OrderTypeId:        o.OrderTypeID,
			CodFee:             o.CODFee,
			PromoDiscount:      o.PromoDiscount,
			Discount:           o.Discount,
			DeliveryFee:        o.DeliveryFee,
			OrderStatus:        o.Status,
			OrderType:          o.OrderType,
			ItemType:           o.ItemType,
			StoreName:          o.StoreName,
			StoreContactPhone:  o.StoreContactPhone,
			CodAmount:          o.CODAmount,
			DeliveryCharge:     o.DeliveryCharge,
			StoreId:            o.StoreID,
			RecipientCity:      o.RecipientCity,
			RecipientZone:      o.RecipientZone,
			RecipientArea:      o.RecipientArea,
			DeliveryType:       o.DeliveryType,
			ItemQuantity:       o.ItemQuantity,
			ItemWeight:         o.ItemWeight,
			AmountToCollect:    o.AmountToCollect,
		})
	}

	lastPage := int64(math.Ceil(float64(total) / float64(req.Limit)))
	return &pb.ListOrdersResponse{
		Message: "Orders successfully fetched.",
		Type:    "success",
		Code:    200,
		Data: &pb.OrdersData{
			Orders:      pbOrders,
			Total:       total,
			CurrentPage: req.Page,
			PerPage:     req.Limit,
			TotalInPage: int64(len(orders)),
			LastPage:    lastPage,
		},
	}, nil
}

func (s *Server) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.CancelOrderResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}

	err = s.orderService.CancelOrder(ctx, req.ConsignmentId, userID)
	if err != nil {
		return &pb.CancelOrderResponse{Message: err.Error(), Type: "error", Code: 400}, nil
	}
	return &pb.CancelOrderResponse{Message: "Order Cancelled Successfully", Type: "success", Code: 200}, nil
}

func (s *Server) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}

	err = s.authService.Logout(ctx, userID)
	if err != nil {
		return &pb.LogoutResponse{Message: err.Error(), Type: "error", Code: 400}, nil
	}
	return &pb.LogoutResponse{Message: "Successfully logged out", Type: "success", Code: 200}, nil
}

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == "/order.OrderService/Login" || info.FullMethod == "/order.OrderService/Signup" {
		return handler(ctx, req)
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}
	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing authorization")
	}
	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	claims, err := auth.ValidateToken(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	ctx = context.WithValue(ctx, "userID", claims.UserID)
	return handler(ctx, req)
}

func getUserIDFromContext(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, errors.New("missing metadata")
	}
	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return 0, errors.New("missing authorization")
	}
	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	claims, err := auth.ValidateToken(token)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
