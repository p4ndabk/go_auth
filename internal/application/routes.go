package application

import "github.com/gin-gonic/gin"

// RegisterRoutes mounts the admin application endpoints.
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("/applications", h.Create)
	rg.GET("/applications", h.List)
	rg.GET("/applications/:id", h.Get)
	rg.PUT("/applications/:id", h.Update)
	rg.DELETE("/applications/:id", h.Delete)
}
