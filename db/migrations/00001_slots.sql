-- +goose Up
-- +goose StatementBegin
CREATE TABLE slots(
    id UUID PRIMARY KEY,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    UNIQUE (start_time, end_time)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE slots;
-- +goose StatementEnd
