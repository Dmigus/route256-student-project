-- name: createOrder :one
INSERT INTO "order"(user_id, status, are_items_reserved)
VALUES ($1, $2, $3)
RETURNING id;

-- name: selectOrder :one
SELECT user_id, status, are_items_reserved
FROM "order"
WHERE id = $1
    FOR UPDATE;

-- name: selectOrderItems :many
SELECT sku_id, count
FROM order_item
WHERE order_id = $1;

-- name: updateOrder :exec
UPDATE "order"
SET status = $2, are_items_reserved = $3
where id = $1;

-- name: insertOrderItem :copyfrom
INSERT INTO order_item(order_id, sku_id, count)
VALUES ($1, $2, $3);