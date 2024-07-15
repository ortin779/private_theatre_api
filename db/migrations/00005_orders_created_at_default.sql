-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
    ALTER COLUMN ordered_at
    SET DEFAULT CURRENT_TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
    ADD COLUMN ordered_at TIMESTAMP NOT NULL;
-- +goose StatementEnd
