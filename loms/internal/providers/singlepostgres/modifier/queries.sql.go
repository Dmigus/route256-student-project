// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: queries.sql

package modifier

import (
	"context"
)

const createOrder = `-- name: createOrder :one
INSERT INTO "order"(user_id, status, are_items_reserved)
VALUES ($1, $2, $3)
RETURNING id
`

type createOrderParams struct {
	UserID           int64
	Status           string
	AreItemsReserved bool
}

func (q *Queries) createOrder(ctx context.Context, arg createOrderParams) (int64, error) {
	row := q.db.QueryRow(ctx, createOrder, arg.UserID, arg.Status, arg.AreItemsReserved)
	var id int64
	err := row.Scan(&id)
	return id, err
}

type insertOrderItemParams struct {
	OrderID int64
	SkuID   int64
	Count   int32
}

const insertStock = `-- name: insertStock :exec
INSERT INTO item_unit(sku_id, total, reserved)
VALUES ($1, $2, $3)
ON CONFLICT (sku_id)
    DO UPDATE SET total=$2, reserved=$3
`

type insertStockParams struct {
	SkuID    int64
	Total    int32
	Reserved int32
}

func (q *Queries) insertStock(ctx context.Context, arg insertStockParams) error {
	_, err := q.db.Exec(ctx, insertStock, arg.SkuID, arg.Total, arg.Reserved)
	return err
}

const pullEvents = `-- name: pullEvents :many
DELETE FROM event_outbox
WHERE id IN
(SELECT id from event_outbox ORDER BY id LIMIT $1)
RETURNING id, order_id, message, at
`

func (q *Queries) pullEvents(ctx context.Context, limit int32) ([]EventOutbox, error) {
	rows, err := q.db.Query(ctx, pullEvents, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []EventOutbox
	for rows.Next() {
		var i EventOutbox
		if err := rows.Scan(
			&i.ID,
			&i.OrderID,
			&i.Message,
			&i.At,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const pushEvent = `-- name: pushEvent :exec
INSERT INTO event_outbox(order_id, message, at)
VALUES ($1, $2, clock_timestamp())
`

type pushEventParams struct {
	OrderID int64
	Message string
}

func (q *Queries) pushEvent(ctx context.Context, arg pushEventParams) error {
	_, err := q.db.Exec(ctx, pushEvent, arg.OrderID, arg.Message)
	return err
}

const selectCount = `-- name: selectCount :one
SELECT total, reserved
FROM item_unit
WHERE sku_id = $1
FOR UPDATE
`

type selectCountRow struct {
	Total    int32
	Reserved int32
}

func (q *Queries) selectCount(ctx context.Context, skuID int64) (selectCountRow, error) {
	row := q.db.QueryRow(ctx, selectCount, skuID)
	var i selectCountRow
	err := row.Scan(&i.Total, &i.Reserved)
	return i, err
}

const selectOrder = `-- name: selectOrder :one
SELECT user_id, status, are_items_reserved
FROM "order"
WHERE id = $1
FOR UPDATE
`

type selectOrderRow struct {
	UserID           int64
	Status           string
	AreItemsReserved bool
}

func (q *Queries) selectOrder(ctx context.Context, id int64) (selectOrderRow, error) {
	row := q.db.QueryRow(ctx, selectOrder, id)
	var i selectOrderRow
	err := row.Scan(&i.UserID, &i.Status, &i.AreItemsReserved)
	return i, err
}

const selectOrderItems = `-- name: selectOrderItems :many
SELECT sku_id, count
FROM order_item
WHERE order_id = $1
`

type selectOrderItemsRow struct {
	SkuID int64
	Count int32
}

func (q *Queries) selectOrderItems(ctx context.Context, orderID int64) ([]selectOrderItemsRow, error) {
	rows, err := q.db.Query(ctx, selectOrderItems, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []selectOrderItemsRow
	for rows.Next() {
		var i selectOrderItemsRow
		if err := rows.Scan(&i.SkuID, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateOrder = `-- name: updateOrder :exec
UPDATE "order"
SET status = $2, are_items_reserved = $3
where id = $1
`

type updateOrderParams struct {
	ID               int64
	Status           string
	AreItemsReserved bool
}

func (q *Queries) updateOrder(ctx context.Context, arg updateOrderParams) error {
	_, err := q.db.Exec(ctx, updateOrder, arg.ID, arg.Status, arg.AreItemsReserved)
	return err
}

const updateReserved = `-- name: updateReserved :exec
UPDATE item_unit
SET reserved = $2
WHERE sku_id = $1
`

type updateReservedParams struct {
	SkuID    int64
	Reserved int32
}

func (q *Queries) updateReserved(ctx context.Context, arg updateReservedParams) error {
	_, err := q.db.Exec(ctx, updateReserved, arg.SkuID, arg.Reserved)
	return err
}

const updateTotalReserved = `-- name: updateTotalReserved :exec
UPDATE item_unit
SET total = $2, reserved = $3
WHERE sku_id = $1
`

type updateTotalReservedParams struct {
	SkuID    int64
	Total    int32
	Reserved int32
}

func (q *Queries) updateTotalReserved(ctx context.Context, arg updateTotalReservedParams) error {
	_, err := q.db.Exec(ctx, updateTotalReserved, arg.SkuID, arg.Total, arg.Reserved)
	return err
}
