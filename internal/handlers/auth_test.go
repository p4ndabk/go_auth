package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"go_auth/internal/handlers"
	"go_auth/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
	"go_auth/internal/database"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("erro abrindo banco: %v", err)
	}

	err = database.RunMigrations(db)
	if err != nil {
		t.Fatalf("erro criando tabela users: %v", err)
	}

	err = models.CreateUser(db, &models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "123456",
	})
	if err != nil {
		t.Fatalf("erro criando usuário teste: %v", err)
	}
	return db
}

func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	loginPayload := handlers.LoginRequest{
		Email:    "test@example.com",
		Password: "123456",
	}

	body, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r := gin.Default()
	r.POST("/login", handlers.LoginHandler(db))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp, "token")
	assert.NotEmpty(t, resp["token"])
}

func TestLoginHandlerInvalidPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	loginPayload := handlers.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	body, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r := gin.Default()
	r.POST("/login", handlers.LoginHandler(db))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginHandlerUserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	loginPayload := handlers.LoginRequest{
		Email:    "nouser@example.com",
		Password: "123456",
	}

	body, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r := gin.Default()
	r.POST("/login", handlers.LoginHandler(db))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
