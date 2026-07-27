// Package health is an infra check, not a business domain: no
// model.go/service.go, the handler talks to *gorm.DB directly.
package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

// Check returns the health status of the API
// @Summary Health check
// @Description Returns the health status of the API and database connection
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/health [get]
func (h *Handler) Check(c *gin.Context) {
	sqlDB, err := h.DB.DB()
	if err != nil || sqlDB.Ping() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":    "unhealthy",
			"message":   "database connection failed",
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"message":   "API is running normally",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
