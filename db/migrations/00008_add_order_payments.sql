-- +goose Up
-- +goose StatementBegin
CREATE TABLE payments(
    razorpay_order_id TEXT PRIMARY KEY,
    razorpay_payment_id TEXT NOT NULL,
    razorpay_signature TEXT NOT NULL,
    status TEXT NOT NULl
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE payments;
-- +goose StatementEnd
