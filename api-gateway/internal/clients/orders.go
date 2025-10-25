// Package clients is a nice package
package clients

import (
	"context"

	pb "github.com/odlev/animal-delivery/contracts/gen/go/animaldelivery"
	"google.golang.org/grpc"
)

type OrdersClient struct {
	grpc pb.OrderServiceClient
}

func NewOrdersClient(conn *grpc.ClientConn) *OrdersClient {
	return &OrdersClient{grpc: pb.NewOrderServiceClient(conn)}
}

func (c *OrdersClient) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	return c.grpc.CreateOrder(ctx, req)
}

func (c *OrdersClient) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	return c.grpc.GetOrder(ctx, req)
}

func (c *OrdersClient) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	return c.grpc.DeleteOrder(ctx, req)
}
