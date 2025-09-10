package routes

import (
	"go_auth/internal/auth"
	"go_auth/internal/handlers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes configura todas as rotas da aplicação
func SetupRoutes(h *handlers.Handler, adminH *handlers.AdminHandler, authService *auth.Service) *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(corsMiddleware())

	// Public routes
	setupPublicRoutes(router, h)

	// Protected routes
	setupProtectedRoutes(router, h, authService)

	// Admin routes
	setupAdminRoutes(router, adminH, authService)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

// corsMiddleware configura o middleware de CORS
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// setupPublicRoutes configura as rotas públicas (sem autenticação)
func setupPublicRoutes(router *gin.Engine, h *handlers.Handler) {
	router.POST("/register", h.Register)
	router.POST("/login", h.Login)
	router.GET("/health", h.Health)
}

// setupProtectedRoutes configura as rotas protegidas (com autenticação)
func setupProtectedRoutes(router *gin.Engine, h *handlers.Handler, authService *auth.Service) {
	protected := router.Group("/")
	protected.Use(authService.AuthMiddleware())
	{
		protected.GET("/me", h.Me)
	}
}

// setupAdminRoutes configura as rotas administrativas
func setupAdminRoutes(router *gin.Engine, adminH *handlers.AdminHandler, authService *auth.Service) {
	admin := router.Group("/admin")
	admin.Use(authService.AuthMiddleware())
	{
		// User routes
		setupUserRoutes(admin, adminH)

		// Application routes
		setupApplicationRoutes(admin, adminH)

		// Role routes
		setupRoleRoutes(admin, adminH)

		// Permission routes
		setupPermissionRoutes(admin, adminH)

		// Role-Permission relationship routes
		setupRolePermissionRoutes(admin, adminH)

		// User-Role relationship routes
		setupUserRoleRoutes(admin, adminH)

		// User-Application relationship routes
		setupUserApplicationRoutes(admin, adminH)

		// Application-User relationship routes
		setupApplicationUserRoutes(admin, adminH)

		// User-Application-Role relationship routes (New Model)
		setupUserApplicationRoleRoutes(admin, adminH)
	}
}

// setupApplicationRoutes configura as rotas de aplicações
func setupApplicationRoutes(admin *gin.RouterGroup, adminH *handlers.AdminHandler) {
	admin.POST("/applications", adminH.CreateApplication)
	admin.GET("/applications", adminH.GetApplications)
	admin.GET("/applications/:id", adminH.GetApplication)
	admin.PUT("/applications/:id", adminH.UpdateApplication)
	admin.DELETE("/applications/:id", adminH.DeleteApplication)
}

// setupRoleRoutes configura as rotas de roles
func setupRoleRoutes(admin *gin.RouterGroup, adminH *handlers.AdminHandler) {
	admin.POST("/roles", adminH.CreateRole)
	admin.GET("/roles", adminH.GetRoles) // ?application_id=1 for filtering
	admin.GET("/roles/:id", adminH.GetRole)
	admin.PUT("/roles/:id", adminH.UpdateRole)
	admin.DELETE("/roles/:id", adminH.DeleteRole)
}

// setupPermissionRoutes configura as rotas de permissões
func setupPermissionRoutes(admin *gin.RouterGroup, adminH *handlers.AdminHandler) {
	admin.POST("/permissions", adminH.CreatePermission)
	admin.GET("/permissions", adminH.GetPermissions) // ?application_id=1 for filtering
	admin.GET("/permissions/:id", adminH.GetPermission)
	admin.PUT("/permissions/:id", adminH.UpdatePermission)
	admin.DELETE("/permissions/:id", adminH.DeletePermission)
}

// setupRolePermissionRoutes configura as rotas de relacionamento role-permission
func setupRolePermissionRoutes(admin *gin.RouterGroup, adminH *handlers.AdminHandler) {
	admin.POST("/roles/:id/permissions", adminH.AssignPermissionToRole)
	admin.DELETE("/roles/:id/permissions/:permissionId", adminH.RemovePermissionFromRole)
	admin.GET("/roles/:id/permissions", adminH.GetRolePermissions)
}

// setupUserRoleRoutes configura as rotas de relacionamento user-role
func setupUserRoleRoutes(admin *gin.RouterGroup, adminH *handlers.AdminHandler) {
	admin.POST("/users/:id/roles", adminH.AssignRoleToUser)
	admin.DELETE("/users/:id/roles/:roleId", adminH.RemoveRoleFromUser)
}

// setupUserApplicationRoutes configura as rotas de relacionamento user-application
func setupUserApplicationRoutes(admin *gin.RouterGroup, adminH *handlers.AdminHandler) {
	admin.POST("/users/:id/applications", adminH.AssignApplicationToUser)
	admin.DELETE("/users/:id/applications/:applicationId", adminH.RemoveApplicationFromUser)
	admin.GET("/users/:id/applications", adminH.GetUserApplications)
}

// setupApplicationUserRoutes configura as rotas de relacionamento application-user
func setupApplicationUserRoutes(admin *gin.RouterGroup, adminH *handlers.AdminHandler) {
	admin.GET("/applications/:id/users", adminH.GetApplicationUsers)
}

// setupUserApplicationRoleRoutes configura as rotas do novo modelo user-application-role
func setupUserApplicationRoleRoutes(admin *gin.RouterGroup, adminH *handlers.AdminHandler) {
	admin.POST("/users/:id/application-roles", adminH.AssignUserApplicationRole)
	admin.DELETE("/users/:id/application-roles/:applicationId/:roleId", adminH.RemoveUserApplicationRole)
	admin.GET("/users/:id/application-roles", adminH.GetUserApplicationRoles)
	admin.GET("/application-roles/users/:userId/applications/:applicationId/roles", adminH.GetUserRolesInApplication)
}

// setupUserRoutes configura as rotas de usuários
func setupUserRoutes(admin *gin.RouterGroup, adminH *handlers.AdminHandler) {
	admin.GET("/users", adminH.GetUsers)
}
