package user

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UUID      string    `gorm:"size:36;uniqueIndex" json:"uuid"`
	Username  string    `gorm:"size:50;uniqueIndex" json:"username"`
	Email     string    `gorm:"size:255;uniqueIndex" json:"email"`
	Password  string    `gorm:"size:255" json:"-"`
	CreatedAt time.Time `json:"created_at"`
}
