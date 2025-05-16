package database

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db   *sql.DB
	once sync.Once
)

func InitSQLite(dbPath string) *sql.DB {
	once.Do(func() {
		var err error
		db, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatalf("Erro ao conectar ao SQLite: %v", err)
		}
		
		if err := RunMigrations(db); err != nil {
			log.Fatalf("Erro ao executar migrações: %v", err)
		}

		if err = db.Ping(); err != nil {
			log.Fatalf("Erro ao verificar conexão com SQLite: %v", err)
		}

		log.Println("Conectado ao banco SQLite com sucesso!")
	})

	return db
}

func GetDB() *sql.DB {
	if db == nil {
		log.Fatal("Banco de dados ainda não inicializado. Chame InitSQLite primeiro.")
	}
	return db
}
