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

CREATE TABLE IF NOT EXISTS articles
(
    id    SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    text  TEXT
);

CREATE TABLE IF NOT EXISTS article_relations
(
    article_id  INT NOT NULL REFERENCES articles (id) ON DELETE CASCADE,
    crop_id     INT NOT NULL REFERENCES crops (id) ON DELETE CASCADE,
    category_id INT NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
    PRIMARY KEY (article_id, crop_id, category_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS article_relations;
DROP TABLE IF EXISTS crops;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS articles;
-- +goose StatementEnd
