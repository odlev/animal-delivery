// Package postgres is a nice package
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/odlev/animal-delivery/orders/internal/domain"
	"github.com/odlev/animal-delivery/orders/internal/repo"
	"github.com/pressly/goose/v3"
)

type Storage struct {
	pool *pgxpool.Pool
}

func Init(ctx context.Context, dsn string) (*Storage, error) {
	if err := applyMigrations(ctx, dsn); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres.Init: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres.Init: ping: %w", err)
	}

	return &Storage{pool: pool}, nil
}

func (s *Storage) Close() {
	if s != nil && s.pool != nil {
		s.pool.Close()
	}
}

func applyMigrations(ctx context.Context, dsn string) error {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse config pgx: %w", err)
	}
	db := stdlib.OpenDB(*cfg.ConnConfig)
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.UpContext(ctx, db, "./migrations")
}

func (s *Storage) CreateOrder(ctx context.Context, order domain.Order) error {
	const op = "repo.Postgres.CreateOrder"

	const query = `
        INSERT INTO orders (
            id,
            status,
            customer_id,
            animal_type,
            animal_age,
            delete_reason,
            updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	_, err := s.pool.Exec(
		ctx,
		query,
		order.OrderID,
		order.Status,
		order.CustomerID,
		order.AnimalType,
		order.AnimalAge,
		order.DeleteReason,
		order.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetOrder(ctx context.Context, id uuid.UUID) (domain.Order, error) {
	const op = "repo.Postgres.GetOrder"

	const query = `
        SELECT
            id,
            status,
            customer_id,
            animal_type,
            animal_age,
            delete_reason,
            updated_at
        FROM orders
        WHERE id = $1
    `

	var row orderRow
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&row.ID,
		&row.Status,
		&row.CustomerID,
		&row.AnimalType,
		&row.AnimalAge,
		&row.DeleteReason,
		&row.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Order{}, fmt.Errorf("%s: %w", op, repo.ErrNotFound)
	}
	if err != nil {
		return domain.Order{}, fmt.Errorf("%s: %w", op, err)
	}

	return row.toDomain()
}

func (s *Storage) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	const op = "repo.Postgres.DeleteOrder"

	const query = `DELETE FROM orders WHERE id = $1`

	tag, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if tag.RowsAffected() == 0 {
		return repo.ErrNotFound
	}

	return nil
}

func (r orderRow) toDomain() (domain.Order, error) {
	return domain.Order{
		OrderID:      r.ID,
		Status:       r.Status,
		CustomerID:   r.CustomerID,
		AnimalType:   r.AnimalType,
		AnimalAge:    r.AnimalAge,
		DeleteReason: r.DeleteReason,
		UpdatedAt:    r.UpdatedAt,
	}, nil
}
