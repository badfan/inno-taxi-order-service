// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: orders.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
)

const createOrder = `-- name: CreateOrder :one
INSERT INTO orders (user_uuid, driver_uuid, origin, destination, taxi_type)
VALUES ($1, $2, $3, $4, $5) RETURNING id, user_uuid, driver_uuid, origin, destination, taxi_type, is_finished, created_at, updated_at
`

type CreateOrderParams struct {
	UserUuid    uuid.UUID `json:"user_uuid"`
	DriverUuid  uuid.UUID `json:"driver_uuid"`
	Origin      string    `json:"origin"`
	Destination string    `json:"destination"`
	TaxiType    TaxiType  `json:"taxi_type"`
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, createOrder,
		arg.UserUuid,
		arg.DriverUuid,
		arg.Origin,
		arg.Destination,
		arg.TaxiType,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserUuid,
		&i.DriverUuid,
		&i.Origin,
		&i.Destination,
		&i.TaxiType,
		&i.IsFinished,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteOrder = `-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1
`

func (q *Queries) DeleteOrder(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteOrder, id)
	return err
}

const getOrdersByUUID = `-- name: GetOrdersByUUID :many
SELECT id, user_uuid, driver_uuid, origin, destination, taxi_type, is_finished, created_at, updated_at FROM orders
WHERE user_uuid=$1 OR driver_uuid=$1
`

func (q *Queries) GetOrdersByUUID(ctx context.Context, userUuid uuid.UUID) ([]Order, error) {
	rows, err := q.db.QueryContext(ctx, getOrdersByUUID, userUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Order
	for rows.Next() {
		var i Order
		if err := rows.Scan(
			&i.ID,
			&i.UserUuid,
			&i.DriverUuid,
			&i.Origin,
			&i.Destination,
			&i.TaxiType,
			&i.IsFinished,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateOrderStatus = `-- name: UpdateOrderStatus :one
UPDATE orders
SET is_finished = TRUE
WHERE id=$1 RETURNING id, user_uuid, driver_uuid, origin, destination, taxi_type, is_finished, created_at, updated_at
`

func (q *Queries) UpdateOrderStatus(ctx context.Context, id int32) (Order, error) {
	row := q.db.QueryRowContext(ctx, updateOrderStatus, id)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserUuid,
		&i.DriverUuid,
		&i.Origin,
		&i.Destination,
		&i.TaxiType,
		&i.IsFinished,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
