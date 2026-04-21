package controller

import (
	"log/slog"
	"net/http"
	"net/url"

	"github.com/inodaf/neoman/internal/daemon/usecase"
)

func (c *controller) ViewDocPage(w http.ResponseWriter, r *http.Request) {
	author := r.PathValue("author")
	repo := r.PathValue("repo")
	path := r.PathValue("path")

	decodedPath, err := url.PathUnescape(path)
	if err != nil {
		http.Error(w, "Invalid path encoding", http.StatusBadRequest)
		return
	}

	output, err := c.useCase.ViewDocPage(usecase.ViewDocPageInput{
		Author:       author,
		Repository:   repo,
		RelativePath: decodedPath,
	})

	switch err {
	case nil:
		w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(output.Content))
		return
	case usecase.ErrViewDocPageDocsNotFound:
		http.Error(w, "Documentation not found", http.StatusNotFound)
		return
	case usecase.ErrViewDocPageNotFound:
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	case usecase.ErrViewDocPageAuthorRequired, usecase.ErrViewDocPageRepositoryRequired, usecase.ErrViewDocPagePathRequired:
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	default:
		slog.Error("failed to view doc page", "error", err)
		http.Error(w, "Unable to view page", http.StatusInternalServerError)
		return
	}
}
