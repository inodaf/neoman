package daemon

import (
	"database/sql"
	"net/http"
)

type Handlers struct {
	DB *sql.DB
}

func (h *Handlers) Pong(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("pong"))
}
