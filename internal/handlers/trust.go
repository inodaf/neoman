package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
)

// TrustHandler handles the Trusted Account scope. It is responsible for
// simply checking and adding new Git remote accounts as trusted so that subsequent
// Git operations for the account is safely verified.
type TrustHandler struct {
	DB *sql.DB
}

// Verify checks if the given [account] name is
// present in the database. Returns false when the account was not marked as trusted.
func (h *TrustHandler) Verify(account string) (bool, error) {
	if len(account) == 0 {
		return false, errors.New("Verify: arg 'account' not provided")
	}
	var count int
	err := h.DB.QueryRow("SELECT COUNT(*) FROM TrustedAccount WHERE account = ?", account).Scan(&count)
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

// Trust flags the provided [account] as trusted and
// does nothing if the given account is trusted already.
func (h *TrustHandler) Trust(account string) error {
	if len(account) == 0 {
		return errors.New("Trust: arg 'account' not provided")
	}
	ok, err := h.Verify(account)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	_, err = h.DB.Exec(`INSERT INTO TrustedAccount (account) VALUES (?)`, account)
	return err
}

// TrustHandler implements [http.Handler]
func (h TrustHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	account := r.URL.Query().Get("account")
	switch r.Method {
	case http.MethodGet:
		log.Printf("GET /trust '%s'", account)
		ok, err := h.Verify(account)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		if ok {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
		w.Write(nil)
	case http.MethodPost:
		log.Printf("POST /trust '%s'", account)
		err := h.Trust(account)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
		w.Write(nil)
	}
}
