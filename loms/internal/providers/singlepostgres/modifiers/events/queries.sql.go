// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: queries.sql

package events

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const currentDatetime = `-- name: currentDatetime :one
SELECT clock_timestamp()::timestamp
`

func (q *Queries) currentDatetime(ctx context.Context) (pgtype.Timestamp, error) {
	row := q.db.QueryRow(ctx, currentDatetime)
	var column_1 pgtype.Timestamp
	err := row.Scan(&column_1)
	return column_1, err
}

const pullEvents = `-- name: pullEvents :many
DELETE FROM message_outbox
WHERE id IN
      (SELECT id from message_outbox ORDER BY id LIMIT $1 FOR UPDATE)
RETURNING id, partition_key, payload, tracing
`

func (q *Queries) pullEvents(ctx context.Context, limit int32) ([]MessageOutbox, error) {
	rows, err := q.db.Query(ctx, pullEvents, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MessageOutbox
	for rows.Next() {
		var i MessageOutbox
		if err := rows.Scan(
			&i.ID,
			&i.PartitionKey,
			&i.Payload,
			&i.Tracing,
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
INSERT INTO message_outbox(partition_key, payload, tracing)
VALUES ($1, $2, $3)
`

type pushEventParams struct {
	PartitionKey []byte
	Payload      []byte
	Tracing      pgtype.Hstore
}

func (q *Queries) pushEvent(ctx context.Context, arg pushEventParams) error {
	_, err := q.db.Exec(ctx, pushEvent, arg.PartitionKey, arg.Payload, arg.Tracing)
	return err
}
