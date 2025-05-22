package application

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"go_auth/pkg/logs"
)

type Application struct {
	ID          int     `json:"id"`
	UUID        string  `json:"uuid"`
	Slug        string  `json:"slug"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Active      bool    `json:"active"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func CreateApplication(db *sql.DB, app *Application) error {
	app.UUID = uuid.New().String()
	app.CreatedAt = time.Now().Format(time.RFC3339)
	app.UpdatedAt = time.Now().Format(time.RFC3339)

	stmt, err := db.Prepare(`
		INSERT INTO applications (uuid, slug, name, description, active, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		logs.Error("create application prepare statement", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(app.UUID, app.Slug, app.Name, app.Description, app.Active, app.CreatedAt, app.UpdatedAt)
	if err != nil {
		logs.Error("create application exec statement", err)
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		logs.Error("create application last insert id", err)
		return err
	}

	app.ID = int(lastID)
	return nil
}


func GetApplicationByID(db *sql.DB, id int) (*Application, error) {
	row := db.QueryRow(`
		SELECT id, uuid, slug, name, description, active, created_at, updated_at 
		FROM applications WHERE id = ?
	`, id)

	var app Application
	err := row.Scan(
		&app.ID,
		&app.UUID,
		&app.Slug,
		&app.Name,
		&app.Description,
		&app.Active,
		&app.CreatedAt,
		&app.UpdatedAt,
	)
	if err != nil {
		logs.Error("get application by id", err)
		return nil, err
	}

	return &app, nil
}

func UpdateApplication(db *sql.DB, app *Application) error {
	app.UpdatedAt = time.Now().Format(time.RFC3339)

	_, err := db.Exec(`
		UPDATE applications 
		SET slug = ?, name = ?, description = ?, active = ?, updated_at = ?
		WHERE id = ?
	`, app.Slug, app.Name, app.Description, app.Active, app.UpdatedAt, app.ID)

	if err != nil {
		logs.Error("update application", err)
		return err
	}

	return nil
}

func DeleteApplication(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM applications WHERE id = ?", id)
	if err != nil {
		logs.Error("delete application", err)
		return err
	}
	return nil
}


