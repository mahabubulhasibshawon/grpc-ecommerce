package order

import (
	"context"
	"math"
	"time"

	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/pkg/auth"
	pb "github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *OrderHandler) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}

	orders, total, err := h.server.OrderService.ListOrders(ctx, userID, req.Limit, req.Page)
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
