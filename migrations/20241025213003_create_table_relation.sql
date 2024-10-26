-- +goose Up
-- +goose StatementBegin
CREATE TABLE relation (
    from_id INTEGER NOT NULL,
    to_id   INTEGER NOT NULL,
    UNIQUE(from_id, to_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE relation;
-- +goose StatementEnd
