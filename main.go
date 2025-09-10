// @title Go Auth API
// @version 1.0
// @description API de autenticação e autorização em Go com JWT, suportando múltiplas aplicações
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:9001
// @BasePath /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Digite "Bearer " seguido do token JWT

package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go_auth/auth"
	"go_auth/config"
	"go_auth/db"
	_ "go_auth/docs"
	"go_auth/handlers"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	var database *sql.DB
	var err error

	switch cfg.DatabaseType {
	case "mysql":
		database, err = sql.Open("mysql", cfg.DatabaseURL)
	case "sqlite":
		fallthrough
	default:
		database, err = sql.Open("sqlite3", cfg.DatabaseURL)
	}

	if err != nil {
		log.Fatalf("Failed to connect to %s database: %v", cfg.DatabaseType, err)
	}
	defer database.Close()

	// Test database connection
	if err := database.Ping(); err != nil {
		log.Fatalf("Failed to ping %s database: %v", cfg.DatabaseType, err)
	}

	log.Printf("Connected to %s database successfully", cfg.DatabaseType)

	// Initialize database schema
	if err := db.InitSchema(database, cfg.DatabaseType); err != nil {
		log.Fatal("Failed to initialize database schema:", err)
	}

	// Initialize services
	dbService := db.New(database)
	authService := auth.New(cfg.JWTSecret)

	// Initialize handlers
	h := handlers.New(dbService, authService)
	adminH := handlers.NewAdminHandler(dbService)

	// Setup Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Public routes
	router.POST("/register", h.Register)
	router.POST("/login", h.Login)
	router.GET("/health", h.Health)

	// Protected routes
	protected := router.Group("/")
	protected.Use(authService.AuthMiddleware())
	{
		protected.GET("/me", h.Me)
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Admin routes (also protected)
	admin := router.Group("/admin")
	admin.Use(authService.AuthMiddleware())
	{
		// Application routes
		admin.POST("/applications", adminH.CreateApplication)
		admin.GET("/applications", adminH.GetApplications)
		admin.GET("/applications/:id", adminH.GetApplication)
		admin.PUT("/applications/:id", adminH.UpdateApplication)
		admin.DELETE("/applications/:id", adminH.DeleteApplication)

		// Role routes
		admin.POST("/roles", adminH.CreateRole)
		admin.GET("/roles", adminH.GetRoles) // ?application_id=1 for filtering
		admin.GET("/roles/:id", adminH.GetRole)
		admin.PUT("/roles/:id", adminH.UpdateRole)
		admin.DELETE("/roles/:id", adminH.DeleteRole)

		// Permission routes
		admin.POST("/permissions", adminH.CreatePermission)
		admin.GET("/permissions", adminH.GetPermissions) // ?application_id=1 for filtering
		admin.GET("/permissions/:id", adminH.GetPermission)
		admin.PUT("/permissions/:id", adminH.UpdatePermission)
		admin.DELETE("/permissions/:id", adminH.DeletePermission)

		// Role-Permission relationship routes
		admin.POST("/roles/:id/permissions", adminH.AssignPermissionToRole)
		admin.DELETE("/roles/:id/permissions/:permissionId", adminH.RemovePermissionFromRole)
		admin.GET("/roles/:id/permissions", adminH.GetRolePermissions)

		// User-Role relationship routes
		admin.POST("/users/:id/roles", adminH.AssignRoleToUser)
		admin.DELETE("/users/:id/roles/:roleId", adminH.RemoveRoleFromUser)

		// User-Application relationship routes
		admin.POST("/users/:id/applications", adminH.AssignApplicationToUser)
		admin.DELETE("/users/:id/applications/:applicationId", adminH.RemoveApplicationFromUser)
		admin.GET("/users/:id/applications", adminH.GetUserApplications)

		// Application-User relationship routes
		admin.GET("/applications/:id/users", adminH.GetApplicationUsers)

		// User-Application-Role relationship routes (New Model)
		admin.POST("/users/:id/application-roles", adminH.AssignUserApplicationRole)
		admin.DELETE("/users/:id/application-roles/:applicationId/:roleId", adminH.RemoveUserApplicationRole)
		admin.GET("/users/:id/application-roles", adminH.GetUserApplicationRoles)
		admin.GET("/application-roles/users/:userId/applications/:applicationId/roles", adminH.GetUserRolesInApplication)
	}

	// Setup server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
