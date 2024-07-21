-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
    ALTER COLUMN total_price TYPE INTEGER USING (total_price::integer),
    ADD COLUMN razorpay_order_id TEXT NOT NULL REFERENCES payments(razorpay_order_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
    ALTER COLUMN total_price TYPE DOUBLE PRECISION,
    DROP COLUMN razorpay_order_id;
-- +goose StatementEnd
