package operations

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/inodaf/neoman/internal/management"
)

// ViewDocument retrieves and displays a specific document from the database
// with metadata header and raw markdown content
func ViewDocument(project, documentPath string) (string, error) {
	// Parse project to extract owner and repo
	var owner, repo string

	if project == "." {
		// Resolve local project from current directory name
		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("could not get current working directory")
		}
		repo = path.Base(wd)
		owner = ""
		project = repo // Update project for output formatting
	} else {
		// Parse owner/repo format
		parts := strings.Split(project, "/")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid project format. Use 'owner/repo' or '.'")
		}
		owner = parts[0]
		repo = parts[1]
	}

	// Validate repo is not empty
	if repo == "" {
		return "", fmt.Errorf("repository name cannot be empty")
	}

	// URL-encode the document path to handle spaces and special characters
	encodedPath := url.PathEscape(documentPath)

	// Build socket request URL
	var endpoint string
	if owner == "" {
		// Local project
		endpoint = fmt.Sprintf("/view/%s/%s", repo, encodedPath)
	} else {
		endpoint = fmt.Sprintf("/view/%s/%s/%s", owner, repo, encodedPath)
	}

	// Make HTTP GET request to daemon via socket
	resource := url.URL{Host: "unix", Scheme: "http", Path: endpoint}
	resp, err := management.UnixSockClient.Get(resource.String())
	if err != nil {
		return "", fmt.Errorf("could not connect to daemon: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not read daemon response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("document not found")
	}

	// Parse JSON response
	var response struct {
		Content string `json:"content"`
		Title   string `json:"title"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("could not parse daemon response: %w", err)
	}

	// Format output with YAML-like header
	var output strings.Builder
	output.WriteString("---\n")
	output.WriteString(fmt.Sprintf("project: %s\n", project))
	output.WriteString(fmt.Sprintf("title: %s\n", response.Title))

	// Normalize path to always start with /
	normalizedPath := documentPath
	if !strings.HasPrefix(normalizedPath, "/") {
		normalizedPath = "/" + normalizedPath
	}
	output.WriteString(fmt.Sprintf("path: %s\n", normalizedPath))
	output.WriteString("---\n\n")
	output.WriteString(response.Content)

	return output.String(), nil
}
