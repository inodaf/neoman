package controller

import (
	"log/slog"
	"net/http"

	"github.com/inodaf/neoman/internal/daemon/usecase"
)

func (c *controller) AddDocs(w http.ResponseWriter, r *http.Request) {
	author := r.PathValue("author")
	repo := r.PathValue("repo")

	err := c.useCase.AddRemoteDocs(usecase.AddDocsInput{
		Author:     author,
		Repository: repo,
	})

	switch err {
	case usecase.ErrAddRemoteDocsAlreadyExist, nil:
		// Ignored: Docs already exist or no error
		break
	case usecase.ErrAddRemoteDocsAuthorRequired, usecase.ErrAddRemoteDocsRepositoryRequired:
		http.Error(w, "Missing author or repository", http.StatusBadRequest)
		return
	case usecase.ErrAddRemoteDocsDirNotFound:
		http.Error(w, "Repo does not have a 'docs/' directory", http.StatusNotFound)
		return
	default:
		slog.Error("failed to add docs", "error", err)
		http.Error(w, "Unable to add docs", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
