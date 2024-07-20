-- +goose Up
-- +goose StatementBegin
ALTER TABLE slots
    ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN created_by UUID NOT NULL REFERENCES users(id),
    ADD COLUMN updated_by UUID NOT NULL REFERENCES users(id);

ALTER TABLE theatres
    ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN created_by UUID NOT NULL REFERENCES users(id),
    ADD COLUMN updated_by UUID NOT NULL REFERENCES users(id);

ALTER TABLE addons
    ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN created_by UUID NOT NULL REFERENCES users(id),
    ADD COLUMN updated_by UUID NOT NULL REFERENCES users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE slots
    DROP COLUMN updated_at,
    DROP COLUMN create_at,
    DROP COLUMN created_by,
    DROP COLUMN updated_by;

ALTER TABLE theatres
    DROP COLUMN updated_at,
    DROP COLUMN create_at,
    DROP COLUMN created_by,
    DROP COLUMN updated_by;

ALTER TABLE addons
    DROP COLUMN updated_at,
    DROP COLUMN create_at,
    DROP COLUMN created_by,
    DROP COLUMN updated_by;
-- +goose StatementEnd
