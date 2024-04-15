

-- name: selectCount :one
SELECT total, reserved
FROM item_unit
WHERE sku_id = $1;