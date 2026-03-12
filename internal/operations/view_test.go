package operations

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// TestQueryDocumentContent_Success tests QueryDocumentContent with valid document
func TestQueryDocumentContent_Success(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	db.Exec(`INSERT INTO docpages 
	         (author, repository, relative_path, title, content, last_modified_at)
	         VALUES (?, ?, ?, ?, ?, datetime('now'))`,
		"n26", "repo", "Authoring Docs/4. Publishing.md", "Publishing", "# Publishing\n\nContent here")

	content, title, err := QueryDocumentContent(db, "n26", "repo", "Authoring Docs/4. Publishing.md")
	if err != nil {
		t.Fatalf("QueryDocumentContent() error = %v", err)
	}

	if title != "Publishing" {
		t.Errorf("title = %q, want %q", title, "Publishing")
	}

	if content != "# Publishing\n\nContent here" {
		t.Errorf("content = %q, want %q", content, "# Publishing\n\nContent here")
	}
}

// TestQueryDocumentContent_CaseInsensitive tests case-insensitive path matching
func TestQueryDocumentContent_CaseInsensitive(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data with mixed case
	db.Exec(`INSERT INTO docpages 
	         (author, repository, relative_path, title, content, last_modified_at)
	         VALUES (?, ?, ?, ?, ?, datetime('now'))`,
		"n26", "repo", "Authoring Docs/4. Publishing.md", "Publishing", "# Publishing\n\nContent")

	// Query with different case
	content, title, err := QueryDocumentContent(db, "n26", "repo", "authoring docs/4. publishing.md")
	if err != nil {
		t.Fatalf("QueryDocumentContent() error = %v", err)
	}

	if title != "Publishing" {
		t.Errorf("title = %q, want %q", title, "Publishing")
	}

	expectedContent := "# Publishing\n\nContent"
	if content != expectedContent {
		t.Errorf("content mismatch")
	}
}

// TestQueryDocumentContent_NotFound tests QueryDocumentContent with non-existent document
func TestQueryDocumentContent_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, _, err := QueryDocumentContent(db, "n26", "repo", "nonexistent.md")
	if err == nil {
		t.Fatalf("QueryDocumentContent() expected error, got nil")
	}
}

// TestQueryDocumentContent_WrongProject tests QueryDocumentContent with wrong project
func TestQueryDocumentContent_WrongProject(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	db.Exec(`INSERT INTO docpages 
	         (author, repository, relative_path, title, content, last_modified_at)
	         VALUES (?, ?, ?, ?, ?, datetime('now'))`,
		"n26", "repo1", "test.md", "Test", "content")

	// Query different project
	_, _, err := QueryDocumentContent(db, "n26", "repo2", "test.md")
	if err == nil {
		t.Fatalf("QueryDocumentContent() expected error, got nil")
	}
}

// TestQueryDocumentContent_WrongAuthor tests QueryDocumentContent with wrong author
func TestQueryDocumentContent_WrongAuthor(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	db.Exec(`INSERT INTO docpages 
	         (author, repository, relative_path, title, content, last_modified_at)
	         VALUES (?, ?, ?, ?, ?, datetime('now'))`,
		"n26", "repo", "test.md", "Test", "content")

	// Query different author
	_, _, err := QueryDocumentContent(db, "different", "repo", "test.md")
	if err == nil {
		t.Fatalf("QueryDocumentContent() expected error, got nil")
	}
}

// TestQueryDocumentContent_MultipleDocuments tests with multiple documents
func TestQueryDocumentContent_MultipleDocuments(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert multiple documents
	docs := []struct {
		path    string
		title   string
		content string
	}{
		{"README.md", "Welcome", "# Welcome\n\nMain doc"},
		{"API.md", "API Guide", "# API\n\nAPI doc"},
		{"Getting Started.md", "Getting Started", "# Getting Started\n\nGuide"},
	}

	for _, doc := range docs {
		db.Exec(`INSERT INTO docpages 
		         (author, repository, relative_path, title, content, last_modified_at)
		         VALUES (?, ?, ?, ?, ?, datetime('now'))`,
			"myorg", "myrepo", doc.path, doc.title, doc.content)
	}

	// Query first document
	content1, title1, err := QueryDocumentContent(db, "myorg", "myrepo", "README.md")
	if err != nil {
		t.Fatalf("QueryDocumentContent() error = %v", err)
	}
	if title1 != "Welcome" {
		t.Errorf("title = %q, want %q", title1, "Welcome")
	}

	// Query second document
	_, title2, err := QueryDocumentContent(db, "myorg", "myrepo", "API.md")
	if err != nil {
		t.Fatalf("QueryDocumentContent() error = %v", err)
	}
	if title2 != "API Guide" {
		t.Errorf("title = %q, want %q", title2, "API Guide")
	}

	// Query with spaces in path
	content3, title3, err := QueryDocumentContent(db, "myorg", "myrepo", "Getting Started.md")
	if err != nil {
		t.Fatalf("QueryDocumentContent() error = %v", err)
	}
	if title3 != "Getting Started" {
		t.Errorf("title = %q, want %q", title3, "Getting Started")
	}

	expectedContent1 := "# Welcome\n\nMain doc"
	if content1 != expectedContent1 {
		t.Errorf("content mismatch")
	}

	expectedContent3 := "# Getting Started\n\nGuide"
	if content3 != expectedContent3 {
		t.Errorf("content mismatch")
	}
}
