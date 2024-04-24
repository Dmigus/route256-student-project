-- name: createOrder :one
INSERT INTO "order"(id, user_id, status, are_items_reserved)
VALUES (nextval('order_id_seq')*1000 +$1, $2, $3, $4)
RETURNING id;

-- name: insertOrderItem :copyfrom
INSERT INTO order_item(order_id, sku_id, count)
VALUES ($1, $2, $3);