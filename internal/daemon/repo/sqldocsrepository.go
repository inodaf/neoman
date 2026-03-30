package repo

import (
	"database/sql"
	"fmt"

	"github.com/inodaf/neoman/internal/daemon/domain"
)

type docsRepository struct {
	db *sql.DB
}

func NewDocsRepository(db *sql.DB) DocsRepository {
	return &docsRepository{db: db}
}

func (r *docsRepository) Save(entry domain.RemoteDocs) error {
	query := `
		INSERT INTO remote_docs (author, repository, source, indexing, last_sync_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(author, repository, source) DO UPDATE SET
			indexing = excluded.indexing,
			last_sync_at = excluded.last_sync_at,
			updated_at = excluded.updated_at
	`

	_, err := r.db.Exec(
		query,
		entry.Author,
		entry.Repository,
		entry.Source,
		entry.Indexing,
		entry.LastSyncAt,
		entry.CreatedAt,
		entry.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save remote docs: %w", err)
	}

	return nil
}

func (r *docsRepository) Exists(author, repository string) (bool, error) {
	query := `SELECT 1 FROM remote_docs WHERE author = ? AND repository = ? LIMIT 1`

	var exists int
	err := r.db.QueryRow(query, author, repository).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check if remote docs exists: %w", err)
	}

	return true, nil
}

func (r *docsRepository) GetOne(author, repository string, source domain.RemoteSource) (domain.RemoteDocs, error) {
	query := `
		SELECT author, repository, source, indexing, last_sync_at, created_at, updated_at
		FROM remote_docs
		WHERE author = ? AND repository = ? AND source = ?
	`

	var docs domain.RemoteDocs
	err := r.db.QueryRow(query, author, repository, source).Scan(
		&docs.Author,
		&docs.Repository,
		&docs.Source,
		&docs.Indexing,
		&docs.LastSyncAt,
		&docs.CreatedAt,
		&docs.UpdatedAt,
	)
	if err != nil {
		return domain.RemoteDocs{}, fmt.Errorf("failed to get remote docs: %w", err)
	}

	return docs, nil
}
