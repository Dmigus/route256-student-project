-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "order"
(
    id INTEGER NOT NULL primary key,
    user_id INTEGER NOT NULL,
    status order_status NOT NULL ,
    are_items_reserved BOOLEAN NOT NULL
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS ix_order_id ON "order"
(
    id
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS ix_order_id;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS "order";
-- +goose StatementEnd