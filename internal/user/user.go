package user

import (
	"database/sql"
	"go_auth/internal/permission"
	"go_auth/pkg/logs"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int     `json:"id"`
	UUID         string  `json:"uuid"`
	ExternalUUID *string `json:"external_uuid,omitempty"`
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	Password     string  `json:"password"`
	Active       bool    `json:"active"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

func CreateUser(db *sql.DB, user *User) error {
	user.UUID = uuid.New().String()
	hashedPassword, err := GenerateHashPassword(user.Password)
	if err != nil {
		logs.Error("create user hash password", err)
		return err
	}
	user.Password = string(hashedPassword)
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

func GenerateHashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
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

func GetUserByID(db *sql.DB, id int) (*User, error) {
	row := db.QueryRow("SELECT id, uuid, external_uuid, name, email, active, created_at, updated_at FROM users WHERE id = ?", id)

	var user User
	err := row.Scan(&user.ID, &user.UUID, &user.ExternalUUID, &user.Name, &user.Email, &user.Active, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		logs.Error("get user by id", err)
		return nil, err
	}

	return &user, nil
}

func (*User) GetUserPermissions(db *sql.DB, userID int) (*[]permission.Permission, error) {
	rows, err := db.Query(`
		SELECT p.id, p.name, p.application_id, p.slug, p.active, p.created_at, p.updated_at
			FROM user_application_role uar
			JOIN roles r ON r.id = uar.role_id
			JOIN role_permission rp ON rp.role_id = r.id
			JOIN permissions p ON p.id = rp.permission_id
			WHERE uar.user_id = ?`, userID)

	if err != nil {
		logs.Error("get user permissions", err)
		return nil, err
	}
	defer rows.Close()

	var permissions []permission.Permission
	for rows.Next() {
		var permission permission.Permission
		if err := rows.Scan(&permission.ID, &permission.Name, &permission.ApplicationID, &permission.Slug, &permission.Active, &permission.CreatedAt, &permission.UpdatedAt); err != nil {
			logs.Error("scan user permission", err)
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return &permissions, nil
}