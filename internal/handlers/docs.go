package handlers

import (
	"database/sql"
	"net/http"
)

type DocsHandler struct {
	DB *sql.DB
}

func (h *DocsHandler) Save() {

}

func (h DocsHandler) ServeHTTP(w http.ResponseWriter, r http.Request) {

}
