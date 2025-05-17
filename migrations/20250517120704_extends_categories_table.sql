-- +goose Up
-- +goose StatementBegin
ALTER TABLE categories
    ADD COLUMN description VARCHAR,
    ADD COLUMN icon        VARCHAR;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE categories
    DROP COLUMN description,
    DROP COLUMN icon;
-- +goose StatementEnd
