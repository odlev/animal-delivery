// Package mongo is a nice package
package mongo

import (
	"time"

	"github.com/odlev/animal-delivery/orders/internal/domain"
)

type mongoOrder struct {
    ID          string        `bson:"_id"`
    Status      domain.Status `bson:"status"`
    CustomerID  string        `bson:"customer_id"`
    AnimalType  string        `bson:"animal_type"`
    AnimalAge   int32         `bson:"animal_age"`
    DeleteReason string       `bson:"delete_reason"`
    UpdatedAt   time.Time     `bson:"updated_at"`
}
