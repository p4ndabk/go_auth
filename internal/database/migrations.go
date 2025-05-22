package database

import (
	"database/sql"
	"log"
)

func RunMigrations(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT NOT NULL UNIQUE,
		external_uuid TEXT UNIQUE DEFAULT NULL,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		active BOOLEAN DEFAULT 1,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS applications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT NOT NULL UNIQUE,
		slug TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		description TEXT,
		active BOOLEAN DEFAULT 1,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	log.Println("Migration 'users' executada com sucesso.")
	return nil
}
