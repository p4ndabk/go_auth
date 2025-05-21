package routes

import (
	"database/sql"
	"go_auth/internal/auth"
	"go_auth/internal/health"
	"go_auth/internal/user"

	"github.com/gin-gonic/gin"
)

func SetupRouter(env string, db *sql.DB) *gin.Engine {
	router := gin.Default()

	router.GET("/health", health.HealthHandler(env))

	router.POST("/register", gin.WrapF(user.RegisterUserHandler(db)))
	router.POST("/login", auth.LoginHandler(db))

	return router
}
