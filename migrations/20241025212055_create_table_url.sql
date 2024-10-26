-- +goose Up
-- +goose StatementBegin
CREATE TABLE url (
    id  INTEGER PRIMARY KEY AUTOINCREMENT,
    url TEXT UNIQUE NOT NULL,
    status TEXT NOT NULL DEFAULT 'PENDING',
    error TEXT
);

CREATE UNIQUE INDEX idx_url_unique ON url(url);
CREATE INDEX idx_url_status ON url(status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE url;
-- +goose StatementEnd
