-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS item_unit
(
    sku_id INTEGER PRIMARY KEY,
    total INTEGER NOT NULL,
    reserved INTEGER NOT NULL
);
-- +goose StatementEnd
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
-- +goose StatementBegin
DROP TABLE IF EXISTS item_unit;
-- +goose StatementEnd