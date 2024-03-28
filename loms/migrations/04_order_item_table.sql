-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_item
(
    order_id INTEGER NOT NULL,
    sku_id INTEGER NOT NULL,
    count INTEGER NOT NULL,
    PRIMARY KEY (order_id, sku_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_item;
-- +goose StatementEnd