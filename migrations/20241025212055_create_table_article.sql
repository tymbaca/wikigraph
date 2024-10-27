-- +goose Up
-- +goose StatementBegin
CREATE TABLE article (
    id  INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    url TEXT UNIQUE NOT NULL,
    status TEXT NOT NULL DEFAULT 'PENDING',
    error TEXT
);

CREATE UNIQUE INDEX idx_article_unique ON article(url);
CREATE INDEX idx_article_status ON article(status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE article;
-- +goose StatementEnd
