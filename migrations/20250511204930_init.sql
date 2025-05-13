-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS crops
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS categories
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS crops_categories
(
    id          SERIAL PRIMARY KEY,
    crop_id     INT NOT NULL REFERENCES crops (id) ON DELETE CASCADE,
    category_id INT NOT NULL REFERENCES categories (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS articles
(
    id          SERIAL PRIMARY KEY,
    crop_id     INT            REFERENCES crops (id) ON DELETE SET NULL,
    category_id INT            REFERENCES categories (id) ON DELETE SET NULL,
    title       VARCHAR UNIQUE NOT NULL,
    text        TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS crops_categories;
DROP TABLE IF EXISTS crops;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS articles;
-- +goose StatementEnd
