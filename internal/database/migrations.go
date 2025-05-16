package database

import (
	"database/sql"
	"log"
)

func RunMigrations(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	log.Println("Migration 'users' executada com sucesso.")
	return nil
}
