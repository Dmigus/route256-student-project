-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS event_outbox
(
    order_id BIGINT NOT NULL,
    type VARCHAR NOT NULL,
    at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS event_outbox;
-- +goose StatementEnd