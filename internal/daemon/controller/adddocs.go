package controller

import (
	"log/slog"
	"net/http"

	"github.com/inodaf/neoman/internal/daemon/usecase"
)

func (c *controller) AddDocs(w http.ResponseWriter, r *http.Request) {
	author := r.PathValue("author")
	repo := r.PathValue("repo")

	err := c.useCase.AddRemoteDocs(&usecase.AddDocsInput{
		Author:     author,
		Repository: repo,
	})
	if err != nil {
		switch err {
		case usecase.ErrAuthorRequired, usecase.ErrRepositoryRequired:
			http.Error(w, "Missing author or repository", http.StatusBadRequest)
			return
		case usecase.ErrDocsDirNotFound:
			http.Error(w, "Repository does not have a 'docs/' directory", http.StatusNotFound)
			return
		case usecase.ErrDocsAlreadyExist:
			break
		default:
			slog.Error("failed to add docs", "error", err)
			http.Error(w, "Unable to add docs", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("docs added"))
}
