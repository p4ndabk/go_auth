package routes

import (
	"go_auth/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(env string) *gin.Engine {
	router := gin.Default()

	router.GET("/health", handlers.HealthHandler(env))

	return router
}
