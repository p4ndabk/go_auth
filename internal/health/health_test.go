package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	env := "test"

	router := gin.Default()
	router.GET("/health", HealthHandler(env))

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expected := `{
    "success": true,
		"data": {
			"env": "test",
			"status": "ok"
		}
	}`
	assert.JSONEq(t, expected, w.Body.String())
}
