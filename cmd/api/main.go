// @title Go Auth API
// @version 1.0
// @description API de autenticação e autorização em Go com JWT, suportando múltiplas aplicações
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Digite "Bearer " seguido do token JWT

package main

import (
	"log"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "go_auth/docs"
	"go_auth/internal/access"
	"go_auth/internal/application"
	"go_auth/internal/auth"
	"go_auth/internal/config"
	"go_auth/internal/database"
	"go_auth/internal/docs"
	"go_auth/internal/health"
	"go_auth/internal/permission"
	"go_auth/internal/role"
	"go_auth/internal/user"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: .env file not found: %v", err)
	}

	cfg := config.Load()

	db, err := database.Connect(cfg.DBDriver, cfg.DBPath, cfg.DBDSN)
	if err != nil {
		log.Fatalf("failed to connect to %s database: %v", cfg.DBDriver, err)
	}
	log.Printf("connected to %s database successfully", cfg.DBDriver)

	userService := user.NewService(db)
	applicationService := application.NewService(db)
	roleService := role.NewService(db)
	permissionService := permission.NewService(db)
	accessService := access.NewService(db)
	authService := auth.New(cfg.JWTSecret)

	userHandler := user.NewHandler(userService, authService)
	applicationHandler := application.NewHandler(applicationService)
	roleHandler := role.NewHandler(roleService, applicationService)
	permissionHandler := permission.NewHandler(permissionService, applicationService)
	accessHandler := access.NewHandler(accessService, userService, applicationService, roleService, authService)
	healthHandler := health.NewHandler(db)

	router := gin.Default()
	router.Use(corsConfig(cfg.CORSAllowedOrigins))

	api := router.Group("/api")

	health.RegisterRoutes(api, healthHandler)
	docs.RegisterRoutes(api)

	user.RegisterPublicRoutes(api, userHandler)
	access.RegisterPublicRoutes(api, accessHandler)

	protected := api.Group("/")
	protected.Use(authService.AuthMiddleware())
	access.RegisterProtectedRoutes(protected, accessHandler)

	admin := api.Group("/admin")
	admin.Use(authService.AuthMiddleware())
	user.RegisterAdminRoutes(admin, userHandler)
	application.RegisterRoutes(admin, applicationHandler)
	role.RegisterRoutes(admin, roleHandler)
	permission.RegisterRoutes(admin, permissionHandler)
	access.RegisterAdminRoutes(admin, accessHandler)

	log.Printf("server starting on port %s", cfg.Port)
	log.Fatal(router.Run(":" + cfg.Port))
}

// corsConfig wires gin-contrib/cors from CORS_ALLOWED_ORIGINS ("*" or a
// comma-separated allowlist). Authorization is explicitly allowed since the
// API is JWT bearer-token based.
func corsConfig(allowedOrigins string) gin.HandlerFunc {
	cfg := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
	}

	if allowedOrigins == "*" {
		cfg.AllowAllOrigins = true
	} else {
		cfg.AllowOrigins = strings.Split(allowedOrigins, ",")
	}

	return cors.New(cfg)
}
