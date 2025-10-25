// Package transport is a nice package
package transport

import (
	"time"

	"github.com/google/uuid"
	"github.com/odlev/animal-delivery/contracts/status"
)

type Order struct {
	ID         uuid.UUID          `json:"id"`
	CustomerID uuid.UUID          `json:"customer_id"`
	AnimalType string             `json:"animal_type"`
	AnimalAge  int                `json:"animal_age"`
	Status     status.OrderStatus `json:"status"`
	UpdatedAt  time.Time          `json:"updated_at"`
}

type CreateOrderRequest struct {
	AnimalType string `json:"animal_type"`
	AnimalAge  int    `json:"animal_age"`
}

type CreateOrderResponse struct {
	OrderID string `json:"order_id"`
}
