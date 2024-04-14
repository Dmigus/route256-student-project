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
