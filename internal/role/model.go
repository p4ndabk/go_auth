package role

import "time"

type Role struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	ApplicationID uint   `gorm:"uniqueIndex:idx_role_app_slug" json:"application_id"`
	UUID          string `gorm:"size:36;uniqueIndex" json:"uuid"`
	Slug          string `gorm:"size:50;uniqueIndex:idx_role_app_slug" json:"slug"`
	Name          string `gorm:"size:100" json:"name"`
	Description   string `json:"description"`
	Active        bool   `json:"active"`
}

// Permission is a join record granting a permission to a role, scoped to
// the application both belong to.
type Permission struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	RoleID        uint      `gorm:"uniqueIndex:idx_role_permission" json:"role_id"`
	PermissionID  uint      `gorm:"uniqueIndex:idx_role_permission" json:"permission_id"`
	ApplicationID uint      `gorm:"uniqueIndex:idx_role_permission" json:"application_id"`
	CreatedAt     time.Time `json:"created_at"`
}

func (Permission) TableName() string {
	return "role_permissions"
}
