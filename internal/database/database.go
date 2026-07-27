// Package database owns the gorm.DB connection. Schema changes live in
// internal/database/migrations, applied via cmd/migrate — this package only
// opens the connection.
package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Connect opens a gorm.DB for the given driver ("sqlite" or "mysql").
// path is used for sqlite (the database file path), dsn for mysql.
// TranslateError is enabled so services can use errors.Is(err,
// gorm.ErrDuplicatedKey) to detect unique-constraint violations.
func Connect(driver, path, dsn string) (*gorm.DB, error) {
	cfg := &gorm.Config{TranslateError: true}

	switch driver {
	case "mysql":
		return gorm.Open(mysql.Open(dsn), cfg)
	case "sqlite":
		fallthrough
	default:
		if dir := filepath.Dir(path); dir != "." {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return nil, fmt.Errorf("failed to create database directory %s: %w", dir, err)
			}
		}
		return gorm.Open(sqlite.Open(path), cfg)
	}
}
