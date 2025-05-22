package database

import (
	"database/sql"

	"go_auth/internal/database/sqlite"
	"go_auth/pkg/logs"
)

func InitDB() *sql.DB {
	db := sqlite.InitSQLite("data.db")

	if err := RunMigrations(db); err != nil {
		logs.Error("Erro ao executar migrações: %v", err)
	}

	return db
}