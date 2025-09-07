package database

import (
	"database/sql"
	"embed"
	"log"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations runs all database migrations
func RunMigrations(db *sql.DB) error {
	// Set the base filesystem for migrations
	goose.SetBaseFS(migrationsFS)

	// Run migrations
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}
