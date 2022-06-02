CREATE TYPE taxi_type AS ENUM ('economy', 'comfort', 'business', 'electro');

CREATE TABLE IF NOT EXISTS orders
(
    id SERIAL PRIMARY KEY,
    user_uuid uuid NOT NULL,
    driver_uuid uuid NOT NULL,
    origin VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    taxi_type taxi_type NOT NULL,
    is_finished BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- name: CreateOrder :one
INSERT INTO orders (user_uuid, driver_uuid, origin, destination, taxi_type)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetOrders :many
SELECT * FROM orders;

-- name: UpdateOrderStatus :one
UPDATE orders
SET is_finished = TRUE
WHERE id=$1 RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;
