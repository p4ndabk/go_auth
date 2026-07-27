// Package access owns the user-application-role join: which role a user
// holds within a given application. It replaces the older, separately
// tracked user_roles / user_applications tables (see git history) with a
// single unified relationship.
package access

import (
	"time"

	"go_auth/internal/application"
)

// UserApplicationAccess is the aggregate view of what a user can do within
// one application — the shape embedded in the login response and the JWT.
type UserApplicationAccess struct {
	Application application.Application `json:"application"`
	Roles       []string                `json:"roles"`
	Permissions []string                `json:"permissions"`
}

type UserApplicationRole struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"uniqueIndex:idx_user_app_role" json:"user_id"`
	ApplicationID uint      `gorm:"uniqueIndex:idx_user_app_role" json:"application_id"`
	RoleID        uint      `gorm:"uniqueIndex:idx_user_app_role" json:"role_id"`
	ProfileID     *uint     `json:"profile_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

func (UserApplicationRole) TableName() string {
	return "user_application_role"
}
