package operations

import (
	"database/sql"
	"fmt"

	"github.com/inodaf/neoman/internal/management"
	"github.com/inodaf/neoman/pkg/git"
)

func GetDocs(owner, repo string) error {
	remote := git.NewGitHubClient()

	err := remote.IsDocsDirPresent(owner, repo)
	if err != nil {
		return fmt.Errorf("neoman: Could not locate 'docs/' from '%s/%s' on GitHub.\nMake sure you have reading rights", owner, repo)
	}

	fmt.Printf("neoman:	Fetching docs for '%s/%s' from GitHub...\n", owner, repo)

	err = management.RegistryAddEntry(management.RegistryEntry{
		Scope:   management.RegistryTypeRemote,
		Owner:   owner,
		Project: repo,
	})
	if err != nil {
		return err
	}

	// Sync docs to database after successful clone
	fmt.Printf("neoman:	Syncing docs to database for '%s/%s'...\n", owner, repo)

	db, err := management.NewSQLiteDatabase()
	if err != nil {
		return fmt.Errorf("neoman: could not open database for sync: %w", err)
	}
	defer db.Close()

	return Sync(owner, repo, db)
}

// QueryDocumentContent retrieves document content and title from the database
// using case-insensitive path matching
func QueryDocumentContent(db *sql.DB, author, repository, relativePath string) (content, title string, err error) {
	// Query with case-insensitive path matching
	query := `
		SELECT content, title FROM docpages 
		WHERE author = ? AND repository = ? 
		  AND LOWER(relative_path) = LOWER(?)
		LIMIT 1
	`

	row := db.QueryRow(query, author, repository, relativePath)
	err = row.Scan(&content, &title)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("document not found")
		}
		return "", "", fmt.Errorf("database query failed: %w", err)
	}

	return content, title, nil
}
