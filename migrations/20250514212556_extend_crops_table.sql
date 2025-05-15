-- +goose Up
-- +goose StatementBegin
ALTER TABLE crops
    ADD COLUMN description VARCHAR,
    ADD COLUMN img         VARCHAR;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE crops
    DROP COLUMN description,
    DROP COLUMN img;
-- +goose StatementEnd
