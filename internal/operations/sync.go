package operations

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/inodaf/neoman/internal/management"
	"github.com/inodaf/neoman/internal/models"
	"github.com/inodaf/neoman/pkg/config"
)

// Sync traverses the docs directory of a repository, parses markdown files,
// creates DocPage models, and syncs them to the database. It also removes
// stale entries from the database for files that no longer exist on the filesystem.
func Sync(author string, repo string, db *sql.DB) error {
	if author == "" {
		return fmt.Errorf("neoman: sync error - author cannot be empty")
	}

	if repo == "" {
		return fmt.Errorf("neoman: sync error - repo cannot be empty")
	}

	baseDir, err := management.RegistryEntryDirPath(management.RegistryEntry{
		Scope:   management.RegistryTypeRemote,
		Owner:   author,
		Project: repo,
	})
	if err != nil {
		return fmt.Errorf("neoman: sync error - failed to get registry entry directory path: %w", err)
	}

	docsDir := path.Join(baseDir, config.PrimaryDocsDirName)

	// Verify docs directory exists
	if _, err := os.Stat(docsDir); os.IsNotExist(err) {
		return fmt.Errorf("neoman: sync error - docs directory not found at %s", docsDir)
	}

	// Track all markdown files found during traversal
	filesOnDisk := make(map[string]bool)

	// Walk the docs directory and sync each markdown file
	err = filepath.Walk(docsDir, func(filePath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip hidden files
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		// Only process markdown files
		if !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		// Calculate relative path from docs directory
		relPath, err := filepath.Rel(docsDir, filePath)
		if err != nil {
			return fmt.Errorf("failed to calculate relative path for %s: %w", filePath, err)
		}

		// Track this file
		filesOnDisk[relPath] = true

		// Read file content
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", filePath, err)
		}

		contentStr := string(content)

		// Extract title from markdown
		title := parseMarkdownTitle(contentStr, info.Name())

		// Check if document already exists in database
		existingLastModified, err := getDocPageLastModified(db, author, repo, relPath)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("failed to query database for %s: %w", relPath, err)
		}

		// If document exists and modification time hasn't changed, skip it
		if err == nil && existingLastModified.Equal(info.ModTime()) {
			return nil
		}

		// Create DocPage model
		docPage := models.DocPage{
			Author:         author,
			Repository:     repo,
			RelativePath:   relPath,
			Title:          title,
			Content:        contentStr,
			LastModifiedAt: info.ModTime(),
			Vector:         nil, // Placeholder for future semantic search
		}

		// Insert or update document in database
		if err == sql.ErrNoRows {
			// Document doesn't exist, insert it
			err = insertDocPage(db, docPage)
		} else {
			// Document exists, update it
			err = updateDocPage(db, docPage)
		}

		if err != nil {
			return fmt.Errorf("failed to sync document %s: %w", relPath, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("neoman: sync error - failed during directory traversal: %w", err)
	}

	// Clean up stale entries from database
	err = cleanupStaleEntries(db, author, repo, filesOnDisk)
	if err != nil {
		return fmt.Errorf("neoman: sync error - failed during cleanup: %w", err)
	}

	return nil
}

// parseMarkdownTitle extracts the title from markdown content using a priority system:
// 1. YAML frontmatter (if present, looks for title field)
// 2. First H1 heading (# Title)
// 3. Filename without .md extension
func parseMarkdownTitle(content string, filename string) string {
	lines := strings.Split(content, "\n")

	// Check for YAML frontmatter at the start
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		// Look for title field in YAML frontmatter
		for i := 1; i < len(lines); i++ {
			line := strings.TrimSpace(lines[i])

			// End of frontmatter
			if line == "---" {
				break
			}

			// Check for title field
			if strings.HasPrefix(line, "title:") {
				title := strings.TrimSpace(strings.TrimPrefix(line, "title:"))
				// Remove quotes if present
				title = strings.Trim(title, "\"'")
				if title != "" {
					return title
				}
			}
		}
	}

	// Check for first H1 heading
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			title := strings.TrimSpace(strings.TrimPrefix(trimmed, "# "))
			if title != "" {
				return title
			}
		}
	}

	// Fallback: use filename without .md extension
	filenameWithoutExt := strings.TrimSuffix(filename, ".md")
	return filenameWithoutExt
}

// getDocPageLastModified queries the database to retrieve the last_modified_at
// timestamp for a document. Returns sql.ErrNoRows if document doesn't exist.
func getDocPageLastModified(db *sql.DB, author, repo, relPath string) (time.Time, error) {
	query := `SELECT last_modified_at FROM docpages
	          WHERE author = ? AND repository = ? AND relative_path = ?`

	var lastModified time.Time
	err := db.QueryRow(query, author, repo, relPath).Scan(&lastModified)

	return lastModified, err
}

// insertDocPage inserts a new DocPage record into the database.
func insertDocPage(db *sql.DB, doc models.DocPage) error {
	query := `INSERT INTO docpages
	          (author, repository, relative_path, title, content, last_modified_at, vector)
	          VALUES (?, ?, ?, ?, ?, ?, NULL)`

	_, err := db.Exec(query, doc.Author, doc.Repository, doc.RelativePath, doc.Title, doc.Content, doc.LastModifiedAt)

	return err
}

// updateDocPage updates an existing DocPage record in the database.
func updateDocPage(db *sql.DB, doc models.DocPage) error {
	query := `UPDATE docpages
	          SET title = ?, content = ?, last_modified_at = ?
	          WHERE author = ? AND repository = ? AND relative_path = ?`

	_, err := db.Exec(query, doc.Title, doc.Content, doc.LastModifiedAt, doc.Author, doc.Repository, doc.RelativePath)

	return err
}

// cleanupStaleEntries removes DocPage records from the database for files that
// no longer exist in the filesystem.
func cleanupStaleEntries(db *sql.DB, author, repo string, filesOnDisk map[string]bool) error {
	// Query all documents for this repository
	query := `SELECT relative_path FROM docpages
	          WHERE author = ? AND repository = ?`

	rows, err := db.Query(query, author, repo)
	if err != nil {
		return err
	}
	defer rows.Close()

	var relPaths []string
	for rows.Next() {
		var relPath string
		if err := rows.Scan(&relPath); err != nil {
			return err
		}
		relPaths = append(relPaths, relPath)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// Delete entries that are no longer on disk
	deleteQuery := `DELETE FROM docpages
	                WHERE author = ? AND repository = ? AND relative_path = ?`

	for _, relPath := range relPaths {
		if !filesOnDisk[relPath] {
			_, err := db.Exec(deleteQuery, author, repo, relPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
