package permission

type Permission struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	ApplicationID uint   `gorm:"uniqueIndex:idx_permission_app_slug" json:"application_id"`
	Name          string `gorm:"size:100" json:"name"`
	Slug          string `gorm:"size:50;uniqueIndex:idx_permission_app_slug" json:"slug"`
	Description   string `json:"description"`
	Active        bool   `json:"active"`
}
