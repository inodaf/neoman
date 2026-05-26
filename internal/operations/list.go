package operations

import (
	"database/sql"
	"fmt"
	"strings"
)

type availableDoc struct {
	Title        string
	RelativePath string
}

// ListAllProjects retrieves all documentation projects in the system
func ListAllProjects(db *sql.DB) ([]string, error) {
	return queryAllProjects(db)
}

// ListProjectsInOrg retrieves all projects in a specific organization
func ListProjectsInOrg(db *sql.DB, org string) ([]string, error) {
	if org == "" {
		return nil, fmt.Errorf("organization name cannot be empty")
	}
	return queryProjectsInOrg(db, org)
}

// ListDocsInProject retrieves all documents in a specific project
func ListDocsInProject(db *sql.DB, author, repo string) ([]availableDoc, error) {
	if author == "" {
		return nil, fmt.Errorf("author cannot be empty")
	}
	if repo == "" {
		return nil, fmt.Errorf("repository cannot be empty")
	}
	return queryDocsInProject(db, author, repo)
}

// HandleListAllProjects handles the /list endpoint
func HandleListAllProjects(db *sql.DB) string {
	projects, err := ListAllProjects(db)
	if err != nil {
		return fmt.Sprintf("neoman: database error: %v", err)
	}

	if len(projects) == 0 {
		return FormatErrorNoProjects()
	}

	return FormatAllProjects(projects)
}

// HandleListProjectsInOrg handles the /list/{org}/* endpoint
func HandleListProjectsInOrg(db *sql.DB, org string) string {
	if org == "" {
		return "neoman: organization name cannot be empty"
	}

	projects, err := ListProjectsInOrg(db, org)
	if err != nil {
		return fmt.Sprintf("neoman: database error: %v", err)
	}

	if len(projects) == 0 {
		return FormatErrorOrgNotFound(org)
	}

	return FormatProjectsInOrg(org, projects)
}

// HandleListDocsInProject handles the /list/{owner}/{repo} endpoint
func HandleListDocsInProject(db *sql.DB, author, repo string) string {
	if author == "" {
		return "neoman: author cannot be empty"
	}
	if repo == "" {
		return "neoman: repository cannot be empty"
	}

	docs, err := ListDocsInProject(db, author, repo)
	if err != nil {
		return fmt.Sprintf("neoman: database error: %v", err)
	}

	if len(docs) == 0 {
		return FormatErrorProjectNotFound(author, repo)
	}

	return FormatDocsInProject(author, repo, docs)
}

// queryAllProjects retrieves distinct (author, repository) pairs sorted alphabetically
func queryAllProjects(db *sql.DB) ([]string, error) {
	query := `SELECT DISTINCT author, repository FROM docpages ORDER BY author, repository`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all projects: %w", err)
	}
	defer rows.Close()

	var projects []string
	for rows.Next() {
		var author, repo string
		if err := rows.Scan(&author, &repo); err != nil {
			return nil, fmt.Errorf("failed to scan project row: %w", err)
		}
		projects = append(projects, fmt.Sprintf("%s/%s", author, repo))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating projects: %w", err)
	}

	return projects, nil
}

// queryProjectsInOrg retrieves repositories in a specific organization, sorted alphabetically
func queryProjectsInOrg(db *sql.DB, org string) ([]string, error) {
	query := `SELECT DISTINCT repository FROM docpages WHERE author = ? ORDER BY repository`

	rows, err := db.Query(query, org)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects in organization: %w", err)
	}
	defer rows.Close()

	var repos []string
	for rows.Next() {
		var repo string
		if err := rows.Scan(&repo); err != nil {
			return nil, fmt.Errorf("failed to scan repository row: %w", err)
		}
		repos = append(repos, repo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating repositories: %w", err)
	}

	return repos, nil
}

// queryDocsInProject retrieves all documents in a project, sorted by title
func queryDocsInProject(db *sql.DB, author, repo string) ([]availableDoc, error) {
	query := `SELECT title, relative_path FROM docpages WHERE author = ? AND repository = ? ORDER BY title`

	rows, err := db.Query(query, author, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to query documents: %w", err)
	}
	defer rows.Close()

	var docs []availableDoc
	for rows.Next() {
		var title, relativePath string
		if err := rows.Scan(&title, &relativePath); err != nil {
			return nil, fmt.Errorf("failed to scan document row: %w", err)
		}
		docs = append(docs, availableDoc{Title: title, RelativePath: relativePath})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating documents: %w", err)
	}

	return docs, nil
}

// FormatAllProjects formats all projects for display
func FormatAllProjects(projects []string) string {
	var output strings.Builder

	output.WriteString("docs:\n")

	for _, project := range projects {
		output.WriteString(fmt.Sprintf("  - %s\n", project))
	}

	return output.String()
}

// FormatProjectsInOrg formats organization projects for display
func FormatProjectsInOrg(org string, projects []string) string {
	var output strings.Builder

	output.WriteString(fmt.Sprintf("author: %s\n", org))
	output.WriteString("docs:\n")

	for _, project := range projects {
		output.WriteString(fmt.Sprintf("  - %s\n", project))
	}

	return output.String()
}

// FormatDocsInProject formats project documents for display
func FormatDocsInProject(author, repo string, docs []availableDoc) string {
	var output strings.Builder

	output.WriteString(fmt.Sprintf("project: %s/%s\n", author, repo))
	output.WriteString("docs:\n")

	for _, doc := range docs {
		output.WriteString(fmt.Sprintf("  - %s: '%s'\n", doc.Title, doc.RelativePath))
	}

	return output.String()
}

// FormatErrorNoProjects returns error message when no projects found
func FormatErrorNoProjects() string {
	return "neoman: No documentation found. Download your first documentation with: nman owner/repo"
}

// FormatErrorOrgNotFound returns error message when organization not found
func FormatErrorOrgNotFound(org string) string {
	return fmt.Sprintf("neoman: No documentation found for organization \"%s\"", org)
}

// FormatErrorProjectNotFound returns error message when project not found
func FormatErrorProjectNotFound(author, repo string) string {
	return fmt.Sprintf("neoman: No documentation found for project \"%s/%s\"", author, repo)
}
