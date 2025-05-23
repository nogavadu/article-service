-- +goose Up
-- +goose StatementBegin
ALTER TABLE article_relations RENAME TO  articles_relations;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE articles_relations RENAME TO  article_relations;
-- +goose StatementEnd
