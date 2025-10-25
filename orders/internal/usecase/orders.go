// Package usecase is a nice package
package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/odlev/animal-delivery/orders/internal/domain"
	"github.com/odlev/animal-delivery/orders/internal/repo"
)

type UseCase struct {
	OrdersManager
}

type OrdersManager interface {
	CreateOrder(ctx context.Context, order domain.Order) error
	GetOrder(ctx context.Context, uuid uuid.UUID) (domain.Order, error)
	DeleteOrder(ctx context.Context, uuid uuid.UUID) error
	Close()
}

func New(orders OrdersManager) *UseCase {
	return &UseCase{
		OrdersManager: orders,
	}
}

func (u *UseCase) CreateOrder(ctx context.Context, order domain.Order) (uuid.UUID, error) {
	const op = "usecase.CreateOrder"

	id := uuid.New()
	order.OrderID = id
	if err := u.OrdersManager.CreateOrder(ctx, order); err != nil {
		OrdersCreatedCounter.WithLabelValues("error", "ru").Inc()
		return uuid.Nil, fmt.Errorf("%s: failed to create order: %s", op, err)
	}
	OrdersCreatedCounter.WithLabelValues("success", "ru").Inc()

	return id, nil
}

func (u *UseCase) GetOrder(ctx context.Context, strUUID string) (domain.Order, error) {
	const op = "usecase.GetOrder"

	id, err := uuid.Parse(strUUID)
	if err != nil {
		return domain.Order{}, ErrInvalidUUID
	}

	var order domain.Order
	order, err = u.OrdersManager.GetOrder(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return domain.Order{}, ErrNotFound
		}
		return domain.Order{}, fmt.Errorf("%s: failed to get order:%w", op, err)
	}

	return order, nil
}

func (u *UseCase) DeleteOrder(ctx context.Context, strUUID string) error {
	const op = "usecase.DeleteOrder"

	id, err := uuid.Parse(strUUID)
	if err != nil {
		return ErrInvalidUUID
	}
	if err := u.OrdersManager.DeleteOrder(ctx, id); err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("%s: failed to delete order:%w", op, err)
	}

	return nil
}
