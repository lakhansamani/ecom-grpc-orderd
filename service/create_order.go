package service

import (
	"context"
	"errors"

	order "github.com/lakhansamani/ecom-grpc-apis/order/v1"
	"github.com/lakhansamani/ecom-grpc-orderd/db"
)

// CreateOrder API to create a new order
// Permission: authenticated user
func (s *service) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	// Authorizer user
	userID, err := s.authorize(ctx)
	if err != nil {
		return nil, err
	}
	// Validate request
	product := req.GetProduct()
	quantity := req.GetQuantity()
	if product == "" {
		return nil, errors.New("product is required")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity should be greater than 0")
	}
	// Static Price
	price := float64(10.5)
	// Save order to database
	resOrder, err := s.DBProvider.CreateOrder(&db.Order{
		UserID:    userID,
		Product:   product,
		Quantity:  quantity,
		UnitPrice: price,
	})
	if err != nil {
		return nil, err
	}
	return &order.CreateOrderResponse{
		Order: resOrder.AsAPIOrder(),
	}, nil
}
