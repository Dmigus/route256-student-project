-- name: selectOrder :one
SELECT user_id, status, are_items_reserved
FROM "order"
WHERE id = $1;

-- name: selectOrderItems :many
SELECT sku_id, count
FROM order_item
WHERE order_id = $1;

-- name: selectAllOrders :many
SELECT id, user_id, status, are_items_reserved
FROM "order"
ORDER BY id desc;

-- name: selectAllOrdersItems :many
SELECT order_id, sku_id, count
FROM order_item
ORDER BY order_id desc;