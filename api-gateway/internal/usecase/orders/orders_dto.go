// Package orders is a nice package
package orders

import (
	"time"

	"github.com/google/uuid"
	"github.com/odlev/animal-delivery/contracts/status"
)

type CreateOrderRequest struct {
	CustomerID uuid.UUID
	AnimalType string
	AnimalAge  int
}

type CreateOrderResponse struct {
	OrderID uuid.UUID
}

type Order struct {
	ID         uuid.UUID
	CustomerID uuid.UUID
	AnimalType string
	AnimalAge  int
	Status     status.OrderStatus
	UpdatedAt  time.Time
}
