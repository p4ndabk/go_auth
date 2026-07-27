package permission

import "github.com/gin-gonic/gin"

// RegisterRoutes mounts the admin permission endpoints.
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("/permissions", h.Create)
	rg.GET("/permissions", h.List)
	rg.GET("/permissions/:id", h.Get)
	rg.PUT("/permissions/:id", h.Update)
	rg.DELETE("/permissions/:id", h.Delete)
}
