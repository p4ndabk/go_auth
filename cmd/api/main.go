// @title Go Auth API
// @version 1.0
// @description API de autenticação e autorização em Go com JWT, suportando múltiplas aplicações
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:9001
// @BasePath /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Digite "Bearer " seguido do token JWT

package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "go_auth/docs"
	"go_auth/internal/auth"
	"go_auth/internal/config"
	db "go_auth/internal/database"
	"go_auth/internal/handlers"
	"go_auth/internal/routes"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	var database *sql.DB
	var err error

	switch cfg.DatabaseType {
	case "mysql":
		database, err = sql.Open("mysql", cfg.DatabaseURL)
	case "sqlite":
		fallthrough
	default:
		database, err = sql.Open("sqlite3", cfg.DatabaseURL)
	}

	if err != nil {
		log.Fatalf("Failed to connect to %s database: %v", cfg.DatabaseType, err)
	}
	defer database.Close()

	// Test database connection
	if err := database.Ping(); err != nil {
		log.Fatalf("Failed to ping %s database: %v", cfg.DatabaseType, err)
	}

	log.Printf("Connected to %s database successfully", cfg.DatabaseType)

	// Initialize database schema
	if err := db.InitSchema(database, cfg.DatabaseType); err != nil {
		log.Fatal("Failed to initialize database schema:", err)
	}

	// Initialize services
	dbService := db.New(database)
	authService := auth.New(cfg.JWTSecret)

	// Initialize handlers
	h := handlers.New(dbService, authService)
	adminH := handlers.NewAdminHandler(dbService)

	// Setup routes
	router := routes.SetupRoutes(h, adminH, authService)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
