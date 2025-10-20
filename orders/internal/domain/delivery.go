// Package domain is a nice package
package domain

import (
	"time"

	"github.com/google/uuid"
)

type Delivery struct {
	DeliveryID      uuid.UUID
	OrderID         uuid.UUID
	OrderStatus     string
	Description     string
	DeliveryAddress string
	PickupAddress   string
	UpdatedAt       time.Time
}
