// Package postgres is a nice package
package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/odlev/animal-delivery/orders/internal/domain"
)

type orderRow struct {
	ID           uuid.UUID
	Status       domain.Status
	CustomerID   uuid.UUID
	AnimalType   string
	AnimalAge    int32
	DeleteReason string
	UpdatedAt    time.Time // заменив на time.Time, если в domain.Time нет особой логики
}
