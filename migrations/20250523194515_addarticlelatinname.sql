-- +goose Up
-- +goose StatementBegin
ALTER TABLE articles
    ADD COLUMN IF NOT EXISTS latin_name VARCHAR;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE articles
    DROP COLUMN IF EXISTS latin_name;
-- +goose StatementEnd
