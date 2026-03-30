-- +goose Up
CREATE TABLE IF NOT EXISTS remote_docs (
    author TEXT NOT NULL,
    repository TEXT NOT NULL,
    source TEXT NOT NULL,
    indexing BOOLEAN NOT NULL DEFAULT 0,
    last_sync_at DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    PRIMARY KEY (author, repository, source)
);

CREATE INDEX IF NOT EXISTS idx_remote_docs_author_repo
    ON remote_docs(author, repository);

-- +goose Down
DROP INDEX IF EXISTS idx_remote_docs_author_repo;
DROP TABLE IF EXISTS remote_docs;
