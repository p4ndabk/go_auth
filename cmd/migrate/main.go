package main

import (
	"flag"
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/joho/godotenv"

	"go_auth/internal/config"
	"go_auth/internal/database"
	"go_auth/internal/database/migrations"
)

func main() {
	rollback := flag.Bool("rollback", false, "roll back the most recently applied migration")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Printf("warning: .env file not found: %v", err)
	}

	cfg := config.Load()

	db, err := database.Connect(cfg.DBDriver, cfg.DBPath, cfg.DBDSN)
	if err != nil {
		log.Fatalf("failed to connect to %s database: %v", cfg.DBDriver, err)
	}

	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations.All)

	if *rollback {
		if err := m.RollbackLast(); err != nil {
			log.Fatalf("rollback failed: %v", err)
		}
		log.Println("rolled back last migration successfully")
		return
	}

	if err := m.Migrate(); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	log.Println("migrations applied successfully")
}
