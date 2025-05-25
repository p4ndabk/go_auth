package routes

import (
	"database/sql"
	"go_auth/internal/auth"
	"go_auth/internal/health"
	"go_auth/internal/user"
	"go_auth/internal/application"

	"github.com/gin-gonic/gin"
	"go_auth/internal/permission"
	"go_auth/internal/role"
)

func SetupRouter(env string, db *sql.DB) *gin.Engine {
	router := gin.Default()

	router.GET("/health", health.HealthHandler(env))

	//user
	router.POST("/register", user.RegisterUserHandler(db))
	router.GET("/user/:id", user.ShowUserHandler(db))
	router.POST("/user-application-role", user.CreateUserApplicationRoleHandler(db))

	//login
	router.POST("/login", auth.LoginHandler(db))

	//application
	router.POST("/application", application.RegisterApplicationHandler(db))
	router.GET("/application/:id", application.ShowApplicationHandler(db))

	//permission
	router.GET("/permission", permission.IndexPermissionHandler(db))
	router.POST("/permission", permission.RegisterPermissionHandler(db))
	router.GET("/permission/:id", permission.ShowPermissionHandler(db))
	router.PUT("/permission/:id", permission.UpdatePermissionHandler(db))
	router.DELETE("/permission/:id", permission.DeletePermissionHandler(db))

	router.GET("/role", role.IndexRolesHandler(db))
	router.POST("/role", role.RegisterRoleHandler(db))
	router.GET("/role/:id", role.ShowRoleHandler(db))
	router.PUT("/role/:id", role.UpdateRoleHandler(db))
	router.DELETE("/role/:id", role.DeleteRoleHandler(db))
	router.POST("/role/:id/attach-permissions", role.AttachPermissionsHandler(db))


	return router
}
