-- +goose Up
CREATE TABLE IF NOT EXISTS docpages (
    author TEXT NOT NULL,
    repository TEXT NOT NULL,
    relative_path TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    last_modified_at DATETIME NOT NULL,
    vector BLOB
);

-- +goose Down
DROP TABLE IF EXISTS docpages;
