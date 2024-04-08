-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS ix_event_outbox_at ON event_outbox USING btree
(
     at
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS ix_event_outbox_at;
-- +goose StatementEnd