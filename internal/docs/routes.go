// Package docs mounts the Swagger UI. Infra, no model/service — same
// exception as health.
package docs

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
