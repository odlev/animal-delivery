// Package transport is a nice package
package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	usorders "github.com/odlev/animal-delivery/api-gateway/internal/usecase/orders"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

type OrdersManager interface {
	CreateOrder(ctx context.Context, req usorders.CreateOrderRequest) (usorders.CreateOrderResponse, error)
	GetOrder(ctx context.Context, id uuid.UUID) (usorders.Order, error)
	DeleteOrder(ctx context.Context, id uuid.UUID) (bool, error)
}

type Handler struct {
	OrdersManager
	log *zerolog.Logger
}

func NewHandler(orders OrdersManager, logger *zerolog.Logger) *Handler {
	return &Handler{OrdersManager: orders, log: logger}
}

const (
	serviceName string = "api-gateway"
)

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler layer: CreateOrder")

	customerid, ok := ctx.Value(customerID).(string)
	if !ok {
		http.Error(w, "failed asserting customer_id to string", http.StatusBadRequest)
	}
	h.log.Info().Str("customer_id", customerid).Msg("creating order")
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error().Err(err).Msg("failed to decode requst body")
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}
	usreq, err := req.toUSLayer()
	if err != nil {
		http.Error(w, "error uuid parse, invalid uuid", http.StatusBadRequest)
	}
	customerID, err := uuid.Parse(customerid)
	if err != nil {
		http.Error(w, "failed to parse customer uuid", http.StatusBadRequest)
	}
	usreq.CustomerID = customerID
	created, err := h.OrdersManager.CreateOrder(ctx, usreq)
	if err != nil {
		h.log.Error().
			Err(err).
			Msg("failed to create order")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = writeJSON(w, http.StatusCreated, created)
	if err != nil {
		h.log.Error().
			Err(err).
			Msg("failed to write json")
		http.Error(w, "failed to write json", http.StatusInternalServerError)
	}
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	h.log.Info().Msg("getting order")
	orderID := r.PathValue("id")
	id, err := uuid.Parse(orderID)
	if err != nil {
		h.log.Error().
			Str("order_id", orderID).
			Err(err).
			Msg("failed to parse uuid")
		http.Error(w, "failed to parse uuid", http.StatusBadRequest)
		return
	}
	order, err := h.OrdersManager.GetOrder(r.Context(), id)
	if err != nil {
		h.log.Error().
			Str("order_id", orderID).
			Err(err).
			Msg("failed to get order")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = writeJSON(w, http.StatusFound, order)
	if err != nil {
		h.log.Error().
			Str("order_id", orderID).
			Err(err).
			Msg("failed to write json")
		http.Error(w, "failed to write json", http.StatusInternalServerError)
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		return fmt.Errorf("failed to encode: %w", err)
	}
	return nil
}

func (req *CreateOrderRequest) toUSLayer() (usorders.CreateOrderRequest, error) {
	return usorders.CreateOrderRequest{
		AnimalType: req.AnimalType,
		AnimalAge:  req.AnimalAge,
	}, nil
}
