package access

import "github.com/gin-gonic/gin"

// RegisterPublicRoutes mounts the unauthenticated login endpoint.
func RegisterPublicRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("/login", h.Login)
}

// RegisterProtectedRoutes mounts endpoints that just need a valid token.
func RegisterProtectedRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.GET("/me", h.Me)
}

// RegisterAdminRoutes mounts the admin user-application-role endpoints.
func RegisterAdminRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.GET("/users/:id/access", h.UserAccess)

	rg.POST("/users/:id/application-roles", h.Assign)
	rg.GET("/users/:id/application-roles", h.ListForUser)
	rg.DELETE("/users/:id/application-roles/:applicationId/:roleId", h.Remove)

	rg.GET("/users/:id/applications", h.ListApplicationsForUser)
	rg.DELETE("/users/:id/applications/:applicationId", h.RemoveAllForUserInApplication)

	rg.GET("/applications/:id/users", h.ListUsersForApplication)

	rg.GET("/application-roles/users/:userId/applications/:applicationId/roles", h.ListRolesForUserInApplication)
}
