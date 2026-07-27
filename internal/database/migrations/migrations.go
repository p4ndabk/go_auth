// Package migrations holds the ordered gormigrate migration list. Never
// edit an ID that has already shipped — anyone who ran it has that ID
// recorded as applied in the migrations table gormigrate creates. Append a
// new entry instead.
package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"go_auth/internal/access"
	"go_auth/internal/application"
	"go_auth/internal/permission"
	"go_auth/internal/role"
	"go_auth/internal/user"
)

var All = []*gormigrate.Migration{
	{
		ID: "20260101000001_create_users",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&user.User{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&user.User{})
		},
	},
	{
		ID: "20260101000002_create_applications",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&application.Application{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&application.Application{})
		},
	},
	{
		ID: "20260101000003_create_roles",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&role.Role{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&role.Role{})
		},
	},
	{
		ID: "20260101000004_create_permissions",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&permission.Permission{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&permission.Permission{})
		},
	},
	{
		ID: "20260101000005_create_role_permissions",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&role.Permission{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&role.Permission{})
		},
	},
	{
		ID: "20260101000006_create_user_application_role",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&access.UserApplicationRole{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&access.UserApplicationRole{})
		},
	},
}
