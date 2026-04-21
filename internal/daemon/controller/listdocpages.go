package controller

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/inodaf/neoman/internal/daemon/domain"
	"github.com/inodaf/neoman/internal/daemon/usecase"
)

func (c *controller) ListDocPages(w http.ResponseWriter, r *http.Request) {
	author := r.PathValue("author")
	repo := r.PathValue("repo")

	accept := r.Header.Get("Accept")
	if accept != "" && accept != "application/yaml" && accept != "*/*" {
		http.Error(w, "Unsupported media type. Use Accept: application/yaml", http.StatusNotAcceptable)
		return
	}

	output, err := c.useCase.ListDocPages(usecase.ListDocPagesInput{
		Author:     author,
		Repository: repo,
	})

	switch err {
	case nil:
		if len(output.Pages) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "application/yaml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(formatPagesAsYAML(author, repo, output.Pages)))
		return
	case usecase.ErrListDocPagesDocsNotFound:
		http.Error(w, "Docs not found", http.StatusNotFound)
		return
	case usecase.ErrListDocPagesAuthorRequired, usecase.ErrListDocPagesRepositoryRequired:
		http.Error(w, "Missing author or repository", http.StatusBadRequest)
		return
	default:
		slog.Error("failed to list doc pages", "error", err)
		http.Error(w, "Unable to list doc pages", http.StatusInternalServerError)
		return
	}
}

func formatPagesAsYAML(author, repo string, pages []domain.DocsPage) string {
	var yaml strings.Builder

	fmt.Fprintf(&yaml, "documentation: %s/%s\n", author, repo)
	yaml.WriteString("pages:\n")

	for _, page := range pages {
		fmt.Fprintf(&yaml, "  - %s: '%s'\n", page.Title, page.RelativePath)
	}

	return yaml.String()
}
