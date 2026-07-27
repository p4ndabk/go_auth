package user

import "github.com/gin-gonic/gin"

// RegisterPublicRoutes mounts the unauthenticated user endpoints.
func RegisterPublicRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("/register", h.Register)
}

// RegisterAdminRoutes mounts the admin-only user endpoints.
func RegisterAdminRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.GET("/users", h.List)
	rg.GET("/users/:id", h.Get)
}
