package health

import (
	"go_auth/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthHandler(env string) gin.HandlerFunc {
	return func(c *gin.Context) {
		response.Success(c.Writer, http.StatusOK, gin.H{
			"status": "ok",
			"env":    env,
		})
	}
}
