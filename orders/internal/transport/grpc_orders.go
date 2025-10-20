// Package transport is a nice package
package transport

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/odlev/animal-delivery/contracts/gen/go/animaldelivery"
	"github.com/odlev/animal-delivery/orders/internal/domain"
	"github.com/odlev/animal-delivery/orders/internal/usecase"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Orders interface {
	CreateOrder(ctx context.Context, order domain.Order) (uuid.UUID, error)
	GetOrder(ctx context.Context, strUUID string) (domain.Order, error)
	DeleteOrder(ctx context.Context, strUUID string) error
}

type server struct {
	pb.UnimplementedOrderServiceServer
	Orders
	log *logrus.Logger
}

func NewServer(orders Orders, log *logrus.Logger) pb.OrderServiceServer {
	return &server{
		UnimplementedOrderServiceServer: pb.UnimplementedOrderServiceServer{},
		Orders:                          orders,
		log:                             log,
	}
}

func Register(srv *grpc.Server, orders Orders, log *logrus.Logger) {
	pb.RegisterOrderServiceServer(srv, NewServer(orders, log))
}

func (s *server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	customerID, err := uuid.Parse(req.GetCustomerId())
	if err != nil {
		s.log.Error("failed to parse uuid: ")
		return nil, status.Error(codes.InvalidArgument, "field customer id: invalid uuid")
	}
	id, err := s.Orders.CreateOrder(ctx, domain.Order{
		CustomerID: customerID,
		Status:     domain.StatusCreated,
		AnimalType: req.GetAnimalType(),
		AnimalAge:  req.GetAnimalAge(),
	})
	if err != nil {
		s.log.Error("failed to create order: ", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	s.log.Infof("order created, id: %s", id)
	return &pb.CreateOrderResponse{OrderId: id.String()}, nil
}

func (s *server) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order, err := s.Orders.GetOrder(ctx, req.GetOrderId())
	if err != nil {
		s.log.Error("failed to get order: ", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.GetOrderResponse{
		OrderId: order.OrderID.String(),
		Status: toPbStatus(order.Status),
		CustomerId: order.CustomerID.String(),
		AnimalType: order.AnimalType,
		AnimalAge: order.AnimalAge,
		UpdatedAt: order.UpdatedAt.String(),
	}, nil
}

func (s *server) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	if err := s.Orders.DeleteOrder(ctx, req.GetOrderId()); err != nil {
		s.log.Info("failed to delete order: ", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.DeleteOrderResponse{Ok: true}, nil
}

func toDomainStatus(s pb.OrderStatus) (domain.Status, error) {
	// выбор между мапой и свитчей пал на свитч. скорость одинакова, а типов не слишком много для мапы
	switch s {
	case pb.OrderStatus_ORDER_STATUS_CREATED:
		return domain.StatusCreated, nil
	case pb.OrderStatus_ORDER_STATUS_IN_PROGRESS:
		return domain.StatusInProgress, nil
	case pb.OrderStatus_ORDER_STATUS_COMPLETED:
		return domain.StatusCompleted, nil
	case pb.OrderStatus_ORDER_STATUS_DELETED:
		return domain.StatusDeleted, nil
	default:
		return "", usecase.ErrInvalidStatus
	}
}

func toPbStatus(s domain.Status) pb.OrderStatus {
	switch s {
	case domain.StatusCreated:
		return pb.OrderStatus_ORDER_STATUS_CREATED
	case domain.StatusInProgress:
		return pb.OrderStatus_ORDER_STATUS_IN_PROGRESS
	case domain.StatusCompleted:
		return pb.OrderStatus_ORDER_STATUS_COMPLETED
	case domain.StatusDeleted:
		return pb.OrderStatus_ORDER_STATUS_DELETED
	default:
		return pb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}
