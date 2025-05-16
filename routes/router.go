package routes

import (
	"database/sql"
	"go_auth/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(env string, db *sql.DB) *gin.Engine {
	router := gin.Default()

	router.GET("/health", handlers.HealthHandler(env))

	router.POST("/register", gin.WrapF(handlers.RegisterUserHandler(db)))

	return router
}
