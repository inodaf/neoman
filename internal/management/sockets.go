package management

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/inodaf/neoman/pkg/config"
)

var UnixSockClient http.Client = http.Client{
	Transport: &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial("unix", config.AppSockPath)
		},
	},
}

// SocketClientPing tries to ping the daemon through the Unix socket.
// If the daemon is not running, it returns an error.
func SocketClientPing() error {
	resource := url.URL{Host: "unix", Scheme: "http", Path: "/ping"}
	_, err := UnixSockClient.Get(resource.String())

	return err
}

// SocketServeIPC serves the IPC Unix Domain Socket (UDS) for communication
// between nman and the nmand (daemon).
func SocketServeIPC(db *sql.DB) {
	if err := os.RemoveAll(config.AppSockPath); err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("unix", config.AppSockPath)
	if err != nil {
		log.Fatalln("neoman: could not listen to the socket", err)
	}

	defer listener.Close()
	log.Println("Listening to socket at", config.AppSockPath)

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("pong"))
	})

	// List command handlers
	// GET /list - List all projects
	mux.HandleFunc("GET /list", func(w http.ResponseWriter, r *http.Request) {
		response := handleListAllProjects(db)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(response))
	})

	// GET /list/{org} - List projects in organization
	mux.HandleFunc("GET /list/{org}", func(w http.ResponseWriter, r *http.Request) {
		org := r.PathValue("org")
		response := handleListProjectsInOrg(db, org)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(response))
	})

	// GET /list/{owner}/{repo} - List documents in project
	mux.HandleFunc("GET /list/{owner}/{repo}", func(w http.ResponseWriter, r *http.Request) {
		owner := r.PathValue("owner")
		repo := r.PathValue("repo")
		response := handleListDocsInProject(db, owner, repo)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(response))
	})

	// View command handlers
	// GET /view/{owner}/{repo}/* - View document content
	mux.HandleFunc("GET /view/{owner}/{repo}/{path...}", func(w http.ResponseWriter, r *http.Request) {
		owner := r.PathValue("owner")
		repo := r.PathValue("repo")
		pathParam := r.PathValue("path")
		handleViewDocument(w, db, owner, repo, pathParam)
	})

	if err := http.Serve(listener, mux); err != nil {
		log.Fatalln("neoman: could not serve IPC Unix Socket", err)
	}
}

// handleListAllProjects is the internal wrapper for list all projects
// We define it here to avoid circular imports
func handleListAllProjects(db *sql.DB) string {
	// Query all projects
	query := `SELECT DISTINCT author, repository FROM docpages ORDER BY author, repository`
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Sprintf("neoman: database error: %v", err)
	}
	defer rows.Close()

	var projects []string
	for rows.Next() {
		var author, repo string
		if err := rows.Scan(&author, &repo); err != nil {
			return fmt.Sprintf("neoman: database error: %v", err)
		}
		projects = append(projects, fmt.Sprintf("%s/%s", author, repo))
	}

	if len(projects) == 0 {
		return "neoman: No documentation found. Download your first documentation with: nman owner/repo"
	}

	var output strings.Builder
	output.WriteString("docs:\n")
	for _, project := range projects {
		output.WriteString(fmt.Sprintf("  - %s\n", project))
	}

	return output.String()
}

// handleListProjectsInOrg is the internal wrapper for list org projects
func handleListProjectsInOrg(db *sql.DB, org string) string {
	if org == "" {
		return "neoman: organization name cannot be empty"
	}

	// Query projects in org
	query := `SELECT DISTINCT repository FROM docpages WHERE author = ? ORDER BY repository`
	rows, err := db.Query(query, org)
	if err != nil {
		return fmt.Sprintf("neoman: database error: %v", err)
	}
	defer rows.Close()

	var repos []string
	for rows.Next() {
		var repo string
		if err := rows.Scan(&repo); err != nil {
			return fmt.Sprintf("neoman: database error: %v", err)
		}
		repos = append(repos, repo)
	}

	if len(repos) == 0 {
		return fmt.Sprintf("neoman: No documentation found for organization \"%s\"", org)
	}

	var output strings.Builder
	output.WriteString(fmt.Sprintf("author: %s\n", org))
	output.WriteString("docs:\n")
	for _, repo := range repos {
		output.WriteString(fmt.Sprintf("  - %s\n", repo))
	}

	return output.String()
}

// handleListDocsInProject is the internal wrapper for list project docs
func handleListDocsInProject(db *sql.DB, author, repo string) string {
	if author == "" {
		return "neoman: author cannot be empty"
	}
	if repo == "" {
		return "neoman: repository cannot be empty"
	}

	// Query docs in project
	query := `SELECT title, relative_path FROM docpages WHERE author = ? AND repository = ? ORDER BY title`
	rows, err := db.Query(query, author, repo)
	if err != nil {
		return fmt.Sprintf("neoman: database error: %v", err)
	}
	defer rows.Close()

	type availableDoc struct {
		Title        string
		RelativePath string
	}

	var docs []availableDoc
	for rows.Next() {
		var title, relativePath string
		if err := rows.Scan(&title, &relativePath); err != nil {
			return fmt.Sprintf("neoman: database error: %v", err)
		}
		docs = append(docs, availableDoc{Title: title, RelativePath: relativePath})
	}

	if len(docs) == 0 {
		return fmt.Sprintf("neoman: No documentation found for project \"%s/%s\"", author, repo)
	}

	var output strings.Builder
	output.WriteString(fmt.Sprintf("project: %s/%s\n", author, repo))
	output.WriteString("docs:\n")
	for _, doc := range docs {
		output.WriteString(fmt.Sprintf("  - %s: '%s'\n", doc.Title, doc.RelativePath))
	}

	return output.String()
}

// handleViewDocument is the handler for viewing document content
func handleViewDocument(w http.ResponseWriter, db *sql.DB, owner, repo, pathParam string) {
	// URL-decode the path parameter
	decodedPath, err := url.QueryUnescape(pathParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid path encoding"})
		return
	}

	// Import operations package to call QueryDocumentContent
	// We need to call operations.QueryDocumentContent, but operations is imported at package level
	// So we'll use the internal helper directly here
	query := `
		SELECT content, title FROM docpages 
		WHERE author = ? AND repository = ? 
		  AND LOWER(relative_path) = LOWER(?)
		LIMIT 1
	`

	var content, title string
	row := db.QueryRow(query, owner, repo, decodedPath)
	err = row.Scan(&content, &title)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "document not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("database error: %v", err)})
		return
	}

	// Return successful response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"content": content,
		"title":   title,
	})
}

// Handle
// - neoman.local/:owner
// - neoman.local/:localRepo
// - neoman.local/:owner/:repo/
// - neoman.local/:owner/:repo/:[introduction/getting-started]
// - neoman.local/?q=MyQuery
// - neoman.local/
// - -
//
// edge cases
// owner and repo sanitized + URL escaped/unescaped
// owner is not found
// repo is not found (fetch it)
// default .md file to display (index.md) or fallback to root's README.md
// -- CLI (ask which page want to open by default for that repo)
// what happens if owner name is same as a local repo?
func SocketServeTCP(db *sql.DB) {
	mux := http.NewServeMux()
	log.Println("Listening to web at", config.AppWebAppPort)

	mux.HandleFunc("GET /{localRepo}", func(w http.ResponseWriter, r *http.Request) {
		repo := r.PathValue("localRepo")

		dir, err := config.DocsRegistryDir()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.ServeFile(w, r, path.Join(dir, "local", repo, "docs", "README.md"))
	})

	mux.HandleFunc("GET /{owner}/{repo}", func(w http.ResponseWriter, r *http.Request) {
		owner := r.PathValue("owner")
		repo := r.PathValue("repo")

		dir, err := config.DocsRegistryDir()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.ServeFile(w, r, path.Join(dir, "remote", owner, repo, "README.md"))
	})

	if err := http.ListenAndServe(config.AppWebAppPort, mux); err != nil {
		log.Fatalln("neoman: could not serve web application", err)
	}
}
