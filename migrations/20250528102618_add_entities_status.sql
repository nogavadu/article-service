-- +goose Up
-- +goose StatementBegin
CREATE TABLE entity_status
(
    id     INT PRIMARY KEY,
    status VARCHAR NOT NULL UNIQUE
);

INSERT INTO entity_status (status)
VALUES ('published');
INSERT INTO entity_status (status)
VALUES ('review');
INSERT INTO entity_status (status)
VALUES ('canceled');

ALTER TABLE crops
    ADD COLUMN IF NOT EXISTS status INT REFERENCES entity_status (id) ON DELETE SET DEFAULT DEFAULT (2);
ALTER TABLE categories
    ADD COLUMN IF NOT EXISTS status INT REFERENCES entity_status (id) ON DELETE SET DEFAULT DEFAULT (2);
ALTER TABLE articles
    ADD COLUMN IF NOT EXISTS status INT REFERENCES entity_status (id) ON DELETE SET DEFAULT DEFAULT (2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE crops
    DROP COLUMN IF EXISTS status;
ALTER TABLE categories
    DROP COLUMN IF EXISTS status;
ALTER TABLE articles
    DROP COLUMN IF EXISTS status;

DROP TABLE IF EXISTS entity_status;
-- +goose StatementEnd
