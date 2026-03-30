package repo

import (
	"database/sql"
	"fmt"

	"github.com/inodaf/neoman/internal/daemon/domain"
)

type docsPageRepository struct {
	db *sql.DB
}

func NewDocsPageRepository(db *sql.DB) DocsPageRepository {
	return &docsPageRepository{db: db}
}

func (r *docsPageRepository) Save(page domain.DocsPage) error {
	query := `
		INSERT INTO docpages (author, repository, relative_path, title, content, last_modified_at, vector)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(author, repository, relative_path) DO UPDATE SET
			title = excluded.title,
			content = excluded.content,
			last_modified_at = excluded.last_modified_at
	`

	_, err := r.db.Exec(
		query,
		page.Author,
		page.Repository,
		page.RelativePath,
		page.Title,
		page.Content,
		page.LastModifiedAt,
		nil, // vector handling for later
	)
	if err != nil {
		return fmt.Errorf("failed to save docs page: %w", err)
	}

	return nil
}

func (r *docsPageRepository) SaveMany(pages []domain.DocsPage) error {
	for _, page := range pages {
		err := r.Save(page)
		if err != nil {
			return err
		}
	}
	return nil
}
