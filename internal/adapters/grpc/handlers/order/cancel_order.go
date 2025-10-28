package order

import (
	"context"
	
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/pkg/auth"
	pb "github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *OrderHandler) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.CancelOrderResponse, error) {
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}

	err = h.cancelOrderService.CancelOrder(ctx, req.ConsignmentId, userID)
	if err != nil {
		return &pb.CancelOrderResponse{Message: err.Error(), Type: "error", Code: 400}, nil
	}
	return &pb.CancelOrderResponse{Message: "Order Cancelled Successfully", Type: "success", Code: 200}, nil
}
