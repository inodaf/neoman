package controller

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/inodaf/neoman/internal/daemon/usecase"
)

func (c *controller) ListAuthorDocs(w http.ResponseWriter, r *http.Request) {
	author := r.PathValue("author")

	accept := r.Header.Get("Accept")
	if accept != "" && accept != "application/yaml" && accept != "*/*" {
		http.Error(w, "Unsupported media type. Use Accept: application/yaml", http.StatusNotAcceptable)
		return
	}

	output, err := c.useCase.ListAuthorDocs(usecase.ListAuthorDocsInput{
		Author: author,
	})

	switch err {
	case nil:
		w.Header().Set("Content-Type", "application/yaml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(formatAuthorDocsAsYAML(author, output.Documentations)))
		return
	case usecase.ErrListAuthorDocsAuthorRequired:
		http.Error(w, "Author is required", http.StatusBadRequest)
		return
	case usecase.ErrListAuthorDocsNotFound:
		http.Error(w, "No documentation found for author", http.StatusNotFound)
		return
	default:
		slog.Error("failed to list author docs", "error", err)
		http.Error(w, "Unable to list author docs", http.StatusInternalServerError)
		return
	}
}

func formatAuthorDocsAsYAML(author string, docs []string) string {
	var yaml strings.Builder

	fmt.Fprintf(&yaml, "author: %s\n", author)
	yaml.WriteString("docs:\n")

	for _, doc := range docs {
		fmt.Fprintf(&yaml, "  - %s\n", doc)
	}

	return yaml.String()
}
