-- +goose Up
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