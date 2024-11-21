package handlers

import (
	"database/sql"

	"github.com/inodaf/neoman/internal/git"
)

type RemoteHandler struct {
	DB             *sql.DB
	RemoteProvider git.RemoteProvider
}

func (h *RemoteHandler) Fetch(owner, repo string) {
	h.RemoteProvider.DocsDirExists(owner, repo)
}
