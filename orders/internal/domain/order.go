// Package domain is a nice package
package domain

import (
	"time"

	"github.com/google/uuid"
)

type Status string

func (s Status) Str() string {
	return string(s)
}

type Order struct {
	OrderID      uuid.UUID
	Status       Status
	CustomerID   uuid.UUID
	AnimalType   string
	AnimalAge    int32
	DeleteReason string
	UpdatedAt    time.Time
}

const (
	StatusCreated    Status = "created"
	StatusInProgress Status = "in progress"
	StatusCompleted Status = "completed"
	StatusDeleted Status = "deleted"
)
