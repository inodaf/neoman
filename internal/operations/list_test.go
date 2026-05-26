package operations

import (
	"database/sql"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB creates an in-memory SQLite database with schema for testing
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

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

	return db
}

// TestQueryAllProjects tests the queryAllProjects function
func TestQueryAllProjects(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	testData := []struct {
		author string
		repo   string
		path   string
		title  string
	}{
		{"n26", "valium", "README.md", "Main"},
		{"n26", "another", "API.md", "API"},
		{"org2", "project1", "Guide.md", "Guide"},
		{"org2", "project2", "Intro.md", "Intro"},
	}

	for _, td := range testData {
		db.Exec(`INSERT INTO docpages 
		         (author, repository, relative_path, title, content, last_modified_at)
		         VALUES (?, ?, ?, ?, ?, datetime('now'))`,
			td.author, td.repo, td.path, td.title, "content")
	}

	projects, err := queryAllProjects(db)
	if err != nil {
		t.Fatalf("queryAllProjects() error = %v", err)
	}

	expected := []string{
		"n26/another",
		"n26/repo",
		"org2/project1",
		"org2/project2",
	}

	if len(projects) != len(expected) {
		t.Errorf("Expected %d projects, got %d", len(expected), len(projects))
	}

	for i, proj := range projects {
		if proj != expected[i] {
			t.Errorf("Project[%d] = %q, want %q", i, proj, expected[i])
		}
	}
}

// TestQueryAllProjectsEmpty tests queryAllProjects with empty database
func TestQueryAllProjectsEmpty(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	projects, err := queryAllProjects(db)
	if err != nil {
		t.Fatalf("queryAllProjects() error = %v", err)
	}

	if len(projects) != 0 {
		t.Errorf("Expected 0 projects, got %d", len(projects))
	}
}

// TestQueryProjectsInOrg tests the queryProjectsInOrg function
func TestQueryProjectsInOrg(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	testData := []struct {
		author string
		repo   string
		path   string
		title  string
	}{
		{"n26", "valium", "README.md", "Main"},
		{"n26", "another", "API.md", "API"},
		{"n26", "zebra", "Guide.md", "Guide"},
		{"org2", "project1", "Intro.md", "Intro"},
	}

	for _, td := range testData {
		db.Exec(`INSERT INTO docpages 
		         (author, repository, relative_path, title, content, last_modified_at)
		         VALUES (?, ?, ?, ?, ?, datetime('now'))`,
			td.author, td.repo, td.path, td.title, "content")
	}

	repos, err := queryProjectsInOrg(db, "n26")
	if err != nil {
		t.Fatalf("queryProjectsInOrg() error = %v", err)
	}

	expected := []string{"another", "valium", "zebra"}

	if len(repos) != len(expected) {
		t.Errorf("Expected %d repos, got %d", len(expected), len(repos))
	}

	for i, repo := range repos {
		if repo != expected[i] {
			t.Errorf("Repo[%d] = %q, want %q", i, repo, expected[i])
		}
	}
}

// TestQueryProjectsInOrgNotFound tests queryProjectsInOrg with non-existent org
func TestQueryProjectsInOrgNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repos, err := queryProjectsInOrg(db, "unknown")
	if err != nil {
		t.Fatalf("queryProjectsInOrg() error = %v", err)
	}

	if len(repos) != 0 {
		t.Errorf("Expected 0 repos, got %d", len(repos))
	}
}

// TestQueryDocsInProject tests the queryDocsInProject function
func TestQueryDocsInProject(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data (unsorted to verify sorting works)
	testData := []struct {
		author string
		repo   string
		path   string
		title  string
	}{
		{"n26", "valium", "Publishing.md", "Publishing"},
		{"n26", "valium", "Teams/Private.md", "Private Repositories"},
		{"n26", "valium", "API.md", "API Guide"},
		{"n26", "other", "README.md", "Other Repo"},
	}

	for _, td := range testData {
		db.Exec(`INSERT INTO docpages 
		         (author, repository, relative_path, title, content, last_modified_at)
		         VALUES (?, ?, ?, ?, ?, datetime('now'))`,
			td.author, td.repo, td.path, td.title, "content")
	}

	docs, err := queryDocsInProject(db, "n26", "valium")
	if err != nil {
		t.Fatalf("queryDocsInProject() error = %v", err)
	}

	expected := []availableDoc{
		{Title: "API Guide", RelativePath: "API.md"},
		{Title: "Private Repositories", RelativePath: "Teams/Private.md"},
		{Title: "Publishing", RelativePath: "Publishing.md"},
	}

	if len(docs) != len(expected) {
		t.Errorf("Expected %d docs, got %d", len(expected), len(docs))
	}

	for i, doc := range docs {
		if doc != expected[i] {
			t.Errorf("Doc[%d] = %+v, want %+v", i, doc, expected[i])
		}
	}
}

// TestQueryDocsInProjectNotFound tests queryDocsInProject with non-existent project
func TestQueryDocsInProjectNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	docs, err := queryDocsInProject(db, "unknown", "unknown")
	if err != nil {
		t.Fatalf("queryDocsInProject() error = %v", err)
	}

	if len(docs) != 0 {
		t.Errorf("Expected 0 docs, got %d", len(docs))
	}
}

// TestFormatAllProjects tests the FormatAllProjects function
func TestFormatAllProjects(t *testing.T) {
	projects := []string{
		"n26/repo",
		"another-org/some-project",
	}

	output := FormatAllProjects(projects)

	// Verify structure
	if !strings.Contains(output, "docs:") {
		t.Error("Missing docs header")
	}
	if !strings.Contains(output, "- n26/repo") {
		t.Error("Missing project in list")
	}
	if !strings.Contains(output, "- another-org/some-project") {
		t.Error("Missing project in list")
	}

	// Verify no prose or separators
	if strings.Contains(output, "This is a list of all documentation installed") {
		t.Error("Should not contain descriptive prose")
	}
	if strings.Contains(output, "---") {
		t.Error("Should not contain separators")
	}
	if strings.Contains(output, "Use 'nman list org/repo'") {
		t.Error("Should not contain next step instruction")
	}
}

// TestFormatProjectsInOrg tests the FormatProjectsInOrg function
func TestFormatProjectsInOrg(t *testing.T) {
	projects := []string{"valium", "another-project"}

	output := FormatProjectsInOrg("n26", projects)

	// Verify structure
	if !strings.Contains(output, "author: n26") {
		t.Error("Missing author header")
	}
	if !strings.Contains(output, "docs:") {
		t.Error("Missing docs header")
	}
	if !strings.Contains(output, "- valium") {
		t.Error("Missing repo in list")
	}
	if !strings.Contains(output, "- another-project") {
		t.Error("Missing repo in list")
	}

	// Verify no prose, bold formatting, or separators
	if strings.Contains(output, "**n26**") {
		t.Error("Should not contain bold formatting")
	}
	if strings.Contains(output, "---") {
		t.Error("Should not contain separators")
	}
	if strings.Contains(output, "`$ nman list n26/repo`") {
		t.Error("Should not contain example command")
	}
}

// TestFormatDocsInProject tests the FormatDocsInProject function
func TestFormatDocsInProject(t *testing.T) {
	docs := []availableDoc{
		{Title: "Publishing", RelativePath: "Authoring Docs/Publishing.md"},
		{Title: "Private Repositories", RelativePath: "Teams/Private Repositories.md"},
	}

	output := FormatDocsInProject("n26", "valium", docs)

	// Verify structure
	if !strings.Contains(output, "project: n26/repo") {
		t.Error("Missing project header")
	}
	if !strings.Contains(output, "docs:") {
		t.Error("Missing docs header")
	}
	if !strings.Contains(output, "- Publishing: 'Authoring Docs/Publishing.md'") {
		t.Error("Missing doc in list")
	}
	if !strings.Contains(output, "- Private Repositories: 'Teams/Private Repositories.md'") {
		t.Error("Missing doc in list")
	}

	// Verify no prose, bold formatting, comments, or separators
	if strings.Contains(output, "**n26/repo**") {
		t.Error("Should not contain bold formatting")
	}
	if strings.Contains(output, "---") {
		t.Error("Should not contain separators")
	}
	if strings.Contains(output, "# formatted as \"Title: Path\"") {
		t.Error("Should not contain format comment")
	}
	if strings.Contains(output, "`$ nman view n26/repo \"Authoring Docs/Publishing.md\"`") {
		t.Error("Should not contain example command")
	}
}

// TestFormatDocsInProjectEmpty tests FormatDocsInProject with empty docs
func TestFormatDocsInProjectEmpty(t *testing.T) {
	docs := []availableDoc{}

	output := FormatDocsInProject("n26", "valium", docs)

	// Should still have structure but no docs
	if !strings.Contains(output, "project: n26/repo") {
		t.Error("Missing project header")
	}
	if !strings.Contains(output, "docs:") {
		t.Error("Missing docs header")
	}
	if strings.Contains(output, "---") {
		t.Error("Should not contain separators")
	}
}

// TestFormatErrorNoProjects tests error message for no projects
func TestFormatErrorNoProjects(t *testing.T) {
	output := FormatErrorNoProjects()

	if !strings.Contains(output, "neoman:") {
		t.Error("Missing neoman prefix")
	}
	if !strings.Contains(output, "No documentation found") {
		t.Error("Missing error message")
	}
}

// TestFormatErrorOrgNotFound tests error message for org not found
func TestFormatErrorOrgNotFound(t *testing.T) {
	output := FormatErrorOrgNotFound("unknown")

	if !strings.Contains(output, "neoman:") {
		t.Error("Missing neoman prefix")
	}
	if !strings.Contains(output, "No documentation found for organization \"unknown\"") {
		t.Error("Missing error message with org name")
	}
}

// TestFormatErrorProjectNotFound tests error message for project not found
func TestFormatErrorProjectNotFound(t *testing.T) {
	output := FormatErrorProjectNotFound("n26", "unknown")

	if !strings.Contains(output, "neoman:") {
		t.Error("Missing neoman prefix")
	}
	if !strings.Contains(output, "No documentation found for project \"n26/unknown\"") {
		t.Error("Missing error message with project name")
	}
}

// TestListAllProjects integration test
func TestListAllProjects(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	db.Exec(`INSERT INTO docpages 
	         (author, repository, relative_path, title, content, last_modified_at)
	         VALUES (?, ?, ?, ?, ?, datetime('now'))`,
		"n26", "valium", "README.md", "Main", "content")

	projects, err := ListAllProjects(db)
	if err != nil {
		t.Fatalf("ListAllProjects() error = %v", err)
	}

	if len(projects) != 1 {
		t.Errorf("Expected 1 project, got %d", len(projects))
	}

	if projects[0] != "n26/repo" {
		t.Errorf("Expected 'n26/repo', got %q", projects[0])
	}
}

// TestListProjectsInOrg integration test
func TestListProjectsInOrg(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	db.Exec(`INSERT INTO docpages 
	         (author, repository, relative_path, title, content, last_modified_at)
	         VALUES (?, ?, ?, ?, ?, datetime('now'))`,
		"n26", "valium", "README.md", "Main", "content")

	repos, err := ListProjectsInOrg(db, "n26")
	if err != nil {
		t.Fatalf("ListProjectsInOrg() error = %v", err)
	}

	if len(repos) != 1 {
		t.Errorf("Expected 1 repo, got %d", len(repos))
	}

	if repos[0] != "valium" {
		t.Errorf("Expected 'valium', got %q", repos[0])
	}
}

// TestListDocsInProject integration test
func TestListDocsInProject(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	db.Exec(`INSERT INTO docpages 
	         (author, repository, relative_path, title, content, last_modified_at)
	         VALUES (?, ?, ?, ?, ?, datetime('now'))`,
		"n26", "valium", "Publishing.md", "Publishing", "content")

	docs, err := ListDocsInProject(db, "n26", "valium")
	if err != nil {
		t.Fatalf("ListDocsInProject() error = %v", err)
	}

	if len(docs) != 1 {
		t.Errorf("Expected 1 doc, got %d", len(docs))
	}

	if docs[0].Title != "Publishing" {
		t.Errorf("Expected title 'Publishing', got %q", docs[0].Title)
	}

	if docs[0].RelativePath != "Publishing.md" {
		t.Errorf("Expected path 'Publishing.md', got %q", docs[0].RelativePath)
	}
}

// TestListProjectsInOrgEmptyOrg tests error case for empty org
func TestListProjectsInOrgEmptyOrg(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, err := ListProjectsInOrg(db, "")
	if err == nil {
		t.Error("Expected error for empty org, got nil")
	}
}

// TestListDocsInProjectEmptyAuthor tests error case for empty author
func TestListDocsInProjectEmptyAuthor(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, err := ListDocsInProject(db, "", "repo")
	if err == nil {
		t.Error("Expected error for empty author, got nil")
	}
}

// TestListDocsInProjectEmptyRepo tests error case for empty repo
func TestListDocsInProjectEmptyRepo(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, err := ListDocsInProject(db, "org", "")
	if err == nil {
		t.Error("Expected error for empty repo, got nil")
	}
}
