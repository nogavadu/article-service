-- +goose Up
-- +goose StatementBegin
ALTER TABLE crops_categories
    DROP CONSTRAINT crops_categories_crop_id_fkey,
    DROP CONSTRAINT crops_categories_category_id_fkey,
    ADD CONSTRAINT crops_categories_crop_id_fkey
        FOREIGN KEY (crop_id) REFERENCES crops(id) ON DELETE CASCADE,
    ADD CONSTRAINT crops_categories_category_id_fkey
        FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE crops_categories
    DROP CONSTRAINT crops_categories_crop_id_fkey,
    DROP CONSTRAINT crops_categories_category_id_fkey,
    ADD CONSTRAINT crops_categories_crop_id_fkey
        FOREIGN KEY (crop_id) REFERENCES crops(id),
    ADD CONSTRAINT crops_categories_category_id_fkey
        FOREIGN KEY (category_id) REFERENCES categories(id);
-- +goose StatementEnd
