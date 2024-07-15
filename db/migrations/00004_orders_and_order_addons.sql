-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders(
    id UUID PRIMARY KEY,
    customer_name TEXT NOT NULL,
    customer_email TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    theatre_id UUID NOT NULL REFERENCES theatres(id),
    slot_id UUID NOT NULL REFERENCES slots(id),
    total_price DOUBLE PRECISION NOT NULL,
    no_of_persons INTEGER NOT NULL DEFAULT(1),
    order_date DATE NOT NULL,
    ordered_at TIMESTAMP NOT NULL,
    UNIQUE (theatre_id, slot_id, order_date)
);

CREATE TABLE order_addons(
    order_id UUID REFERENCES orders(id),
    addon_id UUID REFERENCES addons(id),
    quantity INTEGER NOT NULL DEFAULT(1),
    PRIMARY KEY(order_id, addon_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
DROP TABLE order_addons;
-- +goose StatementEnd
