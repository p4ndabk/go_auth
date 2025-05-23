package permission

import (
	"database/sql"
	"time"

	"go_auth/internal/application"
	"go_auth/pkg/logs"
)

type Permission struct {
	ID            int     `json:"id"`
	ApplicationID int     `json:"application_id"`
	Name          string  `json:"name"`
	Slug          string  `json:"slug"`
	Description   *string `json:"description,omitempty"`
	Active        bool    `json:"active"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	Application   *application.Application `json:"application,omitempty"`
}

func CreatePermission(db *sql.DB, p *Permission) error {
	p.CreatedAt = time.Now().Format(time.RFC3339)
	p.UpdatedAt = time.Now().Format(time.RFC3339)

	stmt, err := db.Prepare(`
		INSERT INTO permissions (application_id, name, slug, description, active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		logs.Error("create permission prepare statement", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(p.ApplicationID, p.Name, p.Slug, p.Description, p.Active, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		logs.Error("create permission exec statement", err)
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		logs.Error("create permission last insert id", err)
		return err
	}

	p.ID = int(lastID)
	return nil
}

func GetPermissionByID(db *sql.DB, id int) (*Permission, error) {
	row := db.QueryRow(`
		SELECT 
			p.id AS permission_id,
			p.application_id,
			p.name AS permission_name,
			p.slug AS permission_slug,
			p.description AS permission_description,
			p.active AS permission_active,
			p.created_at AS permission_created_at,
			p.updated_at AS permission_updated_at
		FROM permissions p
		WHERE p.id = ?
	`, id)

	var p Permission

	err := row.Scan(
		&p.ID,
		&p.ApplicationID,
		&p.Name,
		&p.Slug,
		&p.Description,
		&p.Active,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		logs.Error("get permission by id", err)
		return nil, err
	}

	return &p, nil
}

func GetAllPermissions(db *sql.DB) ([]*Permission, error) {
	rows, err := db.Query(`
		SELECT 
			p.id AS permission_id,
			p.application_id,
			p.name AS permission_name,
			p.slug AS permission_slug,
			p.description AS permission_description,
			p.active AS permission_active,
			p.created_at AS permission_created_at,
			p.updated_at AS permission_updated_at
		FROM permissions p
	`)
	if err != nil {
		logs.Error("get all permissions", err)
		return nil, err
	}
	
	defer rows.Close()

	var permissions []*Permission

	for rows.Next() {
		var p Permission
		err := rows.Scan(
			&p.ID,
			&p.ApplicationID,
			&p.Name,
			&p.Slug,
			&p.Description,
			&p.Active,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			logs.Error("scan permission", err)
			return nil, err
		}
		permissions = append(permissions, &p)
	}

	return permissions, nil
}

func UpdatePermission(db *sql.DB, p *Permission) error {
	p.UpdatedAt = time.Now().Format(time.RFC3339)

	_, err := db.Exec(`
		UPDATE permissions
		SET application_id = ?, name = ?, slug = ?, description = ?, active = ?, updated_at = ?
		WHERE id = ?
	`, p.ApplicationID, p.Name, p.Slug, p.Description, p.Active, p.UpdatedAt, p.ID)

	if err != nil {
		logs.Error("update permission", err)
		return err
	}

	return nil
}

func DeletePermission(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM permissions WHERE id = ?`, id)
	if err != nil {
		logs.Error("delete permission", err)
		return err
	}
	return nil
}
