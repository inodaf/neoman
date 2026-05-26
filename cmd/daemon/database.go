package main

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"path"

	"github.com/inodaf/neoman/pkg/config"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

//go:embed db/migrations/*.sql
var migrations embed.FS

func InitDB() (*sql.DB, error) {
	dir, err := config.AppDataDir()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", path.Join(dir, config.AppDBFileName))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	migrationsFS, err := fs.Sub(migrations, "db/migrations")
	if err != nil {
		return nil, err
	}

	provider, err := goose.NewProvider(goose.DialectSQLite3, db, migrationsFS)
	if err != nil {
		return nil, err
	}

	if _, err := provider.Up(context.Background()); err != nil {
		return nil, err
	}

	return db, nil
}
