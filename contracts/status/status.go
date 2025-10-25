// Package status provides shared order status constants and helpers.
package status

import (
	"fmt"

	pb "github.com/odlev/animal-delivery/contracts/gen/go/animaldelivery"
)

type OrderStatus string

const (
	StatusCreated    OrderStatus = "created"
	StatusInProgress OrderStatus = "in_progress"
	StatusCompleted  OrderStatus = "completed"
	StatusDeleted    OrderStatus = "deleted"
)

var statusToProto = map[OrderStatus]pb.OrderStatus{
	StatusCreated:    pb.OrderStatus_ORDER_STATUS_CREATED,
	StatusInProgress: pb.OrderStatus_ORDER_STATUS_IN_PROGRESS,
	StatusCompleted:  pb.OrderStatus_ORDER_STATUS_COMPLETED,
	StatusDeleted:    pb.OrderStatus_ORDER_STATUS_DELETED,
}

var protoToStatus = map[pb.OrderStatus]OrderStatus{
	pb.OrderStatus_ORDER_STATUS_CREATED:     StatusCreated,
	pb.OrderStatus_ORDER_STATUS_IN_PROGRESS: StatusInProgress,
	pb.OrderStatus_ORDER_STATUS_COMPLETED:   StatusCompleted,
	pb.OrderStatus_ORDER_STATUS_DELETED:     StatusDeleted,
}

func (s OrderStatus) String() string {
	return string(s)
}

func (s OrderStatus) ToPb() (pb.OrderStatus, error) {
	if v, ok := statusToProto[s]; ok {
		return v, nil
	}
	return pb.OrderStatus_ORDER_STATUS_UNSPECIFIED, fmt.Errorf("unknown order status %q", s)
}

func FromPb(p pb.OrderStatus) OrderStatus {
	if v, ok := protoToStatus[p]; ok {
		return v
	}
	return StatusCreated
}
