package internal

import (
	"database/sql"
	"log"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

// DB holds a globally accessible reference to the initialized
// *sql.DB. It is agnostic of the SQL driver.
var DB *sql.DB

// NewSQLiteDatabase creates a connection to a SQL DB
// and prepares the base schema.
func NewSQLiteDatabase() (*sql.DB, error) {
	dir, err := AppDataDir()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", path.Join(dir, AppDBFileName))
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
	schema := `CREATE TABLE IF NOT EXISTS TrustedAccount (
		account TEXT NOT NULL
	)`

	log.Println("db: prepare")
	_, err := db.Exec(schema)
	if err != nil {
		return err
	}

	log.Println("db: prepare success")
	return nil
}
