-- name: selectOrder :one
SELECT user_id, status, are_items_reserved
FROM "order"
WHERE id = $1;

-- name: selectOrderItems :many
SELECT sku_id, count
FROM order_item
WHERE order_id = $1;

-- name: selectCount :one
SELECT total, reserved
FROM item_unit
WHERE sku_id = $1;