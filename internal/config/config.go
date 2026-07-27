package config

import "os"

type Config struct {
	Port               string
	DBDriver           string
	DBPath             string
	DBDSN              string
	JWTSecret          string
	CORSAllowedOrigins string
}

func Load() *Config {
	return &Config{
		Port:               getEnv("PORT", "8080"),
		DBDriver:           getEnv("DB_DRIVER", "sqlite"),
		DBPath:             getEnv("DB_PATH", "data/go_auth.db"),
		DBDSN:              getEnv("DB_DSN", "user:password@tcp(localhost:3306)/go_auth?charset=utf8mb4&parseTime=True&loc=Local"),
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
