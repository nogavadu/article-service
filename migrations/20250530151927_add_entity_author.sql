-- +goose Up
-- +goose StatementBegin
ALTER TABLE crops
    ADD COLUMN IF NOT EXISTS author INT;
ALTER TABLE categories
    ADD COLUMN IF NOT EXISTS author INT;
ALTER TABLE articles
    ADD COLUMN IF NOT EXISTS author INT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE crops
    DROP COLUMN IF EXISTS author;
ALTER TABLE categories
    DROP COLUMN IF EXISTS author;
ALTER TABLE articles
    DROP COLUMN IF EXISTS author;
-- +goose StatementEnd
