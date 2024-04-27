// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: queries.sql

package orders

import (
	"context"
)

const selectAllOrders = `-- name: selectAllOrders :many
SELECT id, user_id, status, are_items_reserved
FROM "order"
ORDER BY id desc
`

func (q *Queries) selectAllOrders(ctx context.Context) ([]Order, error) {
	rows, err := q.db.Query(ctx, selectAllOrders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Order
	for rows.Next() {
		var i Order
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Status,
			&i.AreItemsReserved,
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

const selectAllOrdersItems = `-- name: selectAllOrdersItems :many
SELECT order_id, sku_id, count
FROM order_item
ORDER BY order_id desc
`

func (q *Queries) selectAllOrdersItems(ctx context.Context) ([]OrderItem, error) {
	rows, err := q.db.Query(ctx, selectAllOrdersItems)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OrderItem
	for rows.Next() {
		var i OrderItem
		if err := rows.Scan(&i.OrderID, &i.SkuID, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectOrder = `-- name: selectOrder :one
SELECT user_id, status, are_items_reserved
FROM "order"
WHERE id = $1
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
