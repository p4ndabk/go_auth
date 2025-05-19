package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            int    `json:"id"`
	UUID          string `json:"uuid"`
	ExternalUUID  string `json:"external_uuid,omitempty"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Active        bool   `json:"active"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

func CreateUser(db *sql.DB, user *User) error {
	user.UUID = uuid.New().String()
	user.CreatedAt = time.Now().Format(time.RFC3339)
	user.UpdatedAt = time.Now().Format(time.RFC3339)

	stmt, err := db.Prepare("INSERT INTO users (uuid, external_uuid, name, email, password, active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.UUID, user.ExternalUUID, user.Name, user.Email, user.Password, user.Active, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(lastID)
	return nil
}

func EmailExists(db *sql.DB, email string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	row := db.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", email)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}