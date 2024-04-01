-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS ix_item_unit_sku_id ON item_unit
(
     sku_id
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS ix_item_unit_sku_id;
-- +goose StatementEnd