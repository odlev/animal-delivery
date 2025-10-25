// Package orders is a nice package
package orders

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/odlev/animal-delivery/api-gateway/internal/clients"
	pb "github.com/odlev/animal-delivery/contracts/gen/go/animaldelivery"
	"github.com/odlev/animal-delivery/contracts/gen/go/loggingpb"
	"github.com/odlev/animal-delivery/contracts/status"
	"go.opentelemetry.io/otel"
)

type Producer interface {
	Produce(msg, topic string) error
	ProduceLog(event *loggingpb.LogEvent, topic string) error
	Close()
}

type OrdersUsecase struct {
	client       *clients.OrdersClient
	producer     Producer
	createdTopic string
}

func NewOrdersUsecase(ordersClient *clients.OrdersClient, producer Producer, createdTopic string) *OrdersUsecase {
	return &OrdersUsecase{
		client:       ordersClient,
		producer:     producer,
		createdTopic: createdTopic,
	}
}

const (
	serviceName string = "api-gateway"
	levelInfo   string = "info"
)

func (u *OrdersUsecase) CreateOrder(ctx context.Context, req CreateOrderRequest) (CreateOrderResponse, error) {
	const op = "usecase.CreateOrder"

	tracer := otel.Tracer(serviceName)
	ctx, span := tracer.Start(ctx, "business layer: CreateOrder")
	defer span.End()
	pbresp, err := u.client.CreateOrder(ctx, &pb.CreateOrderRequest{
		CustomerId: req.CustomerID.String(),
		AnimalType: req.AnimalType,
		AnimalAge:  int32(req.AnimalAge),
	})

	if err != nil {
		return CreateOrderResponse{}, fmt.Errorf("%s: failed to create order: %w", op, err)
	}

	orderID, err := uuid.Parse(pbresp.OrderId)
	if err != nil {
		return CreateOrderResponse{}, fmt.Errorf("%s: failed to parse uuid: %w", op, err)
	}

	event := &loggingpb.LogEvent{
		Level:     levelInfo,
		Message:   "order created",
		Service:   serviceName,
		TraceId:   span.SpanContext().TraceID().String(),
		SpanId:    span.SpanContext().SpanID().String(),
		Timestamp: time.Now().Unix(),
		Fields:    map[string]string{"order_id": orderID.String()},
	}

	if err := u.producer.ProduceLog(event, u.createdTopic); err != nil {
		return CreateOrderResponse{}, fmt.Errorf("%s: failed to produce kafka message: %w", op, err)
	}

	return CreateOrderResponse{OrderID: orderID}, nil
}

func (u *OrdersUsecase) GetOrder(ctx context.Context, id uuid.UUID) (Order, error) {
	const op = "usecase.GetOrder"

	pbresp, err := u.client.GetOrder(ctx, &pb.GetOrderRequest{OrderId: id.String()})
	if err != nil {
		return Order{}, fmt.Errorf("%s: failed to get order: %w", op, err)
	}

	orderID, err := uuid.Parse(pbresp.OrderId)
	if err != nil {
		return Order{}, fmt.Errorf("%s: failed to parse uuid: %w", op, err)
	}

	customerID, err := uuid.Parse(pbresp.CustomerId)
	if err != nil {
		return Order{}, fmt.Errorf("%s: failed to parse uuid: %w", op, err)
	}

	updatedAT, err := time.Parse(time.RFC3339, pbresp.UpdatedAt)
	if err != nil {
		return Order{}, fmt.Errorf("%s: failed to parse time: %w", op, err)
	}

	return Order{
		ID:         orderID,
		CustomerID: customerID,
		AnimalType: pbresp.AnimalType,
		AnimalAge:  int(pbresp.AnimalAge),
		Status:     status.FromPb(pbresp.Status),
		UpdatedAt:  updatedAT,
	}, nil
}

func (u *OrdersUsecase) DeleteOrder(ctx context.Context, id uuid.UUID) (bool, error) {
	return true, nil
}
