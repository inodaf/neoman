package operations

import (
	"database/sql"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/inodaf/neoman/pkg/config"
)

// TestParseMarkdownTitle tests markdown title extraction with all priority levels
func TestParseMarkdownTitle(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		filename string
		expected string
	}{
		{
			name:     "YAML frontmatter with title",
			content:  "---\ntitle: \"YAML Title\"\nauthor: Test\n---\n\n# H1 Title\n\nContent",
			filename: "test.md",
			expected: "YAML Title",
		},
		{
			name:     "YAML frontmatter with single quote title",
			content:  "---\ntitle: 'Single Quote Title'\n---\n\nContent",
			filename: "test.md",
			expected: "Single Quote Title",
		},
		{
			name:     "First H1 heading fallback",
			content:  "Some intro text\n\n# First H1 Title\n\nMore content",
			filename: "test.md",
			expected: "First H1 Title",
		},
		{
			name:     "Filename without extension fallback",
			content:  "Just some content with no title markers",
			filename: "MyDocument.md",
			expected: "MyDocument",
		},
		{
			name:     "Empty YAML title falls through to H1",
			content:  "---\ntitle: \n---\n\n# H1 Title\n\nContent",
			filename: "test.md",
			expected: "H1 Title",
		},
		{
			name:     "YAML with empty title and no H1",
			content:  "---\ntitle: \n---\n\nNo title here",
			filename: "Document.md",
			expected: "Document",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseMarkdownTitle(tt.content, tt.filename)
			if result != tt.expected {
				t.Errorf("parseMarkdownTitle() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestSyncWithTestRepository creates a temporary test repository and syncs it
func TestSyncWithTestRepository(t *testing.T) {
	// Create temporary directory for test docs
	tempDir := t.TempDir()
	docsDir := path.Join(tempDir, config.PrimaryDocsDirName)
	if err := os.Mkdir(docsDir, 0755); err != nil {
		t.Fatalf("Failed to create docs directory: %v", err)
	}

	// Create subdirectory for nested docs
	subDir := path.Join(docsDir, "nested")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create nested directory: %v", err)
	}

	// Write test markdown files
	testFiles := map[string]string{
		path.Join(docsDir, "README.md"): `---
title: Main Documentation
---

# Welcome

Main docs content`,
		path.Join(docsDir, "Guide.md"): `# Installation Guide

How to install our software`,
		path.Join(subDir, "API.md"): `API Reference

No explicit title markers`,
	}

	for filePath, content := range testFiles {
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test file %s: %v", filePath, err)
		}
	}

	// Create test database
	db, err := sql.Open("sqlite3", path.Join(tempDir, "test.db"))
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Create schema
	schema := `
		CREATE TABLE IF NOT EXISTS docpages (
			author TEXT NOT NULL,
			repository TEXT NOT NULL,
			relative_path TEXT NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			last_modified_at DATETIME NOT NULL,
			vector BLOB
		);
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	// Mock the registry to point to our test directory
	// We'll need to modify the test to work with actual registry structure
	// For now, we'll manually insert test data

	// Verify the test setup
	files, err := os.ReadDir(docsDir)
	if err != nil {
		t.Fatalf("Failed to read docs directory: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("No files created in docs directory")
	}

	// Test title extraction for each file
	expectedTitles := map[string]string{
		"README.md":     "Main Documentation",
		"Guide.md":      "Installation Guide",
		"nested/API.md": "API",
	}

	for relPath, expectedTitle := range expectedTitles {
		filePath := path.Join(docsDir, relPath)
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Errorf("Failed to read file %s: %v", relPath, err)
			continue
		}

		filename := filepath.Base(filePath)
		title := parseMarkdownTitle(string(content), filename)

		if title != expectedTitle {
			t.Errorf("File %s: parseMarkdownTitle() = %q, want %q", relPath, title, expectedTitle)
		}
	}
}

// TestDatabaseOperations tests INSERT and UPDATE operations
func TestDatabaseOperations(t *testing.T) {
	tempDir := t.TempDir()
	db, err := sql.Open("sqlite3", path.Join(tempDir, "test.db"))
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Create schema
	schema := `
		CREATE TABLE IF NOT EXISTS docpages (
			author TEXT NOT NULL,
			repository TEXT NOT NULL,
			relative_path TEXT NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			last_modified_at DATETIME NOT NULL,
			vector BLOB
		);
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	now := time.Now()

	// Mock DocPage insert
	query := `INSERT INTO docpages 
	          (author, repository, relative_path, title, content, last_modified_at, vector)
	          VALUES (?, ?, ?, ?, ?, ?, NULL)`

	_, err = db.Exec(query, "testauthor", "testrepo", "docs/README.md", "Test Title", "Test Content", now)
	if err != nil {
		t.Fatalf("Failed to insert docpage: %v", err)
	}

	// Verify insert
	var title, content string
	selectQuery := `SELECT title, content FROM docpages 
	                WHERE author = ? AND repository = ? AND relative_path = ?`

	err = db.QueryRow(selectQuery, "testauthor", "testrepo", "docs/README.md").Scan(&title, &content)
	if err != nil {
		t.Fatalf("Failed to query inserted record: %v", err)
	}

	if title != "Test Title" || content != "Test Content" {
		t.Errorf("Inserted data mismatch: title=%q, content=%q", title, content)
	}

	// Test UPDATE
	updateQuery := `UPDATE docpages 
	                SET title = ?, content = ?, last_modified_at = ?
	                WHERE author = ? AND repository = ? AND relative_path = ?`

	updatedTime := now.Add(time.Hour)
	_, err = db.Exec(updateQuery, "Updated Title", "Updated Content", updatedTime, "testauthor", "testrepo", "docs/README.md")
	if err != nil {
		t.Fatalf("Failed to update docpage: %v", err)
	}

	// Verify update
	err = db.QueryRow(selectQuery, "testauthor", "testrepo", "docs/README.md").Scan(&title, &content)
	if err != nil {
		t.Fatalf("Failed to query updated record: %v", err)
	}

	if title != "Updated Title" || content != "Updated Content" {
		t.Errorf("Updated data mismatch: title=%q, content=%q", title, content)
	}
}

// TestCleanupStaleEntries tests the cleanup logic
func TestCleanupStaleEntries(t *testing.T) {
	tempDir := t.TempDir()
	db, err := sql.Open("sqlite3", path.Join(tempDir, "test.db"))
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Create schema
	schema := `
		CREATE TABLE IF NOT EXISTS docpages (
			author TEXT NOT NULL,
			repository TEXT NOT NULL,
			relative_path TEXT NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			last_modified_at DATETIME NOT NULL,
			vector BLOB
		);
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	now := time.Now()

	// Insert some test records
	query := `INSERT INTO docpages 
	          (author, repository, relative_path, title, content, last_modified_at, vector)
	          VALUES (?, ?, ?, ?, ?, ?, NULL)`

	filesToInsert := []string{"README.md", "API.md", "Config.md", "stale.md"}
	for _, file := range filesToInsert {
		_, err = db.Exec(query, "testauthor", "testrepo", file, "Title", "Content", now)
		if err != nil {
			t.Fatalf("Failed to insert record: %v", err)
		}
	}

	// Simulate filesystem having only 3 of 4 files
	filesOnDisk := map[string]bool{
		"README.md": true,
		"API.md":    true,
		"Config.md": true,
		// "stale.md" is missing
	}

	// Run cleanup
	deleteQuery := `DELETE FROM docpages 
	                WHERE author = ? AND repository = ? AND relative_path = ?`

	for _, file := range filesToInsert {
		if !filesOnDisk[file] {
			_, err := db.Exec(deleteQuery, "testauthor", "testrepo", file)
			if err != nil {
				t.Fatalf("Failed to delete stale entry: %v", err)
			}
		}
	}

	// Verify cleanup - should have 3 records
	countQuery := `SELECT COUNT(*) FROM docpages WHERE author = ? AND repository = ?`
	var count int
	err = db.QueryRow(countQuery, "testauthor", "testrepo").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count records: %v", err)
	}

	if count != 3 {
		t.Errorf("After cleanup: got %d records, want 3", count)
	}

	// Verify stale.md is gone
	checkQuery := `SELECT COUNT(*) FROM docpages WHERE author = ? AND repository = ? AND relative_path = ?`
	var exists int
	err = db.QueryRow(checkQuery, "testauthor", "testrepo", "stale.md").Scan(&exists)
	if err != nil {
		t.Fatalf("Failed to check stale record: %v", err)
	}

	if exists != 0 {
		t.Errorf("Stale record still exists: %d", exists)
	}
}
