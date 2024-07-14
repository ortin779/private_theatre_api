-- +goose Up
-- +goose StatementBegin
CREATE TABLE theatres(
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    price DOUBLE PRECISION NOT NULL,
    additional_price_per_head DOUBLE PRECISION NOT NULL,
    max_capacity INTEGER NOT NULL,
    min_capacity INTEGER DEFAULT(1),
    default_capacity INTEGER NOT NULL
);

CREATE TABLE theatre_slots(
    theatre_id UUID REFERENCES theatres(id),
    slot_id UUID REFERENCES slots(id),
    PRIMARY KEY(theatre_id, slot_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE theatres;
DROP TABLE theatre_slots;
-- +goose StatementEnd
