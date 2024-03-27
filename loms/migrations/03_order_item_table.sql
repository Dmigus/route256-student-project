-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_item
(
    order_id INTEGER NOT NULL,
    sku_id INTEGER NOT NULL,
    count INTEGER NOT NULL,
    CONSTRAINT UQ_order_item_order_id_sku_id UNIQUE (order_id, sku_id)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS ix_order_item_order_id ON order_item
(
     order_id
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS ix_order_item_order_id;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS order_item;
-- +goose StatementEnd