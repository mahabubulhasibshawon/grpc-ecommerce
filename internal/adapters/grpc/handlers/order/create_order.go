package order

import (
	"context"

	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/pkg/auth"
	pb "github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/grpc/proto"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
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

	created, err := h.server.OrderService.CreateOrder(ctx, order, userID)
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
