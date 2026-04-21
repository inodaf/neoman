package controller

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/inodaf/neoman/internal/daemon/usecase"
)

func (c *controller) ListAllDocs(w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Get("Accept")
	if accept != "" && accept != "application/yaml" && accept != "*/*" {
		http.Error(w, "Unsupported media type. Use Accept: application/yaml", http.StatusNotAcceptable)
		return
	}

	output, err := c.useCase.ListAllDocs()

	switch err {
	case nil:
		w.Header().Set("Content-Type", "application/yaml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(formatAllDocsAsYAML(output.Documentations)))
		return
	case usecase.ErrListAllDocsEmpty:
		http.Error(w, "No documentation found", http.StatusNotFound)
		return
	default:
		slog.Error("failed to list all docs", "error", err)
		http.Error(w, "Unable to list all docs", http.StatusInternalServerError)
		return
	}
}

func formatAllDocsAsYAML(docs []string) string {
	var yaml strings.Builder

	yaml.WriteString("docs:\n")

	for _, doc := range docs {
		fmt.Fprintf(&yaml, "  - %s\n", doc)
	}

	return yaml.String()
}
