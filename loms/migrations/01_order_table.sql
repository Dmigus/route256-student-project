-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "order"
(
    id BIGSERIAL NOT NULL primary key,
    user_id BIGINT NOT NULL,
    status varchar NOT NULL,
    are_items_reserved BOOLEAN NOT NULL
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "order";
-- +goose StatementEnd