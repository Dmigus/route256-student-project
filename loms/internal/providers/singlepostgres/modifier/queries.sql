
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

-- name: insertStock :exec
INSERT INTO item_unit(sku_id, total, reserved)
VALUES ($1, $2, $3)
ON CONFLICT (sku_id)
    DO UPDATE SET total=$2, reserved=$3;

-- name: selectCount :one
SELECT total, reserved
FROM item_unit
WHERE sku_id = $1
FOR UPDATE;

-- name: updateReserved :exec
UPDATE item_unit
SET reserved = $2
WHERE sku_id = $1;

-- name: updateTotalReserved :exec
UPDATE item_unit
SET total = $2, reserved = $3
WHERE sku_id = $1;

-- name: pushEvent :exec
INSERT INTO event_outbox(order_id, message, at)
VALUES ($1, $2, clock_timestamp());

-- name: pullEvents :many
DELETE FROM event_outbox
WHERE id IN
(SELECT id from event_outbox ORDER BY id LIMIT $1)
RETURNING *;

