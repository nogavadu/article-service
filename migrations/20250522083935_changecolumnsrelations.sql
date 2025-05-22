-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS crops_categories
(
    crop_id INT REFERENCES crops(id),
    category_id INT REFERENCES categories(id),
    PRIMARY KEY (crop_id, category_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS crops_categories;
-- +goose StatementEnd
