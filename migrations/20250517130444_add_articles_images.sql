-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS articles_images
(
    id         SERIAL PRIMARY KEY,
    img        VARCHAR NOT NULL,
    article_id INT REFERENCES articles (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS articles_images;
-- +goose StatementEnd
