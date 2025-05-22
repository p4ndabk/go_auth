package user

import (
	"database/sql"
)

func store(db *sql.DB, user *User) error {
	return CreateUser(db, user)
}

func show(db *sql.DB, id int) (*User, error) {
	return GetUserByID(db, id)
}