package management

import (
	"database/sql"
	"log"
	"path"

	"github.com/inodaf/neoman/packages/config"
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
	log.Println("db: prepare")
	schema := ``

	_, err := db.Exec(schema)
	if err != nil {
		return err
	}

	log.Println("db: prepare success")
	return nil
}
