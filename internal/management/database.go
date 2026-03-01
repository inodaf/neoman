package management

import (
	"database/sql"
	"log/slog"
	"path"

	"github.com/inodaf/neoman/pkg/config"
	_ "github.com/mattn/go-sqlite3"
)

// NewSQLiteDatabase creates a connection to a SQL DB
// and prepares the base schema.
func NewSQLiteDatabase() (*sql.DB, error) {
	dir, err := config.AppDataDir()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", path.Join(dir, config.AppDBFileName))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = prepare(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// prepare defines the database schema and executes it
// against the provided [db] reference.
func prepare(db *sql.DB) error {
	schema := `
		CREATE TABLE IF NOT EXISTS docpages (
			author TEXT NOT NULL,
			repository TEXT NOT NULL,
			relative_path TEXT NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			last_modified_at DATETIME NOT NULL,
			vector BLOB
		);
	`

	if _, err := db.Exec(schema); err != nil {
		slog.Error("Database schema preparation failed", "error", err)
		return err
	}

	return nil
}
