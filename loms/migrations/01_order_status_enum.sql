-- +goose Up
-- +goose StatementBegin
CREATE TYPE order_status AS ENUM ('Undefined', 'New', 'AwaitingPayment', 'Failed', 'Payed', 'Cancelled');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS order_status;
-- +goose StatementEnd
