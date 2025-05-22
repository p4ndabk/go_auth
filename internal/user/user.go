package user

import (
	"database/sql"
	"go_auth/pkg/logs"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            int    `json:"id"`
	UUID          string `json:"uuid"`
	ExternalUUID  *string `json:"external_uuid,omitempty"`
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
		logs.Error("create user prepare statement", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.UUID, user.ExternalUUID, user.Name, user.Email, user.Password, user.Active, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		logs.Error("create user exec statement", err)
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		logs.Error("create user last insert id", err)
		return err
	}

	user.ID = int(lastID)
	return nil
}

func EmailExists(db *sql.DB, email string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&exists)
	if err != nil {
		logs.Error("check email exists", err)
		return false, err
	}
	return exists, nil
}

func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	row := db.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", email)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		logs.Error("get user by email", err)
		return nil, err
	}

	return &user, nil
}

func GetUserByID(db *sql.DB, id string) (*User, error) {
	row := db.QueryRow("SELECT id, uuid, external_uuid, name, email, active, created_at, updated_at FROM users WHERE id = ?", id)

	var user User
	err := row.Scan(&user.ID, &user.UUID, &user.ExternalUUID, &user.Name, &user.Email, &user.Active, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		logs.Error("get user by id", err)
		return nil, err
	}

	return &user, nil
}