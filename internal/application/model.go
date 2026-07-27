package application

type Application struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	UUID        string `gorm:"size:36;uniqueIndex" json:"uuid"`
	Slug        string `gorm:"size:50;uniqueIndex" json:"slug"`
	Name        string `gorm:"size:100" json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}
