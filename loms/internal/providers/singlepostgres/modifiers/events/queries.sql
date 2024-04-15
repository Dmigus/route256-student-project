-- name: pushEvent :exec
INSERT INTO message_outbox(partition_key, payload)
VALUES ($1, $2);

-- name: pullEvents :many
DELETE FROM message_outbox
WHERE id IN
      (SELECT id from message_outbox ORDER BY id LIMIT $1 FOR UPDATE)
RETURNING *;

-- name: currentDatetime :one
SELECT clock_timestamp()::timestamp;