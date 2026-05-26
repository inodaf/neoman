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
	case nil:
		w.WriteHeader(http.StatusAccepted)
		return
	case usecase.ErrAddRemoteDocsAlreadyExist:
		http.Error(w, "Docs already exist", http.StatusConflict)
		return
	case usecase.ErrAddRemoteDocsAuthorRequired, usecase.ErrAddRemoteDocsRepositoryRequired:
		http.Error(w, "Missing author or repository", http.StatusBadRequest)
		return
	case usecase.ErrAddRemoteDocsDirNotFound:
		http.Error(w, "Repo does not have a 'docs/' directory", http.StatusNoContent)
		return
	default:
		slog.Error("failed to add docs", "error", err)
		http.Error(w, "Unable to add docs", http.StatusInternalServerError)
		return
	}
}
