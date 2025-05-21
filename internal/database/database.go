package database

import (
	"database/sql"
	"log"

	"go_auth/internal/database/sqlite"
)

func InitDB() *sql.DB {
	db := sqlite.InitSQLite("data.db")

	if err := RunMigrations(db); err != nil {
		log.Fatalf("Erro ao executar migrações: %v", err)
	}

	return db
}