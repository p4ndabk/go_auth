package role

import "github.com/gin-gonic/gin"

// RegisterRoutes mounts the admin role and role-permission endpoints.
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("/roles", h.Create)
	rg.GET("/roles", h.List)
	rg.GET("/roles/:id", h.Get)
	rg.PUT("/roles/:id", h.Update)
	rg.DELETE("/roles/:id", h.Delete)

	rg.POST("/roles/:id/permissions", h.AssignPermission)
	rg.GET("/roles/:id/permissions", h.ListPermissions)
	rg.DELETE("/roles/:id/permissions/:permissionId", h.RemovePermission)
}
