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

func (r *docsPageRepository) GetAll(author, repository string) ([]domain.DocsPage, error) {
	query := `
		SELECT author, repository, relative_path, title, last_modified_at
		FROM docpages
		WHERE author = ? AND repository = ?
		ORDER BY relative_path
	`

	rows, err := r.db.Query(query, author, repository)
	if err != nil {
		return nil, fmt.Errorf("failed to query doc pages: %w", err)
	}
	defer rows.Close()

	var pages []domain.DocsPage
	for rows.Next() {
		var page domain.DocsPage
		err := rows.Scan(
			&page.Author,
			&page.Repository,
			&page.RelativePath,
			&page.Title,
			&page.LastModifiedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan doc page: %w", err)
		}
		pages = append(pages, page)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating doc pages: %w", err)
	}

	return pages, nil
}

func (r *docsPageRepository) GetOne(author, repository, relativePath string) (*domain.DocsPage, error) {
	query := `
		SELECT author, repository, relative_path, title, content, last_modified_at
		FROM docpages
		WHERE author = ? AND repository = ? AND relative_path = ?
	`

	var page domain.DocsPage
	err := r.db.QueryRow(query, author, repository, relativePath).Scan(
		&page.Author,
		&page.Repository,
		&page.RelativePath,
		&page.Title,
		&page.Content,
		&page.LastModifiedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("page not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query doc page: %w", err)
	}

	return &page, nil
}
