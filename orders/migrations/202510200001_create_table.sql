-- +goose Up
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY,
    status TEXT NOT NULL,
    customer_id UUID NOT NULL,
    animal_type TEXT NOT NULL,
    animal_age INTEGER NOT NULL,
    delete_reason TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders (customer_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders (status);
CREATE INDEX IF NOT EXISTS idx_orders_updated_at ON orders (updated_at DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_orders_updated_at;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_customer_id;
DROP TABLE IF EXISTS orders;
