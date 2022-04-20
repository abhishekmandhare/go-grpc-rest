package api

import (
	"context"

	orders "github.com/abhishekmandhare/go-grpc-rest/orders-api/proto"
)

type OrdersAPI struct {
	orders.UnimplementedOrdersAPIServer
}

// NewOrdersAPI creates orders service with given parameters
func NewOrdersAPI() OrdersAPI {
	return OrdersAPI{}
}

// CreateOrder creates order using given parameters.
func (ordersAPI OrdersAPI) CreateOrder(context.Context, *orders.CreateOrderRequest) (*orders.CreateOrderResponse, error) {
	return &orders.CreateOrderResponse{Num: 22}, nil
}
