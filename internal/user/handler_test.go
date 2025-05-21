package user

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go_auth/internal/database"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	database.RunMigrations(db)
	return db
}

func TestRegisterUserHandler(t *testing.T) {
	db := setupTestDB()
	handler := RegisterUserHandler(db)

	user := User{
		Name:     "David Richard",
		Email:    "david@example.com",
		Password: "senha123",
	}

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("esperado status %d, obtido %d", http.StatusCreated, rr.Code)
	}
}

func TestRegisterUserHandler_EmailExists(t *testing.T) {
	db := setupTestDB()
	handler := RegisterUserHandler(db)

	user := User{
		Name:     "David Richard",
		Email:    "david@example.com",
		Password: "senha123",
	}

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("esperado status %d, obtido %d", http.StatusCreated, rr.Code)
	}
}
