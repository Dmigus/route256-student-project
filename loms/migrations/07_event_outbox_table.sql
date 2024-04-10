-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS event_outbox
(
    id BIGSERIAL NOT NULL primary key,
    order_id BIGINT NOT NULL,
    message VARCHAR NOT NULL,
    at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS event_outbox;
-- +goose StatementEnd