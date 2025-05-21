package auth

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"go_auth/internal/user"
	"go_auth/pkg/logs"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"go_auth/internal/database"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	logs.Init("test")
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("erro abrindo banco: %v", err)
	}

	err = database.RunMigrations(db)
	if err != nil {
		t.Fatalf("erro criando tabela users: %v", err)
	}

	err = user.CreateUser(db, &user.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "123456",
	})
	if err != nil {
		t.Fatalf("erro criando usuário teste: %v", err)
	}
	return db
}

type LoginResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
}

func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	loginPayload := LoginRequest{
		Email:    "test@example.com",
		Password: "123456",
	}

	body, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r := gin.Default()
	r.POST("/login", LoginHandler(db))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotEmpty(t, resp.Data.Token)
}

func TestLoginHandlerInvalidPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	loginPayload := LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	body, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r := gin.Default()
	r.POST("/login", LoginHandler(db))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginHandlerUserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	loginPayload := LoginRequest{
		Email:    "nouse1r@example.com",
		Password: "123456",
	}

	body, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r := gin.Default()
	r.POST("/login", LoginHandler(db))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
