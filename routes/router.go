package routes

import (
	"database/sql"
	"go_auth/internal/auth"
	"go_auth/internal/health"
	"go_auth/internal/user"
	"go_auth/internal/application"

	"github.com/gin-gonic/gin"
)

func SetupRouter(env string, db *sql.DB) *gin.Engine {
	router := gin.Default()

	router.GET("/health", health.HealthHandler(env))

	//user
	router.POST("/register", user.RegisterUserHandler(db))
	router.GET("/user/:id", user.ShowUserHandler(db))
	
	//login
	router.POST("/login", auth.LoginHandler(db))

	//application
	router.POST("/application", application.RegisterApplicationHandler(db))
	router.GET("/application/:id", application.ShowApplicationHandler(db))

	return router
}
