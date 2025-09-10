package config

import (
	"os"
	"strings"
)

type Config struct {
	Port         string
	DatabaseType string
	DatabaseURL  string
	JWTSecret    string
}

func Load() *Config {
	dbType := strings.ToLower(getEnv("DATABASE_TYPE", "sqlite"))
	var defaultURL string

	switch dbType {
	case "mysql":
		defaultURL = "user:password@tcp(localhost:3306)/auth_db?charset=utf8mb4&parseTime=True&loc=Local"
	case "sqlite":
		fallthrough
	default:
		defaultURL = "./auth.db"
		dbType = "sqlite"
	}

	return &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseType: dbType,
		DatabaseURL:  getEnv("DATABASE_URL", defaultURL),
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
