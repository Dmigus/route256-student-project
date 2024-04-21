-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION hstore;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE message_outbox ADD COLUMN IF NOT EXISTS tracing hstore;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE message_outbox DROP COLUMN IF EXISTS tracing;
-- +goose StatementEnd