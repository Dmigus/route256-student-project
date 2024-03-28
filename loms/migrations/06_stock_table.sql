-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS item_unit
(
    sku_id INTEGER PRIMARY KEY,
    total INTEGER NOT NULL,
    reserved INTEGER NOT NULL,
    CHECK (total >= 0),
    CHECK (reserved >= 0),
    CHECK (total >= reserved)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS item_unit;
-- +goose StatementEnd