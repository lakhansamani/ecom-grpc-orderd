package service

import (
	"context"
	"errors"

	order "github.com/lakhansamani/ecom-grpc-apis/order/v1"
)

// GetOrder API to get order details
// Permission: authenticated user who created the order
func (s *service) GetOrder(ctx context.Context, req *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	ctx, span := s.trace.Start(ctx, "GetOrder")
	defer span.End()
	// Authorizer user
	userID, err := s.authorize(ctx)
	if err != nil {
		return nil, err
	}
	// Get order from database
	orderID := req.GetId()
	if orderID == "" {
		return nil, errors.New("order id is required")
	}
	resOrder, err := s.DBProvider.GetOrderById(ctx, orderID)
	if err != nil {
		return nil, err
	}
	// Check if user is authorized to get the order
	if resOrder.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	return &order.GetOrderResponse{
		Order: resOrder.AsAPIOrder(),
	}, nil
}
