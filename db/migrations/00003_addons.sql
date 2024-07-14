-- +goose Up
-- +goose StatementBegin
CREATE TABLE addons(
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    price DOUBLE PRECISION NOT NULL,
    meta_data JSONB
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE addons;
-- +goose StatementEnd
