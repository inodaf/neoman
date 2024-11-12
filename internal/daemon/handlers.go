package daemon

import (
	"database/sql"
	"log"
	"net/http"
)

type Handlers struct {
	DB *sql.DB
}

func (h *Handlers) Pong(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("pong"))
}

func (h *Handlers) CheckTrustedRemoteAccount(w http.ResponseWriter, r *http.Request) {
	account := r.URL.Query().Get("account")
	log.Println("Handlers/IsTrusted", account)

	if account == "inodaf" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

	w.Write(nil)
}
