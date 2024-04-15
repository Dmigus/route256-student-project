-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS message_outbox
(
    id BIGSERIAL NOT NULL primary key,
    partition_key BYTEA NOT NULL,
    payload BYTEA NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS message_outbox;
-- +goose StatementEnd