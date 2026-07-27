// cmd/seed populates a freshly migrated database with the same sample data
// the project has always shipped with (see former migrations/002_seed_data.sql):
// a main-app application, admin/user/moderator roles, a handful of
// permissions, and an admin user (admin@admin.com / admin123).
package main

import (
	"log"

	"github.com/joho/godotenv"

	"go_auth/internal/access"
	"go_auth/internal/application"
	"go_auth/internal/auth"
	"go_auth/internal/config"
	"go_auth/internal/database"
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

	userService := user.NewService(db)
	applicationService := application.NewService(db)
	roleService := role.NewService(db)
	permissionService := permission.NewService(db)
	accessService := access.NewService(db)
	authService := auth.New(cfg.JWTSecret)

	app, err := applicationService.Create("main-app", "Main Application", "The main application")
	if err != nil {
		log.Fatalf("failed to create application: %v", err)
	}

	adminRole, err := roleService.Create(app.ID, "admin", "Administrator", "Full system access")
	if err != nil {
		log.Fatalf("failed to create admin role: %v", err)
	}
	userRole, err := roleService.Create(app.ID, "user", "User", "Standard user access")
	if err != nil {
		log.Fatalf("failed to create user role: %v", err)
	}
	moderatorRole, err := roleService.Create(app.ID, "moderator", "Moderator", "Moderate content")
	if err != nil {
		log.Fatalf("failed to create moderator role: %v", err)
	}

	permissionDefs := []struct{ name, slug, description string }{
		{"Read Users", "read_users", "Can view user information"},
		{"Write Users", "write_users", "Can create and edit users"},
		{"Delete Users", "delete_users", "Can delete users"},
		{"Read Posts", "read_posts", "Can view posts"},
		{"Write Posts", "write_posts", "Can create and edit posts"},
		{"Delete Posts", "delete_posts", "Can delete posts"},
		{"Moderate Content", "moderate_content", "Can moderate user content"},
	}
	perms := make([]*permission.Permission, len(permissionDefs))
	for i, def := range permissionDefs {
		p, err := permissionService.Create(app.ID, def.name, def.slug, def.description)
		if err != nil {
			log.Fatalf("failed to create permission %s: %v", def.slug, err)
		}
		perms[i] = p
	}

	grant := func(r *role.Role, slugs ...string) {
		for _, slug := range slugs {
			for _, p := range perms {
				if p.Slug == slug {
					if err := roleService.AssignPermission(r.ID, p.ID); err != nil {
						log.Fatalf("failed to assign %s to role %s: %v", slug, r.Slug, err)
					}
				}
			}
		}
	}
	grant(adminRole, "read_users", "write_users", "delete_users", "read_posts", "write_posts", "delete_posts", "moderate_content")
	grant(userRole, "read_users", "read_posts", "write_posts")
	grant(moderatorRole, "read_users", "read_posts", "write_posts", "moderate_content")

	hashedPassword, err := authService.HashPassword("admin123")
	if err != nil {
		log.Fatalf("failed to hash admin password: %v", err)
	}
	admin, err := userService.Create("admin", "admin@admin.com", hashedPassword)
	if err != nil {
		log.Fatalf("failed to create admin user: %v", err)
	}

	if err := accessService.Assign(admin.ID, app.ID, adminRole.ID, nil); err != nil {
		log.Fatalf("failed to assign admin role to admin user: %v", err)
	}

	log.Println("database seeded successfully")
}
